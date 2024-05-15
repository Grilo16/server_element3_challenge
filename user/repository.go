package user

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewUserRepository(db *sql.DB, ctx context.Context) *UserRepository {
	return &UserRepository{
		db:  db,
		ctx: ctx,
	}
}

func (ur UserRepository) GetUserBySub(sub string) (*User, error) {
	query := `
		SELECT * FROM users WHERE sub = @sub
	`
	row := ur.db.QueryRowContext(ur.ctx, query, sql.Named("sub", sub))
	user, err := FromRow(row)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (ur UserRepository) GetUserById(id int) (*User, error) {
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

	query := "SELECT * FROM users"

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

func (ur UserRepository) DeleteUserById(id int) string {
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
		UPDATE Users SET first_name = @firstName, last_name = @lastName, email = @email, WHERE id = @id
	`
	_, err := ur.db.Exec(query, sql.Named("firstName", user.FirstName), sql.Named("lastName", user.LastName), sql.Named("email", user.Email), sql.Named("id", user.Id))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur UserRepository) CreateNewUser(user User) (*User, error) {
	query := `
	INSERT INTO Users (first_name, last_name, email, sub) OUTPUT INSERTED.* VALUES (@firstName, @lastName, @email, @sub)
`
	row := ur.db.QueryRow(query, sql.Named("firstName", user.FirstName), sql.Named("lastName", user.LastName), sql.Named("email", user.Email), sql.Named("sub", user.Sub))

	newUser, err := FromRow(row)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
