package controllers

import (
	"net/http"
	"os"
	"time"
	"writeapp_api/initializers"
	"writeapp_api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	//Get the request body
	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name, Email and Password are required",
		})
		return
	}

	//HashPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//Add the new user to the database
	newUser := models.User{Name: body.Name, Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    newUser,
	})
}

func Login(c *gin.Context) {
	//Get the email and Password from the request body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//Email and Password is required
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and Password are required",
		})
		return
	}

	//get the user from the database
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	//Check if the user exists
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	//validate the user password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//generate for the user a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	user.Token = tokenString
	initializers.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User logged in successfully",
		"data":    user,
	})
}

func ValidateAuth(c *gin.Context) {
	//Get the user from the context
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User is authenticated",
		"data":    user,
	})
}

func ResetPassword(c *gin.Context) {
	//Get the data from the request body
	var body struct {
		Email       string
		NewPassword string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Email == "" || body.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and NewPassword are required",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}
	user.Password = string(hash)
	initializers.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successfully",
		"data":    nil,
	})
}
