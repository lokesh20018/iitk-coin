package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id       int
	username string
}

// SQL Statement for Create Table

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"hid" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"id" integer,
		"name" TEXT
	  );`
	log.Println("Creating student table...")
	// Prepare SQL Statement
	statement, err := db.Prepare(createStudentTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Execute SQL Statements
	statement.Exec()
	log.Println("student table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertStudent(db *sql.DB, code int, name string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(id, name) VALUES (?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var hid int
		var id int
		var name string

		row.Scan(&hid, &id, &name)
		log.Println("Student: ", "hid: ", hid, "  ", id, " ", name)
	}
}

func main() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating DB")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	createTable(sqliteDatabase)
	// Create Database Tables
	person := User{id: 190460}
	person.username = "Lokesh Bharati"
	// INSERT RECORDS

	insertStudent(sqliteDatabase, person.id, person.username)

	// DISPLAY INSERTED RECORDS
	displayStudents(sqliteDatabase)
}
