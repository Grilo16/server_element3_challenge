package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *UserService
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) InitializeRoutes(router *gin.Engine, privateRoutes *gin.RouterGroup) {
	// Public Routes 
	router.POST("create-account", uc.createNewUserAccountHandler)

	// My Account routes
	privateRoutes.GET("my-account", uc.getMyUserDetailsHandler)
	privateRoutes.DELETE("my-account", uc.deleteMyAccountHandler)
	privateRoutes.PATCH("my-account", uc.editMyUserDetailsHandler)

	// Private Routes
	privateRoutes.GET("users", uc.getAllUsersHandler)
	privateRoutes.GET("users/my-details", uc.getMyUserDetailsHandler)
	privateRoutes.DELETE("users", uc.deleteMyAccountHandler)
	privateRoutes.PATCH("users", uc.editMyUserDetailsHandler)
}


func (uc *UserController) createNewUserAccountHandler(ctx *gin.Context) {
	var newUser User
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
        return
	}
	user, err := uc.userService.CreateNewUser(&newUser)
	if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) getMyUserDetailsHandler(ctx *gin.Context) {
	user, err := uc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) editMyUserDetailsHandler(ctx *gin.Context) {
	user, err := uc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	var editData map[string]interface{}
	if err := ctx.BindJSON(&editData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
	}
	
	editedUser, err := uc.userService.EditUserById(user.Id, editData)
	if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit user"})
        return
    }

	ctx.JSON(http.StatusOK, editedUser)
}

func (uc *UserController) deleteMyAccountHandler(ctx *gin.Context) {
	user, err := uc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	result := uc.userService.DeleteUserById(user.Id)
	ctx.JSON(http.StatusOK, result)
}

func (uc *UserController) getAllUsersHandler(ctx *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// func (uc *UserController) getUserByEmailHandler(ctx *gin.Context) {
// 	claims := jwt.ExtractClaims(ctx)
// 	userId := claims["identity"].(int)
// 	user, err := uc.userService.GetUserByEmail(userId)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, user)
// }


