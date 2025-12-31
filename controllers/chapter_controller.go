package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetChapters - Get all chapters (optionally filter by course, teacher, or search by name)
func GetChapters(c *gin.Context) {
	courseID := c.Query("course_id")
	teacherID := c.Query("teacher_id")
	search := c.Query("search")
	var chapters []models.Chapter

	query := database.DB.Preload("Course").Preload("Teacher.User")

	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	res := query.Order("`order` ASC").Find(&chapters)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}

	helpers.Respond(c, true, chapters, "Chapters retrieved successfully")
}

// GetChapter - Get single chapter
func GetChapter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var chapter models.Chapter

	res := database.DB.Preload("Course").Preload("Teacher.User").Preload("Lessons.Teacher.User").First(&chapter, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Chapter not found")
		return
	}

	helpers.Respond(c, true, chapter, "Chapter retrieved successfully")
}

// CreateChapter - Create a new chapter
func CreateChapter(c *gin.Context) {
	var req models.CreateChapterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Verify course exists
	var course models.Course
	if err := database.DB.First(&course, req.CourseID).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	// Verify teacher exists
	var teacher models.Teacher
	if err := database.DB.First(&teacher, req.TeacherID).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	// Set default order if not provided
	if req.Order == 0 {
		var maxOrder int
		database.DB.Model(&models.Chapter{}).Where("course_id = ?", req.CourseID).Select("COALESCE(MAX(`order`), 0)").Scan(&maxOrder)
		req.Order = maxOrder + 1
	}

	chapter := models.Chapter{
		CourseID:    req.CourseID,
		TeacherID:   req.TeacherID,
		Name:        req.Name,
		Order:       req.Order,
		Description: req.Description,
	}

	if err := database.DB.Create(&chapter).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to create chapter")
		return
	}

	database.DB.Preload("Course").Preload("Teacher.User").First(&chapter, chapter.ID)

	helpers.Respond(c, true, chapter, "Chapter created successfully")
}

// UpdateChapter - Update chapter
func UpdateChapter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var chapter models.Chapter

	if err := database.DB.First(&chapter, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Chapter not found")
		return
	}

	var req models.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	if req.Name != "" {
		chapter.Name = req.Name
	}
	if req.Order != 0 {
		chapter.Order = req.Order
	}
	if req.Description != "" {
		chapter.Description = req.Description
	}

	if err := database.DB.Save(&chapter).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to update chapter")
		return
	}

	database.DB.Preload("Course").Preload("Teacher.User").First(&chapter, chapter.ID)

	helpers.Respond(c, true, chapter, "Chapter updated successfully")
}

// DeleteChapter - Delete chapter
func DeleteChapter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var chapter models.Chapter
	if err := database.DB.First(&chapter, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Chapter not found")
		return
	}

	if err := database.DB.Delete(&chapter).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete chapter")
		return
	}

	helpers.Respond(c, true, nil, "Chapter deleted successfully")
}
