package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLessons - Get all lessons (filter by chapter, teacher, or search by name)
func GetLessons(c *gin.Context) {
	chapterID := c.Query("chapter_id")
	teacherID := c.Query("teacher_id")
	search := c.Query("search")
	var lessons []models.Lesson

	params := helpers.GetPaginationParams(c)
	query := database.DB.Model(&models.Lesson{}).Preload("Chapter.Course").Preload("Teacher.User")

	if chapterID != "" {
		query = query.Where("chapter_id = ?", chapterID)
	}
	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query = query.Order("`order` ASC")
	pagination, err := helpers.Paginate(query, params, &lessons)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to retrieve lessons")
		return
	}

	helpers.RespondWithPagin(c, true, lessons, "Lessons retrieved successfully", pagination)
}

// GetLesson - Get single lesson
func GetLesson(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var lesson models.Lesson

	res := database.DB.Preload("Chapter.Course").Preload("Teacher.User").Preload("Files").Preload("Videos").First(&lesson, id)
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
	if err := database.DB.Preload("Files").Preload("Videos").First(&lesson, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Lesson not found")
		return
	}

	// Delete all files from R2
	for _, file := range lesson.Files {
		helpers.DeleteFromR2(file.FileURL)
	}

	// Delete all videos from R2
	for _, video := range lesson.Videos {
		helpers.DeleteFromR2(video.VideoURL)
	}

	// Delete file and video records from database
	database.DB.Where("lesson_id = ?", id).Delete(&models.LessonFile{})
	database.DB.Where("lesson_id = ?", id).Delete(&models.LessonVideo{})

	if err := database.DB.Delete(&lesson).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete lesson")
		return
	}

	helpers.Respond(c, true, nil, "Lesson deleted successfully")
}

// UploadLessonFiles - Upload multiple files to a lesson
func UploadLessonFiles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var lesson models.Lesson
	if err := database.DB.First(&lesson, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Lesson not found")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		helpers.Respond(c, false, nil, "No files provided")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		helpers.Respond(c, false, nil, "No files provided")
		return
	}

	// Get current max order
	var maxOrder int
	database.DB.Model(&models.LessonFile{}).Where("lesson_id = ?", id).Select("COALESCE(MAX(`order`), 0)").Scan(&maxOrder)

	var uploadedFiles []models.LessonFile
	for i, file := range files {
		fileURL, err := helpers.UploadFileToR2(file, "lessons/files")
		if err != nil {
			helpers.Respond(c, false, nil, "Failed to upload file: "+file.Filename+" - "+err.Error())
			return
		}

		lessonFile := models.LessonFile{
			LessonID: uint(id),
			FileURL:  fileURL,
			FileName: file.Filename,
			FileType: helpers.GetFileExtension(file.Filename),
			FileSize: file.Size,
			Order:    maxOrder + i + 1,
		}

		if err := database.DB.Create(&lessonFile).Error; err != nil {
			helpers.DeleteFromR2(fileURL)
			helpers.Respond(c, false, nil, "Failed to save file record")
			return
		}

		uploadedFiles = append(uploadedFiles, lessonFile)
	}

	helpers.Respond(c, true, uploadedFiles, "Files uploaded successfully")
}

// UploadLessonVideos - Upload multiple videos to a lesson
func UploadLessonVideos(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var lesson models.Lesson
	if err := database.DB.First(&lesson, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Lesson not found")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		helpers.Respond(c, false, nil, "No videos provided")
		return
	}

	videos := form.File["videos"]
	if len(videos) == 0 {
		helpers.Respond(c, false, nil, "No videos provided")
		return
	}

	// Get current max order
	var maxOrder int
	database.DB.Model(&models.LessonVideo{}).Where("lesson_id = ?", id).Select("COALESCE(MAX(`order`), 0)").Scan(&maxOrder)

	var uploadedVideos []models.LessonVideo
	for i, video := range videos {
		videoURL, err := helpers.UploadVideoToR2(video, "lessons/videos")
		if err != nil {
			helpers.Respond(c, false, nil, "Failed to upload video: "+video.Filename+" - "+err.Error())
			return
		}

		lessonVideo := models.LessonVideo{
			LessonID:  uint(id),
			VideoURL:  videoURL,
			VideoName: video.Filename,
			FileSize:  video.Size,
			Order:     maxOrder + i + 1,
		}

		if err := database.DB.Create(&lessonVideo).Error; err != nil {
			helpers.DeleteFromR2(videoURL)
			helpers.Respond(c, false, nil, "Failed to save video record")
			return
		}

		uploadedVideos = append(uploadedVideos, lessonVideo)
	}

	helpers.Respond(c, true, uploadedVideos, "Videos uploaded successfully")
}

// DeleteLessonFile - Delete a single file from a lesson
func DeleteLessonFile(c *gin.Context) {
	lessonID, _ := strconv.Atoi(c.Param("id"))
	fileID, _ := strconv.Atoi(c.Param("file_id"))

	var file models.LessonFile
	if err := database.DB.Where("id = ? AND lesson_id = ?", fileID, lessonID).First(&file).Error; err != nil {
		helpers.Respond(c, false, nil, "File not found")
		return
	}

	// Delete from R2
	if err := helpers.DeleteFromR2(file.FileURL); err != nil {
		helpers.Respond(c, false, nil, "Failed to delete file from storage")
		return
	}

	// Delete from database
	if err := database.DB.Delete(&file).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete file record")
		return
	}

	helpers.Respond(c, true, nil, "File deleted successfully")
}

// DeleteLessonVideo - Delete a single video from a lesson
func DeleteLessonVideo(c *gin.Context) {
	lessonID, _ := strconv.Atoi(c.Param("id"))
	videoID, _ := strconv.Atoi(c.Param("video_id"))

	var video models.LessonVideo
	if err := database.DB.Where("id = ? AND lesson_id = ?", videoID, lessonID).First(&video).Error; err != nil {
		helpers.Respond(c, false, nil, "Video not found")
		return
	}

	// Delete from R2
	if err := helpers.DeleteFromR2(video.VideoURL); err != nil {
		helpers.Respond(c, false, nil, "Failed to delete video from storage")
		return
	}

	// Delete from database
	if err := database.DB.Delete(&video).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete video record")
		return
	}

	helpers.Respond(c, true, nil, "Video deleted successfully")
}

// GetLessonFiles - Get all files for a lesson
func GetLessonFiles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var files []models.LessonFile
	if err := database.DB.Where("lesson_id = ?", id).Order("`order` ASC").Find(&files).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to retrieve files")
		return
	}

	helpers.Respond(c, true, files, "Files retrieved successfully")
}

// GetLessonVideos - Get all videos for a lesson
func GetLessonVideos(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var videos []models.LessonVideo
	if err := database.DB.Where("lesson_id = ?", id).Order("`order` ASC").Find(&videos).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to retrieve videos")
		return
	}

	helpers.Respond(c, true, videos, "Videos retrieved successfully")
}
