package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCourses - Get all courses
func GetCourses(c *gin.Context) {
	var courses []models.Course

	res := database.DB.Preload("Chapters").Preload("Teachers.User").Find(&courses)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}

	helpers.Respond(c, true, courses, "Courses retrieved successfully")
}

// GetCourse - Get single course
func GetCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var course models.Course

	res := database.DB.Preload("Chapters.Lessons").Preload("Teachers.User").First(&course, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	helpers.Respond(c, true, course, "Course retrieved successfully")
}

// CreateCourse - Create a new course (Admin only)
func CreateCourse(c *gin.Context) {
	// Get text fields from form-data
	name := c.PostForm("name")
	description := c.PostForm("description")

	if name == "" {
		helpers.Respond(c, false, nil, "Name is required")
		return
	}

	// Handle optional image upload
	var imageURL string
	file, err := c.FormFile("image")
	if err == nil {
		// Image provided (optional)
		baseURL := os.Getenv("APP_URL")
		if baseURL == "" {
			baseURL = "http://localhost:" + os.Getenv("APP_PORT")
		}

		imageURL, err = helpers.UploadImage(file, "courses", baseURL)
		if err != nil {
			helpers.Respond(c, false, nil, err.Error())
			return
		}
	}

	course := models.Course{
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
	}

	if err := database.DB.Create(&course).Error; err != nil {
		// Delete uploaded image if course creation fails
		if imageURL != "" {
			helpers.DeleteImage(imageURL)
		}
		helpers.Respond(c, false, nil, "Failed to create course")
		return
	}

	helpers.Respond(c, true, course, "Course created successfully")
}

// UpdateCourse - Update course
func UpdateCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var course models.Course

	if err := database.DB.First(&course, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	// Get text fields from form-data
	name := c.PostForm("name")
	description := c.PostForm("description")

	// Update text fields
	if name != "" {
		course.Name = name
	}
	if description != "" {
		course.Description = description
	}

	// Handle optional image upload
	file, err := c.FormFile("image")
	if err == nil {
		// New image provided
		oldImageURL := course.ImageURL

		baseURL := os.Getenv("APP_URL")
		if baseURL == "" {
			baseURL = "http://localhost:" + os.Getenv("APP_PORT")
		}

		// Upload new image
		newImageURL, err := helpers.UploadImage(file, "courses", baseURL)
		if err != nil {
			helpers.Respond(c, false, nil, err.Error())
			return
		}

		course.ImageURL = newImageURL

		// Delete old image after successful upload
		if oldImageURL != "" {
			helpers.DeleteImage(oldImageURL)
		}
	}

	if err := database.DB.Save(&course).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to update course")
		return
	}

	helpers.Respond(c, true, course, "Course updated successfully")
}

// DeleteCourse - Delete course
func DeleteCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	// Delete image if exists
	if course.ImageURL != "" {
		helpers.DeleteImage(course.ImageURL)
	}

	if err := database.DB.Delete(&course).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete course")
		return
	}

	helpers.Respond(c, true, nil, "Course deleted successfully")
}

// UploadCourseImage - Upload or update course image
func UploadCourseImage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		helpers.Respond(c, false, nil, "No image file provided")
		return
	}

	// Save old image URL
	oldImageURL := course.ImageURL

	baseURL := os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + os.Getenv("APP_PORT")
	}

	// Upload new image
	imageURL, err := helpers.UploadImage(file, "courses", baseURL)
	if err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Update course
	course.ImageURL = imageURL
	if err := database.DB.Save(&course).Error; err != nil {
		// Delete newly uploaded image if database update fails
		helpers.DeleteImage(imageURL)
		helpers.Respond(c, false, nil, "Failed to update course")
		return
	}

	// Delete old image after successful update
	if oldImageURL != "" {
		helpers.DeleteImage(oldImageURL)
	}

	helpers.Respond(c, true, course, "Course image uploaded successfully")
}

// DeleteCourseImage - Delete course image
func DeleteCourseImage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	if course.ImageURL == "" {
		helpers.Respond(c, false, nil, "Course has no image")
		return
	}

	// Delete image file
	if err := helpers.DeleteImage(course.ImageURL); err != nil {
		helpers.Respond(c, false, nil, "Failed to delete image file")
		return
	}

	// Update course
	course.ImageURL = ""
	if err := database.DB.Save(&course).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to update course")
		return
	}

	helpers.Respond(c, true, nil, "Course image deleted successfully")
}

// AssignCourseToTeacher - Assign a course to a teacher
func AssignCourseToTeacher(c *gin.Context) {
	var req models.AssignCourseToTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Check if teacher exists
	var teacher models.Teacher
	if err := database.DB.First(&teacher, req.TeacherID).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	// Check if course exists
	var course models.Course
	if err := database.DB.First(&course, req.CourseID).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	// Check if already assigned
	var existing models.TeacherCourse
	err := database.DB.Where("teacher_id = ? AND course_id = ?", req.TeacherID, req.CourseID).First(&existing).Error
	if err == nil {
		helpers.Respond(c, false, nil, "Course already assigned to this teacher")
		return
	}

	// Assign course
	teacherCourse := models.TeacherCourse{
		TeacherID: req.TeacherID,
		CourseID:  req.CourseID,
	}

	if err := database.DB.Create(&teacherCourse).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to assign course")
		return
	}

	helpers.Respond(c, true, teacherCourse, "Course assigned successfully")
}

// GetTeacherCourses - Get courses for a specific teacher
func GetTeacherCourses(c *gin.Context) {
	teacherID, _ := strconv.Atoi(c.Param("teacher_id"))

	var teacher models.Teacher
	if err := database.DB.Preload("Courses.Chapters").First(&teacher, teacherID).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	helpers.Respond(c, true, teacher.Courses, "Teacher courses retrieved successfully")
}
