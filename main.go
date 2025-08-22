package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/models"
)

func mustEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if def != "" {
		return def
	}
	log.Fatalf("missing env %s", key)
	return ""
}

func main() {
	// 1) load .env (dev)
	_ = godotenv.Load()

	// 2) config
	port := mustEnv("APP_PORT", "8080")
	dsn := mustEnv("DATABASE_URL", "")

	// 3) connect DB (GORM + Postgres)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	// 3b) AutoMigrate AFTER successful connect
	if err := db.AutoMigrate(
		&models.User{},
		&models.Tweet{},
		&models.Follow{},
		&models.Like{},
		&models.EditHistory{},
	); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	// optional: test ping at startup (open sql.DB)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("DB handle error: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	// 4) HTTP server (Gin)
	r := gin.Default()

	// CORS (longgar untuk dev)
	c := cors.DefaultConfig()
	c.AllowOrigins = []string{
		"http://localhost:3000", // Next.js default
		"http://localhost:5173", // Vite default
	}
	c.AllowHeaders = append(c.AllowHeaders, "Authorization")
	r.Use(cors.New(c))

	// 5) routes
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from Go backend! 🚀")
	})

	// health check + DB ping
	r.GET("/healthz", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "db": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "db": "up"})
	})

	// 6) run
	addr := ":" + port
	log.Println("listening on " + addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
