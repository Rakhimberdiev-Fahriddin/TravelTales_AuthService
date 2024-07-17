package api

import (
	_ "my_module/api/docs"
	"my_module/api/handler"
	"my_module/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @title User
// @version 1.0
// @description API Gateway of Authorazation
// @host localhost:8085
// BasePath: /
func Router(handler *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", handler.RegisterUser)
		auth.POST("/login", handler.LoginUser)
		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/logout", handler.LogOutUser)
	}

	userAuth := router.Group("/api/v1/auth")
	userAuth.Use(middleware.Check)
	{
		userAuth.POST("/reset-password", handler.ResetUserPassword)
	}
	return router
}
