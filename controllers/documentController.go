package controllers

import (
	"net/http"
	"strconv"
	"writeapp_api/initializers"
	"writeapp_api/models"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	if body.Title == "" || body.Content == "" || body.Author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Content and Author are required"})
		return
	}

	//Get the user from the context
	user, _ := c.Get("user")

	newDocument := models.Document{Title: body.Title, Content: body.Content, Author: body.Author, UserID: user.(models.User).ID}
	result := initializers.DB.Create(&newDocument).Preload("User")

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document created successfully", "data": newDocument})
}

func GetDocuments(c *gin.Context) {
	var documents []models.Document

	initializers.DB.Find(&documents)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Documents fetched successfully", "data": documents})
}

func GetDocument(c *gin.Context) {
	var document models.Document

	if initializers.DB.First(&document, c.Param("id")).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document fetched successfully", "data": document})
}

func UpdateDocument(c *gin.Context) {
	var body struct {
		Title   string
		Content string
		Author  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var document models.Document

	if initializers.DB.First(&document, c.Param("id")).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document not found"})
		return
	}

	initializers.DB.Model(&document).Updates(models.Document{Title: body.Title, Content: body.Content, Author: body.Author})

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document updated successfully", "data": document})
}

func DeleteDocument(c *gin.Context) {
	var document models.Document

	if initializers.DB.First(&document, c.Param("id")).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document not found"})
		return
	}

	initializers.DB.Delete(&document)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document deleted successfully"})
}

func GetDocumentsByAuthor(c *gin.Context) {
	var documents []models.Document

	initializers.DB.Where("author = ?", c.Param("author")).Find(&documents)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Documents fetched successfully", "data": documents})
}

func GetUserDocuments(c *gin.Context) {
	var documents []models.Document
	user_id := c.Param("id")
	user, _ := c.Get("user")
	num64, _ := strconv.ParseUint(user_id, 10, 32)
	num := uint(num64)
	if num != user.(models.User).ID {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	initializers.DB.Where("user_id = ?", user.(models.User).ID).Preload("User").Find(&documents)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Documents fetched successfully", "data": documents})
}

func SearchDocuments(c *gin.Context) {
	var documents []models.Document

	initializers.DB.Where("title LIKE ?", "%"+c.Query("title")+"%").Find(&documents)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Documents fetched successfully", "data": documents})
}
