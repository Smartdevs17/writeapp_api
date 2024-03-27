package controllers

import (
	"net/http"
	"writeapp_api/initializers"
	"writeapp_api/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	initializers.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "all the users found successfully",
		"data":    users,
	})
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	initializers.DB.First(&user, id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user found successfully",
		"data":    user,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	var body struct {
		Role   string
		Status string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid body",
		})
		return
	}

	initializers.DB.First(&user, id)
	user.Role = body.Role
	user.Status = body.Status
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user updated successfully",
		"data":    user,
	})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	initializers.DB.First(&user, id)
	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user deleted successfully",
		"data":    nil,
	})
}
