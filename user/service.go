package user

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	claims := jwt.ExtractClaims(ctx)
	idFloat64 := claims["identity"].(float64)
	userId := int(idFloat64)
	user, err := us.GetUserById(userId)
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
		case "password":
			newPassword, ok := value.(string)
			if !ok {
				return nil, err
			}
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 3)
			if err !=nil {
				return nil, err
			}
			user.Password = string(hashedPassword)
		case "dateOfBirth":
			dobString, ok := value.(string)
			if !ok {
				return nil, err
			}
			dob, err := time.Parse("2006-01-02", dobString)
			if err != nil {
				return nil, err
			}
			
		user.DateOfBirth = dob
		}
	}

	editedUser, err := us.userRepository.EditUser(*user)
	if err != nil {
		return nil, err
	}

	return editedUser, nil
}

func (us *UserService) CreateNewUser(user *User) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 3)
	if err !=nil {
		return nil, err
	}
	
	
	user.Password = string(hashedPassword)
	newUser, err := us.userRepository.CreateNewUser(*user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
