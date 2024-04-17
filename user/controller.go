package user

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
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
	router.POST("users", uc.createNewUserHandler)

	// Private Routes
	privateRoutes.GET("users", uc.getAllUsersHandler)
	privateRoutes.GET("users/my-details", uc.getMyUserDetailsHandler)
	privateRoutes.DELETE("users", uc.deleteUserByIdHandler)
	privateRoutes.PATCH("users", uc.editUserByIdHandler)
}

func (uc *UserController) getMyUserDetailsHandler(ctx *gin.Context) {
	user, err := uc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
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

func (uc *UserController) getAllUsersHandler(ctx *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		fmt.Println("Error Fetching user", err)
	}
	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) deleteUserByIdHandler(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	user, err := uc.userService.GetUserByEmail(claims["identity"].(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	result := uc.userService.DeleteUserById(fmt.Sprint(user.Id))
	ctx.JSON(http.StatusOK, result)
}

func (uc *UserController) editUserByIdHandler(ctx *gin.Context) {
	user, err := uc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var editData map[string]interface{}
	if err := ctx.BindJSON(&editData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
	}
	
	editedUser, err := uc.userService.EditUserById(user.Id, editData)
	if err != nil {
		fmt.Println("Error creating user:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit user"})
        return
    }

	ctx.JSON(http.StatusOK, editedUser)
}

func (uc *UserController) createNewUserHandler(ctx *gin.Context) {
	var newUser User
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
        return
	}

	savedUser, err := uc.userService.CreateNewUser(&newUser)
	if err != nil {
		fmt.Println("Error creating user:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

	ctx.JSON(http.StatusOK, savedUser)
}


