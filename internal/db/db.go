// internal/db/db.go
package db

import (
	"chat-system/internal/db/migrations"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *sql.DB
	GormDB *gorm.DB
	Redis  *redis.Client
)

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to MySQL
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),     // Username
		os.Getenv("DB_PASSWORD"), // Password
		os.Getenv("DB_HOST"),     // Host (e.g., "db")
		os.Getenv("DB_PORT"),     // Port (e.g., "3306")
		os.Getenv("DB_NAME"))     // Database name (e.g., "chat_system")

	dsn := dbURL + "?charset=utf8mb4&parseTime=True&loc=Local"

	var db *gorm.DB
	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database after 3 attempts: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get underlying *sql.DB: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	DB = sqlDB
	GormDB = db

	// Connect to Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	migrations.RunMigrations(db)

	// Connect to Elasticsearch
	setupElasticsearch()
	log.Println("Connected to database and services")
}
