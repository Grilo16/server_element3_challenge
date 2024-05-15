package user

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Grilo16/server_element3_challenge/user"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {

	firstName := "tom"
	lastName := "Britton"
	email := "tom@email.com"
	sub := "mock-sub"
	
	user := user.NewUser(firstName, lastName, email, sub)

	assert.Equal(t, firstName, user.FirstName)
	assert.Equal(t, lastName, user.LastName)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, sub, user.Sub)
}


func TestFromRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "first_name", "last_name", "email", "password", "date_of_birth"}

	expectedValues := []driver.Value{1, "John", "Doe", "john.doe@example.com"}

	mockRow := mock.NewRows(columns).AddRow(expectedValues...)

	mock.ExpectQuery("SELECT").WillReturnRows(mockRow)

	user, err := user.FromRow(db.QueryRow("SELECT * FROM users WHERE id = ?", 1))

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john.doe@example.com", user.Email)
}



func TestFromRows(t *testing.T) {
	var users []user.User
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "first_name", "last_name", "email"}

	expectedValues := [][]driver.Value{
		{1, "John", "Doe", "john.doe@example.com"},
		{2, "Jane", "Smith", "jane.smith@example.com"},
	}

	mockRows := mock.NewRows(columns)
	for _, values := range expectedValues {
		mockRows.AddRow(values...)
	}

	mock.ExpectQuery("SELECT").WillReturnRows(mockRows)

	mockQueryRows, err := db.Query("SELECT * FROM users")
	assert.NoError(t, err)

	for mockQueryRows.Next() {
		user, err := user.FromRows(mockQueryRows)
		assert.NoError(t, err)
		
		users = append(users, *user)
	}
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, len(expectedValues))
	for i, user := range users {
		assert.Equal(t, expectedValues[i][0], user.Id)
		assert.Equal(t, expectedValues[i][1], user.FirstName)
		assert.Equal(t, expectedValues[i][2], user.LastName)
		assert.Equal(t, expectedValues[i][3], user.Email)
	}
}