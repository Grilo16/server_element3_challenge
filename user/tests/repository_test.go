package user

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetUserById(t *testing.T) {
	 // Mock DB setup
	 db, mock, err := sqlmock.New()
	 if err != nil {
		 t.Fatalf("Error creating mock database: %s", err)
	 }
	 defer db.Close()
 
	 // Mock user data
	 expectedUser := &user.User{
		 Id:          1,
		 FirstName:   "John",
		 LastName:    "Doe",
		 Email:       "john.doe@example.com",
		 Password:    "password",
		 DateOfBirth: time.Now(),
	 }
 
	 // Expectation: Mock the row returned by QueryRowContext
	 rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "date_of_birth"}).
		 AddRow(expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email,  expectedUser.Password, expectedUser.DateOfBirth)
 
	 // Mock the query and return the mocked rows
	 mock.ExpectQuery("^SELECT \\* FROM users WHERE id = @id$").
	 WithArgs(1).
	 WillReturnRows(rows)
	 
	 // Create the UserRepository instance with the mock DB
	 userRepository := user.NewUserRepository(db, context.Background())
 
	 // Call the function to be tested
	 retrievedUser, err := userRepository.GetUserById(1)
	 if err != nil {
		 t.Fatalf("Error retrieving user: %s", err)
	 }
 
	 // Assert that the retrieved user matches the expected user
	 if !reflect.DeepEqual(retrievedUser, expectedUser) {
		 t.Errorf("Retrieved user does not match expected user. Got: %+v, Expected: %+v", retrievedUser, expectedUser)
	 }
 
	 // Ensure all expectations are met
	 if err := mock.ExpectationsWereMet(); err != nil {
		 t.Errorf("Unfulfilled expectations: %s", err)
	 }
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	// Mock DB setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	ctx := context.Background()
	defer db.Close()
	// Mock data
	expectedUser := &user.User{
		Id:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		DateOfBirth: time.Now(),
		Password:    "password",
	}
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "date_of_birth"}).
	AddRow(expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email,  expectedUser.Password, expectedUser.DateOfBirth)

// Mock the query and return the mocked rows
	mock.ExpectQuery("^SELECT \\* FROM users WHERE email = @email$").
	WithArgs("john.doe@example.com").
	WillReturnRows(rows)
	// Expectation
	// Repository setup
	repo := user.NewUserRepository(db, ctx)



	// Test method call
	user, err := repo.GetUserByEmail(expectedUser.Email)

	// Assertion
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
}

func TestUserRepository_GetAllUsers(t *testing.T) {
		// Mock DB setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
	
		ctx := context.Background()
		defer db.Close()

	// Mock data
	expectedUsers := []user.User{
		{
			Id:          1,
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			DateOfBirth: time.Now(),
			Password:    "password",
		},
		{
			Id:          2,
			FirstName:   "Jane",
			LastName:    "Smith",
			Email:       "jane.smith@example.com",
			DateOfBirth: time.Now(),
			Password:    "password123",
		},
	}
	
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "date_of_birth"}).
	AddRow(expectedUsers[0].Id, expectedUsers[0].FirstName, expectedUsers[0].LastName, expectedUsers[0].Email,  expectedUsers[0].Password, expectedUsers[0].DateOfBirth).
	AddRow(expectedUsers[1].Id, expectedUsers[1].FirstName, expectedUsers[1].LastName, expectedUsers[1].Email,  expectedUsers[1].Password, expectedUsers[1].DateOfBirth)

	// Mock the query and return the mocked rows
	mock.ExpectQuery("SELECT \\* FROM users").
	WillReturnRows(rows)

	// Repository setup
	repo := user.NewUserRepository(db, ctx)

	// Test method call
	users, err := repo.GetAllUsers()

	// Assertion
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, expectedUsers, users)
}

func TestUserRepository_DeleteUserById(t *testing.T) {
	// Mock DB setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	ctx := context.Background()
	defer db.Close()

	// Mock data
	userID := 1

	// Expectation
	mock.ExpectExec("DELETE FROM Users WHERE id = @id").WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))

	// Repository setup
	repo := user.NewUserRepository(db, ctx)

	// Test method call
	msg := repo.DeleteUserById(userID)

	// Assertion
	assert.Equal(t, "Succesfully deleted user", msg)
}

func TestUserRepository_EditUser(t *testing.T) {
	// Mock DB setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	ctx := context.Background()
	defer db.Close()

	// Mock data
	mockUser := &user.User{
		Id:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		Password:    "password",
		DateOfBirth: time.Now(),
	}

	// "^SELECT \\* FROM users WHERE email = @email$"
	// Expectation
	mock.ExpectExec("^UPDATE Users SET first_name = @firstName, last_name = @lastName, email = @email, password = @password, date_of_birth = @dateOfBirth WHERE id = @id$").
    WithArgs(mockUser.FirstName, mockUser.LastName, mockUser.Email, mockUser.Password, mockUser.DateOfBirth, mockUser.Id).
    WillReturnResult(sqlmock.NewResult(0, 1))

	// Repository setup
	repo := user.NewUserRepository(db, ctx)



	// Test method call
	updatedUser, err := repo.EditUser(*mockUser)

	// Assertion
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, mockUser, updatedUser)
}

func TestUserRepository_CreateNewUser(t *testing.T) {
	// Mock DB setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	ctx := context.Background()
	defer db.Close()

	// Mock data
	mockUser := user.User{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		DateOfBirth: time.Now(),
		Password:    "password",
	}
	expectedUser := &user.User{
		Id:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		DateOfBirth: mockUser.DateOfBirth,
		Password:    "password",
	}

	mock.ExpectQuery("INSERT INTO Users (.+) VALUES (.+)").WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "date_of_birth", "password"}).
	AddRow(expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Password, expectedUser.DateOfBirth))
	
	repo := user.NewUserRepository(db, ctx)

	newUser, err := repo.CreateNewUser(mockUser)

	// Assertion
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, expectedUser, newUser)
}