package user

import (
	"context"
	"database/sql"

	"github.com/Grilo16/server_element3_challenge/database"
)

type UserRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db:  database.DB,
		ctx: context.Background(),
	}
}

func (ur UserRepository) GetUserById(id string) (*User, error) {
	query := `
		SELECT * FROM users WHERE id = @id
	`
	row := ur.db.QueryRowContext(ur.ctx, query, sql.Named("id", id))
	user, err := FromRow(row)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (ur UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT * FROM users WHERE email = @email
	`
	row := ur.db.QueryRowContext(ur.ctx, query, sql.Named("email", email))
	user, err := FromRow(row)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur UserRepository) GetAllUsers() ([]User, error) {
	var users []User

	query := "select * from users"

	rows, err := ur.db.QueryContext(ur.ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user, err := FromRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (ur UserRepository) DeleteUserById(id string) string {
	query := `
		DELETE FROM Users WHERE id = @id
	`
	_, err := ur.db.ExecContext(ur.ctx, query, sql.Named("id", id))
	if err != nil {
		return err.Error()
	}
	return "Succesfully deleted user"
}

func (ur UserRepository) EditUser(user User) (*User, error) {
	query := `
		UPDATE Users SET first_name = @firstName, last_name = @lastName, email = @email, password = @password, date_of_birth = @dateOfBirth WHERE id = @id
	`
	_, err := ur.db.Exec(query, sql.Named("firstName", user.FirstName), sql.Named("lastName", user.LastName), sql.Named("email", user.Email), sql.Named("dateOfBirth", user.DateOfBirth), sql.Named("id", user.Id), sql.Named("password", user.Password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur UserRepository) CreateNewUser(user User) (*User, error) {
	query := `
	INSERT INTO Users (first_name, last_name, email, password, date_of_birth) OUTPUT INSERTED.* VALUES (@firstName, @lastName, @email, @password, @dateOfBirth)
`
	row := ur.db.QueryRow(query, sql.Named("firstName", user.FirstName), sql.Named("lastName", user.LastName), sql.Named("email", user.Email), sql.Named("password", user.Password), sql.Named("dateOfBirth", user.DateOfBirth))

	newUser, err := FromRow(row)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
