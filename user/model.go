package user

import (
	"database/sql"
	"time"
)

type User struct {
	Id int `json:"id,omitempty" db:"id"`
	FirstName string `json:"firstName,omitempty" db:"first_name"`
	LastName string `json:"lastName,omitempty" db:"last_name"`
	Email string `json:"email,omitempty" db:"email"`
	DateOfBirth time.Time `json:"dateOfBirth,omitempty" db:"date_of_birth"`
	Password string `json:"password,omitempty" db:"password"`
}

func NewUser(firstName string, lastName string, email string, dateOfBirth time.Time, password string) *User {
	return &User{
		FirstName: firstName, 
		LastName: lastName, 
		Email: email, 
		DateOfBirth: dateOfBirth, 
		Password: password, 
	}
}

func FromRow(row *sql.Row) (*User, error) {
	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.DateOfBirth)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FromRows(rows *sql.Rows) (*User, error) {
	var user User
	err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.DateOfBirth)
	if err != nil {
		return nil, err
	}
	return &user, nil
}