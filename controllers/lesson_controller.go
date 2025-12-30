package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLessons - Get all lessons (filter by chapter or teacher)
func GetLessons(c *gin.Context) {
	chapterID := c.Query("chapter_id")
	teacherID := c.Query("teacher_id")
	var lessons []models.Lesson

	query := database.DB.Preload("Chapter").Preload("Teacher.User")

	if chapterID != "" {
		query = query.Where("chapter_id = ?", chapterID)
	}
	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}

	res := query.Order("`order` ASC").Find(&lessons)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}

	helpers.Respond(c, true, lessons, "Lessons retrieved successfully")
}

// GetLesson - Get single lesson
func GetLesson(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var lesson models.Lesson

	res := database.DB.Preload("Chapter.Course").Preload("Teacher.User").First(&lesson, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Lesson not found")
		return
	}

	helpers.Respond(c, true, lesson, "Lesson retrieved successfully")
}

// CreateLesson - Create a new lesson
func CreateLesson(c *gin.Context) {
	var req models.CreateLessonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Verify chapter exists
	var chapter models.Chapter
	if err := database.DB.First(&chapter, req.ChapterID).Error; err != nil {
		helpers.Respond(c, false, nil, "Chapter not found")
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
		database.DB.Model(&models.Lesson{}).Where("chapter_id = ?", req.ChapterID).Select("COALESCE(MAX(`order`), 0)").Scan(&maxOrder)
		req.Order = maxOrder + 1
	}

	lesson := models.Lesson{
		ChapterID:   req.ChapterID,
		TeacherID:   req.TeacherID,
		Name:        req.Name,
		Description: req.Description,
		Order:       req.Order,
	}

	if err := database.DB.Create(&lesson).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to create lesson")
		return
	}

	database.DB.Preload("Chapter.Course").Preload("Teacher.User").First(&lesson, lesson.ID)

	helpers.Respond(c, true, lesson, "Lesson created successfully")
}

// UpdateLesson - Update lesson
func UpdateLesson(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var lesson models.Lesson

	if err := database.DB.First(&lesson, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Lesson not found")
		return
	}

	var req models.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	if req.Name != "" {
		lesson.Name = req.Name
	}
	if req.Order != 0 {
		lesson.Order = req.Order
	}
	if req.Description != "" {
		lesson.Description = req.Description
	}

	if err := database.DB.Save(&lesson).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to update lesson")
		return
	}

	database.DB.Preload("Chapter.Course").Preload("Teacher.User").First(&lesson, lesson.ID)

	helpers.Respond(c, true, lesson, "Lesson updated successfully")
}

// DeleteLesson - Delete lesson
func DeleteLesson(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var lesson models.Lesson
	if err := database.DB.First(&lesson, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Lesson not found")
		return
	}

	if err := database.DB.Delete(&lesson).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete lesson")
		return
	}

	helpers.Respond(c, true, nil, "Lesson deleted successfully")
}
