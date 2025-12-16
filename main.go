package main

import (
	"fmt"
	"log"
	"manara/database"
	"manara/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("❌ Error loading %s file: %v", envFile, err)
	}

	log.Printf("🚀 Starting Manara Backend in %s mode...", env)

	err = database.Connect()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer database.Close()

	if os.Getenv("APP_DEBUG") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	routes.AuthRoutes(router)
	routes.RoleRoutes(router)
	routes.TeacherRoutes(router)
	routes.StudentRoutes(router)
	routes.AcademicYearRoutes(router)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Server is running on port %s", port)
	log.Printf("🌍 Environment: %s", env)
	router.Run(":" + port)
}
