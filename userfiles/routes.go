package userfiles

import (
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, authGroup *gin.RouterGroup) {
	UserFilesController := NewUserFilesController()

	authGroup.GET("files", UserFilesController.getAllUserFilesHanlder)
	authGroup.POST("files", UserFilesController.saveUserFileHander)
	authGroup.GET("files/:id", UserFilesController.downloadSavedFileHandler)
	authGroup.DELETE("files/:id", UserFilesController.deleteSavedFileHandler)
}
