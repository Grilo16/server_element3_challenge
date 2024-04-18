package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Grilo16/server_element3_challenge/database"
	"github.com/Grilo16/server_element3_challenge/middleware"
	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/Grilo16/server_element3_challenge/userfiles"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// connString := "server=localhost;database=element3_challenge"
	
	driver := "sqlserver"
	connString := "server=sqlcontainer,1433;database=master;user=SA;password=StrongP4ssword!"
	
	db, err := sql.Open(driver, connString)
	if err != nil {
		fmt.Println("Error connecting to db")
	}
	database.Initialize(db)
	ctx := context.Background()
	
	// Initialize Repositories
	userRepository := user.NewUserRepository(db, ctx)
	userFilesRepository := userfiles.NewUserFilesRepository(db, ctx)
	
	// Initialize Services
	userServices := user.NewUserService(userRepository)
	userFilesServices := userfiles.NewUserFilesService(userFilesRepository)
	
	// Initialize Controllers
	userControllers := user.NewUserController(userServices)
	userFilesControllers := userfiles.NewUserFilesController(userFilesServices, userServices)
	
	
	
	// Create Router Config
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true 
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	
	// Set up router
	router := gin.Default()
	router.Use(cors.New(config))
	authMiddleware := middleware.InitializeAuthMiddleware(userServices)
	
	// Create Private routes group
	privateRoutes := router.Group("auth")
	privateRoutes.Use(authMiddleware.MiddlewareFunc())
	
	// Setup authentication routes
	router.POST("login", authMiddleware.LoginHandler)	
	router.POST("refresh", authMiddleware.RefreshHandler)	
	
	// Initialize Routes
	userControllers.InitializeRoutes(router, privateRoutes)
	userFilesControllers.InitializeRoutes(router, privateRoutes)

	router.Run(":8080")
}
