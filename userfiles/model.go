package userfiles

import "database/sql"

type UserFiles struct {
	Id       int    `json:"id,omitempty" db:"id"`
	UserId   int    `json:"userId,omitempty" db:"user_id"`
	FilePath string `json:"filePath,omitempty" db:"file_path"`
	FileName string `json:"fileName,omitempty" db:"file_name"`
}

func NewUserFiles(userId int, filePath string, fileName string) *UserFiles {
	return &UserFiles{
		UserId:   userId,
		FilePath: filePath,
		FileName: fileName,
	}
}

func FromRow(row *sql.Row) (*UserFiles, error) {
	var userFiles UserFiles
	err := row.Scan(&userFiles.Id, &userFiles.UserId, &userFiles.FilePath, &userFiles.FileName)
	if err != nil {
		return nil, err
	}
	return &userFiles, nil
}

func FromRows(rows *sql.Rows) (*UserFiles, error) {
	var userFiles UserFiles
	err := rows.Scan(&userFiles.Id, &userFiles.UserId, &userFiles.FilePath, &userFiles.FileName)
	if err != nil {
		return nil, err
	}
	return &userFiles, nil
}
