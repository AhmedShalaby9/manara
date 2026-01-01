package database

import (
	"fmt"
	"log"
	"manara/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbCharset := os.Getenv("DB_CHARSET")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbDatabase,
		dbCharset,
	)

	var gormLogger logger.Interface
	if os.Getenv("APP_DEBUG") == "true" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // ← Add this
	})
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Printf("✅ Successfully connected to database: %s", dbDatabase)
	return nil
}

func AutoMigrate() error {
	log.Println("🔄 Running database migrations...")

	err := DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.AcademicYear{},
		&models.Teacher{},
		&models.Student{},
		&models.Course{},
		&models.Chapter{},
		&models.Lesson{},
		&models.LessonFile{},
		&models.LessonVideo{},
		&models.TeacherCourse{},
	)

	if err != nil {
		log.Printf("❌ AutoMigrate error: %v", err)
		return err
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

//GO_ENV=local go run main.go -migrate
