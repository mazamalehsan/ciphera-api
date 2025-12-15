package types

import "time"

type Group struct {
	ID        string
	Name      string
	Admins    []string
	Users     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
