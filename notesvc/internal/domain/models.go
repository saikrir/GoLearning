package domain

import (
	"fmt"
	"time"
)

type NoteRecord struct {
	Id          int64
	Description string
	CreatedAt   time.Time
	Status      string
}

func (n *NoteRecord) String() string {
	return fmt.Sprintf("ID %d, Description %s, CreatedAt %v, Status %s", n.Id, n.Description, n.CreatedAt, n.Status)
}
