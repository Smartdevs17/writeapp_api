package controllers

import (
	"net/http"
	"strconv"
	"writeapp_api/initializers"
	"writeapp_api/models"
	"writeapp_api/utils"

	"github.com/gin-gonic/gin"
)

func CreateDocument(c *gin.Context) {
	var body struct {
		Title   string
		Content string
		Author  string
		Count   int
	}

	if c.Bind(&body) != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "bad request", "Failed to read body")
		return
	}

	if body.Title == "" || body.Content == "" || body.Author == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "bad request", "Title, Content and Author are required")
		return
	}

	user, _ := c.Get("user")

	newDocument := models.Document{Title: body.Title, Content: body.Content, Author: body.Author, Count: body.Count, UserID: user.(models.User).ID}
	result := initializers.DB.Create(&newDocument).Preload("User")

	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "error", "Failed to create document")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Document created successfully", newDocument)
}

func GetDocuments(c *gin.Context) {
	var documents []models.Document

	initializers.DB.Find(&documents)

	utils.SuccessResponse(c, http.StatusOK, "Documents fetched successfully", documents)
}

func GetDocument(c *gin.Context) {
	var document models.Document

	if initializers.DB.First(&document, c.Param("id")).Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Not Found", "Document not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Document fetched successfully", document)
}

func UpdateDocument(c *gin.Context) {
	var body struct {
		Title   string
		Content string
		Author  string
		Count   int
	}

	if c.Bind(&body) != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "bad request", "Failed to read body")
		return
	}

	var document models.Document

	if initializers.DB.First(&document, c.Param("id")).Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Not Found", "Document not found")
		return
	}

	initializers.DB.Model(&document).Updates(models.Document{Title: body.Title, Content: body.Content, Author: body.Author, Count: body.Count})
	utils.SuccessResponse(c, http.StatusOK, "Document updated successfully", document)
}

func DeleteDocument(c *gin.Context) {
	var document models.Document

	if initializers.DB.First(&document, c.Param("id")).Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Not Found", "Document not found")
		return
	}

	initializers.DB.Delete(&document)
	utils.SuccessResponse(c, http.StatusOK, "Document deleted successfully", nil)
}

func GetDocumentsByAuthor(c *gin.Context) {
	var documents []models.Document

	initializers.DB.Where("author = ?", c.Param("author")).Find(&documents)
	utils.SuccessResponse(c, http.StatusOK, "Documents fetched successfully", documents)
}

func GetUserDocuments(c *gin.Context) {
	var documents []models.Document
	user_id := c.Param("id")
	user, _ := c.Get("user")
	num64, _ := strconv.ParseUint(user_id, 10, 32)
	num := uint(num64)
	if num != user.(models.User).ID {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "You are not authorized to access this resource")
		return
	}

	initializers.DB.Where("user_id = ?", user.(models.User).ID).Preload("User").Find(&documents)
	utils.SuccessResponse(c, http.StatusOK, "Documents fetched successfully", documents)
}

func SearchDocuments(c *gin.Context) {
	// Get the authenticated user from the context
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "bad request", "User not found")
		return
	}

	// Extract user ID from the user object
	user_id := user.(models.User).ID

	// Query parameter for search filter
	search := c.Query("search")

	// Define a slice to hold the documents
	var documents []models.Document

	// Construct the base query with user_id
	query := initializers.DB.Where("user_id = ?", user_id)

	// Add conditions for title and/or content filters if provided
	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Execute the query
	if err := query.Find(&documents).Error; err != nil {
		// Handle database error
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch documents", err.Error())
		return
	}

	// Check if no documents were found
	if len(documents) == 0 {
		// Respond with an empty array
		utils.SuccessResponse(c, http.StatusOK, "No documents found", []models.Document{})
		return
	}

	// Respond with the documents found
	utils.SuccessResponse(c, http.StatusOK, "Documents fetched successfully", documents)
}
