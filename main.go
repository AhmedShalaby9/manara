package main

import (
	"flag"
	"fmt"
	"log"
	"manara/database"
	"manara/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

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

	// Connect to database FIRST
	err = database.Connect()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer database.Close()

	// Run migrations if flag is set (AFTER connection)
	if *migrate {
		err = database.AutoMigrate()
		if err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
		log.Println("✅ Migrations completed. Exiting.")
		return
	}

	if os.Getenv("APP_DEBUG") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Serve static files (uploaded images)
	router.Static("/uploads", "./uploads")

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

	api := router.Group("/api")

	routes.AuthRoutes(api)
	routes.RoleRoutes(api)
	routes.TeacherRoutes(api)
	routes.StudentRoutes(api)
	routes.AcademicYearRoutes(api)
	routes.CourseRoutes(api)
	routes.ChapterRoutes(api)
	routes.LessonRoutes(api)
	routes.UserRoutes(api)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Server is running on port %s", port)
	log.Printf("🌍 Environment: %s", env)
	router.Run(":" + port)
}
