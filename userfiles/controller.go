package userfiles

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/gin-gonic/gin"
)

type UserFilesController struct {
	userFilesService *UserFilesService
	userService *user.UserService
}

func NewUserFilesController(userFilesService *UserFilesService, userService *user.UserService) *UserFilesController {
	return &UserFilesController{
		userFilesService: userFilesService,
		userService: userService,
	}
}

func (ufc *UserFilesController) InitializeRoutes(router *gin.Engine, privateRoutes *gin.RouterGroup) {
	privateRoutes.GET("files", ufc.getAllUserFilesHanlder)
	privateRoutes.POST("files", ufc.saveUserFileHander)
	privateRoutes.GET("files/:id", ufc.downloadSavedFileHandler)
	privateRoutes.DELETE("files/:id", ufc.deleteSavedFileHandler)
}

func (ufc *UserFilesController) getAllUserFilesHanlder(ctx *gin.Context) {
	user, err := ufc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userFiles, err := ufc.userFilesService.GetAllUserFilesByUserId(user.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error retrieving logs " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userFiles)
}


func (ufc *UserFilesController) saveUserFileHander(ctx *gin.Context) {
	user, err := ufc.userService.GetAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving file -" + err.Error()})
		return
	}

	filePath := "./public/" + fmt.Sprint(user.Id) + "/" + file.Filename
	userFile := NewUserFiles(user.Id, filePath, file.Filename)
	savedUserFile, err := ufc.userFilesService.SaveUserFile(userFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "File already exists"})
		return
	}

	if err := ctx.SaveUploadedFile(file, savedUserFile.FilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, savedUserFile)
}

func (ufc *UserFilesController) downloadSavedFileHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}


	userFile, err := ufc.userFilesService.GetUserFileById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error retrieving log"})
		return
	}
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
    ctx.Header("Pragma", "no-cache")
    ctx.Header("Expires", "0")
	ctx.FileAttachment(userFile.FilePath, userFile.FileName)
}

func (ufc *UserFilesController) deleteSavedFileHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	userFile, err := ufc.userFilesService.GetUserFileById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error fetching userfile"})
		return
	}

	err = os.Remove(userFile.FilePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error deleting file"})
		return
	}
	result := ufc.userFilesService.DeleteUserFileById(id)
	ctx.JSON(http.StatusOK, result)

}



// func (ufc *UserFilesController) getUserByIdHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	user, err := ufc.userFilesService.GetUserById(id)

// 	if err != nil {
// 		fmt.Println("Error Fetching user", err)
// 	}
// 	ctx.JSON(http.StatusOK, user)
// }

// func (ufc *UserFilesController) getAllUsersHandler(ctx *gin.Context) {
// 	users, err := ufc.userFilesService.GetAllUsers()
// 	if err != nil {
// 		fmt.Println("Error Fetching user", err)
// 	}
// 	ctx.JSON(http.StatusOK, users)
// }

// func (ufc *UserFilesController) deleteUserByIdHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	result := ufc.userFilesService.DeleteUserById(id)
// 	ctx.JSON(http.StatusOK, result)
// }

// func (ufc *UserFilesController) editUserByIdHandler(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	var editData map[string]interface{}
// 	if err := ctx.BindJSON(&editData); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//         return
// 	}

// 	editedUser, err := ufc.userFilesService.EditUserById(id, editData)
// 	if err != nil {
// 		fmt.Println("Error creating user:", err)
//         ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit user"})
//         return
//     }

// 	ctx.JSON(http.StatusOK, editedUser)
// }

// func (ufc *UserFilesController) createNewUserHandler(ctx *gin.Context) {
// 	var newUser User
// 	if err := ctx.BindJSON(&newUser); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
//         return
// 	}

// 	savedUser, err := ufc.userFilesService.CreateNewUser(&newUser)
// 	if err != nil {
// 		fmt.Println("Error creating user:", err)
//         ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
//         return
//     }

// 	ctx.JSON(http.StatusOK, savedUser)
// }
