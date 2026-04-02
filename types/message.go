package types

import (
	"time"
)

type MessageKey struct {
	UserId       string
	EncryptedKey string
	KeyIV        string
}

type Message struct {
	UUID               string
	From               string
	To                 string       // for DM
	GroupId            string       // for group messages
	Content            string       // encrypted ciphertext (base64)
	IV                 string       // AES-GCM IV (base64)
	Keys               []MessageKey // per-member encrypted keys (group only)
	FileId             string
	SenderPublicKey    string // Ed25519 (for verification)
	SenderEncryptionKey string // X25519 (for ECDH decryption)
	FromUsername       string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
