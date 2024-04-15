package database

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

var DB *sql.DB

func Initialize() {
	driver := "sqlserver"
	connString := "server=localhost;database=element3_challenge"

	var err error
	DB, err = sql.Open(driver, connString)
	if err != nil {
		fmt.Println("Error connecting to db")
	}

	createUserTable()
	createUserFilesTable()
}

func createUserTable() {
	query := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='Users')
    create table Users (
        id INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email VARCHAR(255),
		password VARCHAR(255),
		date_of_birth Date
    )
	`
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Println("error: ", err)
	}
}

func createUserFilesTable() {
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
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Println("error: ", err)
	}
}
