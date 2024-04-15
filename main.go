package main

import (
	"github.com/Grilo16/server_element3_challenge/database"
	"github.com/Grilo16/server_element3_challenge/middleware"
	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/Grilo16/server_element3_challenge/userfiles"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.Initialize()
	
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true 
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	
	router := gin.Default()
	router.Use(cors.New(config))
	
	authMiddleware := middleware.InitializeAuthMiddleware()
	
	
	router.POST("login", authMiddleware.LoginHandler)	
	router.POST("refresh", authMiddleware.RefreshHandler)	
	
	auth := router.Group("auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	
	user.InitializeRoutes(router, auth)
	userfiles.InitializeRoutes(router, auth)

	router.Run(":8080")
}
