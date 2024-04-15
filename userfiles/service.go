package userfiles

type UserFilesService struct {
	userFilesRepository *UserFilesRepository
}

func NewUserFilesService() *UserFilesService {
	userFilesRepository := NewUserFilesRepository()
	return &UserFilesService{
		userFilesRepository: userFilesRepository,
	}
}

func (ufs *UserFilesService) SaveUserFile(userFile *UserFiles) (*UserFiles, error) {
	userFile, err := ufs.userFilesRepository.CreateNewFile(*userFile)
	if err != nil {
		return nil, err
	}
	return userFile, nil
}

func (ufs *UserFilesService) GetUserFileById(id string) (*UserFiles, error) {
	userFile, err := ufs.userFilesRepository.GetFileByFileId(id)
	if err != nil {
		return nil, err
	}
	return userFile, nil
}

func (ufs *UserFilesService) GetAllUserFilesByUserId(userId string) ([]UserFiles, error) {
	userFiles, err := ufs.userFilesRepository.GetFilesByUserId(userId)
	if err != nil {
		return nil, err
	}
	return userFiles, nil
}

func (ufs *UserFilesService) DeleteUserFileById(id string) string {
	result := ufs.userFilesRepository.DeleteFileById(id)
	return result
}

// func (ufs *UserFilesService) EditUserById(id string, updates map[string]interface{}) (*User, error) {

// 	user, err := ufs.userFilesRepository.GetUserById(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for key, value := range updates {
// 		switch key {
// 		case "firstName":
// 			user.FirstName = value.(string)
// 		case "lastName":
// 			user.LastName = value.(string)
// 		case "email":
// 			user.Email = value.(string)
// 		case "password":
// 			newPassword, ok := value.(string)
// 			if !ok {
// 				return nil, err
// 			}
// 			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 3)
// 			if err !=nil {
// 				return nil, err
// 			}
// 			user.Password = string(hashedPassword)
// 		case "dateOfBirth":
// 			dobString, ok := value.(string)
// 			if !ok {
// 				return nil, err
// 			}
// 			dob, err := time.Parse("2006-01-02", dobString)
// 			if err != nil {
// 				return nil, err
// 			}
			
// 		user.DateOfBirth = dob
// 		}
// 	}

// 	editedUser, err := ufs.userFilesRepository.EditUser(*user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return editedUser, nil
// }
