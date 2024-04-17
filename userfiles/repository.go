package userfiles

import (
	"context"
	"database/sql"

	"github.com/Grilo16/server_element3_challenge/database"
)

type UserFilesRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewUserFilesRepository() *UserFilesRepository {
	return &UserFilesRepository{
		db:  database.DB,
		ctx: context.Background(),
	}
}

func (ur UserFilesRepository) GetFileByFileId(id int) (*UserFiles, error) {
	query := `
		SELECT * FROM user_files WHERE id = @id
	`
	row := ur.db.QueryRowContext(ur.ctx, query, sql.Named("id", id))
	userFile, err := FromRow(row)

	if err != nil {
		return nil, err
	}

	return userFile, nil
}

func (ur UserFilesRepository) GetFilesByUserId(userId int) ([]UserFiles, error) {
	var userFiles []UserFiles

	query := "SELECT * FROM user_files WHERE user_id = @userId"

	rows, err := ur.db.QueryContext(ur.ctx, query, sql.Named("userId", userId))

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		userFile, err := FromRows(rows)
		if err != nil {
			return nil, err
		}
		userFiles = append(userFiles, *userFile)
	}

	return userFiles, nil
}

func (ur UserFilesRepository) DeleteFileById(id int) string {
	query := `
		DELETE FROM user_files WHERE id = @id
	`
	_, err := ur.db.ExecContext(ur.ctx, query, sql.Named("id", id))
	if err != nil {
		return err.Error()
	}
	return "Succesfully deleted userFile"
}

func (ur UserFilesRepository) EditFile(userFile UserFiles) (*UserFiles, error) {
	query := `
		UPDATE user_files SET file_path = @filePath, file_name = @fileName OUTPUT INSERTED.* WHERE id = @id
	`
	row := ur.db.QueryRow(query, sql.Named("filePath", userFile.FilePath), sql.Named("fileName", userFile.FileName), sql.Named("id", userFile.Id))

	updatedUserFile, err := FromRow(row)
	if err != nil {
		return nil, err
	}

	return updatedUserFile, nil
}

func (ur UserFilesRepository) CreateNewFile(userFile UserFiles) (*UserFiles, error) {
	query := `
	INSERT INTO user_files (user_id, file_path, file_name) OUTPUT INSERTED.* VALUES (@userId, @filePath, @fileName)
`
	row := ur.db.QueryRow(query, sql.Named("userId", userFile.UserId), sql.Named("fileName", userFile.FileName), sql.Named("filePath", userFile.FilePath))

	newUserFile, err := FromRow(row)
	if err != nil {
		return nil, err
	}

	return newUserFile, nil
}
