package repositories

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool  *pgxpool.Pool
	queries *Queries
	once    sync.Once
)

func InitDB() *Queries {

	once.Do(func() {
		log.Printf("InitDB - Starting database initialization")
		connString := os.Getenv("DB_URL")
		if connString == "" {
			log.Fatal("DB_URL environment variable is not set")
		}

		log.Printf("InitDB - Parsing database configuration")
		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			log.Fatal("Unable to parse config:", err)
		}
		config.MaxConns = 50
		config.MinConns = 2
		config.MaxConnIdleTime = 15 * time.Minute
		config.MaxConnLifetime = 1 * time.Hour
		config.HealthCheckPeriod = 2 * time.Minute

		log.Printf("InitDB - Creating connection pool with MaxConns: %d, MinConns: %d", config.MaxConns, config.MinConns)
		dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatal("Unable to connect to database:", err)
		}

		log.Printf("InitDB - Testing database connection")
		if err := dbPool.Ping(context.Background()); err != nil {
			log.Fatal("Unable to ping database:", err)
		}

		queries = New(dbPool)
		log.Println("✅ Database initialized and sqlc queries ready")
	})

	return queries
}

func GetDB() *Queries {
	return queries
}

func Close() {
	if dbPool != nil {
		dbPool.Close()
	}
}
