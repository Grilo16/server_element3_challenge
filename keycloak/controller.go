package keycloak

import (
		"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

type KeyCloakController struct {
	keycloakService *KeycloakService
}

type UserCreationRequest struct {
    gocloak.User
    Password string `json:"Password"`
}

type StageBulkCreateOutput struct {
	gocloak.User
	Exists bool `json:"exists"`
}


func NewKeycloakController(identityManager *KeycloakService) *KeyCloakController {
	return &KeyCloakController{
		keycloakService: identityManager,
	}
}

func (kc *KeyCloakController) InitializeRoutes(router *gin.Engine) {
	router.POST("create-new-account", kc.CreateNewUserHandler)
	router.POST("bulk-create-stage", kc.StageBulkCreaeUserHandler)
	router.GET("create-realm", kc.CreateNewRealmHandler)
	router.GET("create-client", kc.CreateNewClientHandler)

}



func (kc *KeyCloakController) CreateNewClientHandler(ctx *gin.Context) {
	
	realmName := "test-realm"
	clientName := "test-client"
	client, err := kc.keycloakService.CreateNewClient(ctx, realmName, clientName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	ctx.JSON(http.StatusOK, client)

}

func (kc *KeyCloakController) CreateNewRealmHandler(ctx *gin.Context) {
	
	realmName := "test-realm"
	realm, err := kc.keycloakService.CreateNewRealm(ctx, realmName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	ctx.JSON(http.StatusOK, realm)

}

func (kc *KeyCloakController) CreateNewUserHandler(ctx *gin.Context) {
	var user gocloak.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := kc.keycloakService.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, createdUser)
}

func (kc *KeyCloakController) StageBulkCreaeUserHandler(ctx *gin.Context) {
	var users []StageBulkCreateOutput
	if err := ctx.BindJSON(&users); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := kc.keycloakService.loginRestApiClient(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	

	for index := range users {
		user := &users[index]
		params := gocloak.GetUsersParams{
			Email: user.Email,
		}
		userExists, err := kc.keycloakService.client.GetUsers(ctx, token.AccessToken, kc.keycloakService.realm, params)
		if err != nil {
			user.Exists = false
			continue
		}
		if len(userExists) > 0 {
			user.Exists = true
		}
	}
	ctx.JSON(http.StatusOK, users)
}

func (kc *KeyCloakController) BulkCreateUserHandler(ctx *gin.Context) {
	
}
