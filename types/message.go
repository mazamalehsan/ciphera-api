package types

import (
	"time"
)

type Message struct {
	ID        string
	To        string
	FileId    string
	GroupId   string
	From      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
