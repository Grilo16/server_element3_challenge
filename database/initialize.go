package database

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)


func Initialize(db *sql.DB) {
	createUserTable(db)
	createUserFilesTable(db)
}

func createUserTable(db *sql.DB) {
	query := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='Users')
    create table Users (
        id INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email VARCHAR(255),
		sub VARCHAR(255) UNIQUE,
    )
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("error: ", err)
	}
}

func createUserFilesTable(db *sql.DB) {
	query := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='user_files')
    create table user_files (
        id INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
		user_id INT NOT NULL, 
		file_path VARCHAR(255) UNIQUE,
		file_name VARCHAR(255),
		FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    )
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("error: ", err)
	}
}
