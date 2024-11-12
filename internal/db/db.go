// internal/db/db.go
package db

import (
	"chat-system/internal/db/migrations"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *sql.DB
	GormDB *gorm.DB
)

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to MySQL
	dbURL := os.Getenv("DATABASE_URL")
	dsn := dbURL + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
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

	migrations.RunMigrations(db)
	log.Println("Connected to database and services")
}
