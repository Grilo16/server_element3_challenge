package user

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}


func (us *UserService) GetAuthenticatedUser(ctx *gin.Context) (*User, error) {
	userSub, ok:= ctx.Get("sub")
	if !ok {
		return nil, errors.New("User sub not found")
	}
	user, err := us.GetUserBySub(userSub.(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (us *UserService) GetUserBySub(sub string) (*User, error) {
	user, err := us.userRepository.GetUserBySub(sub)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (us *UserService) GetUserById(id int) (*User, error) {
	user, err := us.userRepository.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	user, err := us.userRepository.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetAllUsers() ([]User, error) {
	users, err := us.userRepository.GetAllUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) DeleteUserById(id int) string {
	result := us.userRepository.DeleteUserById(id)
	return result
}

func (us *UserService) EditUserById(id int, updates map[string]interface{}) (*User, error) {
	
	user, err := us.userRepository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	for key, value := range updates {
		switch key {
		case "firstName":
			user.FirstName = value.(string)
		case "lastName":
			user.LastName = value.(string)
		case "email":
			user.Email = value.(string)
		}
	}

	editedUser, err := us.userRepository.EditUser(*user)
	if err != nil {
		return nil, err
	}

	return editedUser, nil
}

func (us *UserService) CreateNewUser(user *User) (*User, error) {
	newUser, err := us.userRepository.CreateNewUser(*user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
