package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Grilo16/server_element3_challenge/database"
	"github.com/Grilo16/server_element3_challenge/keycloak"
	"github.com/Grilo16/server_element3_challenge/middleware"
	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/Grilo16/server_element3_challenge/userfiles"
	"github.com/Nerzal/gocloak/v13"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var client gocloak.GoCloak = *gocloak.NewClient("http://keycloak:8080")

const (
	realm        = "Element3"
	clientID     = "e3-challenge-server"
	clientSecret = "Q4xgljWZtLQGj50FTZqRC4rkEkBUfS0u"
)

type UserAuth struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func main() {

	driver := "sqlserver"
	// connString := "server=localhost;database=element3_challenge"
	connString := "server=sqlserver;database=e3-challenge;user=SA;password=StrongP4ssword!"

	// db, err := sql.Open(driver, connString)
	// if err != nil {
	// 	fmt.Println("Error connecting to db")
	// }

	var db *sql.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = sql.Open(driver, connString)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Println("Failed to connect to database, retrying in 5 seconds:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		fmt.Println("Could not connect to the database:", err)
	}

	fmt.Println("Connected to database.")

	database.Initialize(db)
	ctx := context.Background()

	// Initialize Repositories
	userRepository := user.NewUserRepository(db, ctx)
	userFilesRepository := userfiles.NewUserFilesRepository(db, ctx)

	// Initialize Services
	userServices := user.NewUserService(userRepository)
	userFilesServices := userfiles.NewUserFilesService(userFilesRepository)
	keyCloakServices := keycloak.NewKeycloakService()

	// Initialize Controllers
	userControllers := user.NewUserController(userServices)
	userFilesControllers := userfiles.NewUserFilesController(userFilesServices, userServices)
	keyCloakController := keycloak.NewKeycloakController(keyCloakServices)

	// Create Router Config
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	// Set up router
	router := gin.Default()
	router.Use(cors.New(config))
	authMiddleware := middleware.TokenAuthMiddleware(userServices)

	// Create Private routes group
	privateRoutes := router.Group("realms/Element3/auth")
	privateRoutes.Use(authMiddleware)
	privateRoutes.POST("login", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"login": "attempted"})
	})

	// // Setup authentication routes
	// router.POST("login", authMiddleware.LoginHandler)
	// router.POST("refresh", authMiddleware.RefreshHandler)

	// router.POST("login", func(ctx *gin.Context) {
	// 	var credentials UserAuth
	// 	if err := ctx.BindJSON(&credentials); err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	// 		return
	// 	}

	// 	params := gocloak.TokenOptions{
	// 		GrantType: &grantType, // Or "client_credentials" or another appropriate type
	// 		ClientID:  &clientID,
	// 		ClientSecret: &clientSecret,
	// 		Scope:    &scope,
	// 		Username: &credentials.Email,
	// 		Password: &credentials.Password,
	// 	}

	// 	token, err := client.GetToken(ctx, realm, params)
	// 	if err != nil {
	// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 			return
	// 	}
	// 	ctx.JSON(http.StatusOK, token)
	// })

	// Initialize Routes
	userControllers.InitializeRoutes(router, privateRoutes)
	userFilesControllers.InitializeRoutes(router, privateRoutes)
	keyCloakController.InitializeRoutes(router)

	router.Run(":8080")
}
