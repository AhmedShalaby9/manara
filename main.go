package main

import (
	"flag"
	"fmt"
	"log"
	"manara/database"
	"manara/helpers"
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

	// Initialize R2 cloud storage
	err = helpers.InitR2Client()
	if err != nil {
		log.Printf("⚠️ R2 storage not configured: %v (file uploads will fail)", err)
	} else {
		log.Println("✅ R2 cloud storage initialized good")
	}

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

	// Set max multipart memory to 3GB for large video uploads
	router.MaxMultipartMemory = 3 << 30 // 3GB

	// Serve static files (uploaded images)
	router.Static("/uploads", "./uploads")

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes group
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
