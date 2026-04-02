package types

import "time"

type User struct {
	UUID          string
	Username      string
	PublicKey     string // Ed25519 public key (auth/signing)
	EncryptionKey string // X25519 public key (ECDH encryption)
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserLoginChallenge struct {
	UserUUID  string
	UUID      string
	Challenge string
	ExpiresAt time.Time
	Used      bool
}
