package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetChapters - Get all chapters (optionally filter by course, teacher, or search by name)
// For teachers: automatically filtered to their own chapters
// For students: filtered to their teacher's chapters AND their academic year
// For admins: can filter by teacher_id query param or see all
func GetChapters(c *gin.Context) {
	courseID := c.Query("course_id")
	academicYearID := c.Query("academic_year_id")

	search := c.Query("search")
	var chapters []models.Chapter

	params := helpers.GetPaginationParams(c)
	query := database.DB.Model(&models.Chapter{}).Preload("Course").Preload("Teacher.User").Preload("AcademicYear")

	// Role-based teacher scoping
	if teacherID := helpers.GetEffectiveTeacherID(c); teacherID != nil {
		query = query.Where("teacher_id = ?", *teacherID)
	}

	// For students: filter by their academic year
	roleValue, _ := c.Get("role_value")
	if roleValue == "student" {
		if academicYearID, exists := c.Get("academic_year_id"); exists {
			query = query.Where("academic_year_id = ?", academicYearID)
		}
	}

	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}

	if academicYearID != "" {
		query = query.Where("academic_year_id = ?", academicYearID)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query = query.Order("`order` ASC")
	pagination, err := helpers.Paginate(query, params, &chapters)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to retrieve chapters")
		return
	}

	helpers.RespondWithPagin(c, true, chapters, "Chapters retrieved successfully", pagination)
}

// GetChapter - Get single chapter
func GetChapter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var chapter models.Chapter

	res := database.DB.Preload("Course").Preload("Teacher.User").Preload("AcademicYear").Preload("Lessons.Teacher.User").First(&chapter, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Chapter not found")
		return
	}

	helpers.Respond(c, true, chapter, "Chapter retrieved successfully")
}

// CreateChapter - Create a new chapter
// For teachers: teacher_id and course_id are automatically set from token/teacher record
// For admins: teacher_id and course_id must be provided in request
func CreateChapter(c *gin.Context) {
	var req models.CreateChapterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Get effective teacher_id based on role
	teacherID, ok := helpers.GetTeacherIDForCreate(c, req.TeacherID)
	if !ok {
		helpers.Respond(c, false, nil, "Teacher ID is required")
		return
	}

	// Verify teacher exists and get their course_id
	var teacher models.Teacher
	if err := database.DB.First(&teacher, teacherID).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	// Determine course_id based on role
	var courseID uint
	if helpers.IsTeacher(c) {
		// Teachers use their assigned course
		if teacher.CourseID == nil {
			helpers.Respond(c, false, nil, "Teacher has no assigned course")
			return
		}
		courseID = *teacher.CourseID
	} else {
		// Admins must provide course_id
		if req.CourseID == 0 {
			helpers.Respond(c, false, nil, "Course ID is required")
			return
		}
		courseID = req.CourseID
	}

	// Verify course exists
	var course models.Course
	if err := database.DB.First(&course, courseID).Error; err != nil {
		helpers.Respond(c, false, nil, "Course not found")
		return
	}

	// Set default order if not provided
	if req.Order == 0 {
		var maxOrder int
		database.DB.Model(&models.Chapter{}).Where("course_id = ?", courseID).Select("COALESCE(MAX(`order`), 0)").Scan(&maxOrder)
		req.Order = maxOrder + 1
	}

	chapter := models.Chapter{
		CourseID:    courseID,
		TeacherID:   teacherID,
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
