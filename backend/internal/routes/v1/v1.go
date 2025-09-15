package v1_routes

import (
	"io"
	"os"
	"time"

	ws "backend/internal/controllers/v1/ws"
	"backend/internal/redis"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func V1RoutesRegister(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	hub := ws.NewHub()
	go hub.Run()

	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.Rdb,
		Rate:        time.Minute,
		Limit:       100,
	})

	r.Use(ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
			c.JSON(429, gin.H{"error": "Too many requests"})
		},
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	}))
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r.Use(gin.LoggerWithWriter(f))

	v1 := r.Group("/v1")
	v1.GET("", func(c *gin.Context) {
		c.String(200, "OK")
	})
	userRoutes(v1)
	authRoutes(v1)
	customerRoutes(v1)
	tenantRoutes(v1)
	ticketRoutes(v1, hub)
	invoiceRoutes(v1)

}
