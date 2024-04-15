package user

import (
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, authGroup *gin.RouterGroup ) {
	userController := NewUserController()

	router.POST("users", userController.createNewUserHandler)
	
	authGroup.GET("users", userController.getAllUsersHandler)
	authGroup.GET("users/my-details", userController.getUserByEmailHandler)
	authGroup.DELETE("users", userController.deleteUserByIdHandler)
	authGroup.PATCH("users", userController.editUserByIdHandler)
}
