package main

import (
	"backend/internal/redis"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/jobs"
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env once
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system env vars")
	}
}

func gracefulShutdown(apiServer *http.Server, asynqServer *asynq.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// Shutdown asynq server first
	log.Println("Shutting down asynq server...")
	asynqServer.Shutdown()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	repositories.InitDB()
	redis.InitRedis()
	defer repositories.Close()
	defer redis.Close()

	// Initialize job logger
	if err := jobs.InitLogger("logs/jobs.log"); err != nil {
		log.Fatalf("Failed to initialize job logger: %v", err)
	}
	defer jobs.CloseLogger()

	// Redis configuration for asynq
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}

	// Create asynq client for enqueueing jobs
	client := asynq.NewClient(redisOpt)
	defer client.Close()

	// Create asynq server for processing jobs
	asynqServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5, // Number of concurrent workers
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TypePDFInvoice, jobs.HandlePDFTask)

	// Create HTTP server with routes (pass asynq client to routes)
	server := &http.Server{
		Addr:           ":8080",
		Handler:        routes.RegisterRoutes(), // Pass client to routes
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Use WaitGroup to manage both HTTP server and asynq worker
	var wg sync.WaitGroup
	done := make(chan bool, 1)

	// Start graceful shutdown handler
	go gracefulShutdown(server, asynqServer, done)

	// Start the asynq worker in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("🔄 Starting asynq worker...")
		log.Printf("📋 Worker ready to process tasks: %s", jobs.TypePDFInvoice)
		if err := asynqServer.Run(mux); err != nil {
			log.Printf("❌ Asynq worker error: %v", err)
		}
		log.Println("🔄 Asynq worker stopped")
	}()

	// Start the HTTP server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("🚀 Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
		log.Println("🚀 HTTP server stopped")
	}()

	// Wait for shutdown signal
	<-done

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("✅ Graceful shutdown complete.")
}
