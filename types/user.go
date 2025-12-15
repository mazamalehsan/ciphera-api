package types

import "time"

type User struct {
	UUID      string
	Username  string
	PublicKey string
	Privates  UserPrivates
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPrivates struct {
	Groups       []string
	IssuedTokens []string
}
