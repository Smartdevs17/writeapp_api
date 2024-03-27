package main

import (
	"fmt"
	"writeapp_api/controllers"
	"writeapp_api/initializers"
	"writeapp_api/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	fmt.Println("Server stated running successfully")
	r := gin.Default()
	r.POST("/api/auth/register", controllers.Register)
	r.POST("/api/auth/login", controllers.Login)
	r.GET("/api/auth/validate", middleware.AuthMiddleware, controllers.ValidateAuth)
	r.POST("/api/auth/resetpassword", controllers.ResetPassword)
	r.POST("/api/documents", middleware.AuthMiddleware, controllers.CreateDocument)
	r.GET("/api/documents", middleware.AuthMiddleware, controllers.GetDocuments)
	r.GET("/api/documents/:id", middleware.AuthMiddleware, controllers.GetDocument)
	r.GET("/api/documents/user/:id", middleware.AuthMiddleware, controllers.GetUserDocuments)
	r.PUT("/api/documents/:id", middleware.AuthMiddleware, controllers.UpdateDocument)
	r.DELETE("/api/documents/:id", middleware.AuthMiddleware, controllers.DeleteDocument)
	r.GET("/api/documents/search/:query", middleware.AuthMiddleware, controllers.SearchDocuments)

	r.Run() // listen and serve on localhost:3000
}
