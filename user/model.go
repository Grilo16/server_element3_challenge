package user

import (
	"database/sql"
)

type User struct {
	Id int `json:"id,omitempty" db:"id"`
	FirstName string `json:"firstName,omitempty" db:"first_name"`
	LastName string `json:"lastName,omitempty" db:"last_name"`
	Email string `json:"email,omitempty" db:"email"`
	Sub string `json:"sub,omitempty" db:"sub"`
}

func NewUser(firstName string, lastName string, email string, sub string) *User {
	return &User{
		FirstName: firstName, 
		LastName: lastName, 
		Email: email, 
		Sub: sub,
	}
}

func FromRow(row *sql.Row) (*User, error) {
	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Sub)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FromRows(rows *sql.Rows) (*User, error) {
	var user User
	err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Sub)
	if err != nil {
		return nil, err
	}
	return &user, nil
}