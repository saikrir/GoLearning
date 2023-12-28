package main

import (
	"fmt"
	"log"
	"time"

	"github.com/saikrir/notesvc/internal/database"
	"github.com/saikrir/notesvc/internal/domain"
)

func main() {
	fmt.Println("Hello World")

	conn, err := database.NewConnection()

	if err != nil {
		log.Fatalln("Failed to connect to DB ", err)
	}

	_, err = database.CreateNewNoteEntry(newNoteRecord(), conn)

	if err != nil {
		log.Fatalln("Failed to insert ", err)
	}

	allRows, err := database.GetAllNotes(conn)

	if err != nil {
		log.Fatalln("Failed to read rows from Database ", err)
	}

	for _, newEntry := range allRows {
		fmt.Println("New Record ", newEntry)
	}

	defer conn.Close()
}

func newNoteRecord() domain.NoteRecord {
	return domain.NoteRecord{
		Id:          0,
		Description: "I Love GO",
		CreatedAt:   time.Now(),
		Status:      "ACTIVE",
	}
}
