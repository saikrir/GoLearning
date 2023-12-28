package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/saikrir/notesvc/internal/domain"
	go_ora "github.com/sijms/go-ora/v2"
)

func NewConnection() (*sql.DB, error) {
	user, pwd, err := getCredentials()

	if err != nil {
		log.Fatalln("Cannot establish connection, credentials lookup problem ", err)
	}

	urlOptions := map[string]string{
		"SID": "XE",
	}

	connStr := go_ora.BuildUrl("skrao-db-server", 1521, "", user, pwd, urlOptions)

	db, err := sql.Open("oracle", connStr)

	if err != nil {
		fmt.Println("Failed to connect to Database ", err)
		return nil, err
	}

	if db.Ping() != nil {
		log.Fatalf("Failed to ping database with connection string %s, error is %#v \n", connStr, err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)

	return db, nil
}

func GetAllNotes(db *sql.DB) ([]domain.NoteRecord, error) {
	getAllNoteSQL := `SELECT * FROM APP_USER.T_USER_NOTES ORDER BY ID`

	rows, err := db.Query(getAllNoteSQL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allRows := make([]domain.NoteRecord, 0)
	for rows.Next() {

		var newNoteRecord domain.NoteRecord

		if err := rows.Scan(&newNoteRecord.Id, &newNoteRecord.Description, &newNoteRecord.CreatedAt, &newNoteRecord.Status); err != nil {
			return nil, err
		}

		allRows = append(allRows, newNoteRecord)
	}

	return allRows, nil
}

func CreateNewNoteEntry(noteRecord domain.NoteRecord, db *sql.DB) (int64, error) {
	createNewNoteSQL := "INSERT INTO T_USER_NOTES (DESCRIPTION, CREATED_AT, STATUS) VALUES (:1, :2, :3)"

	fmt.Println("Will Insert ", noteRecord)
	result, err := db.Exec(createNewNoteSQL, noteRecord.Description, noteRecord.CreatedAt, noteRecord.Status)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func getCredentials() (string, string, error) {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")

	if len(user) == 0 || len(password) == 0 {
		return "", "", fmt.Errorf("failed to lookup database credentials from environment")
	}

	return user, password, nil
}
