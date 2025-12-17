package types

import "time"

type User struct {
	UUID      string
	Username  string
	PublicKey string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserLoginChallenge struct {
	UserUUID  string
	UUID      string
	Challenge string
	ExpiresAt time.Time
	Used      bool
}
