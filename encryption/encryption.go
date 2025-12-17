package encryption

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func GenerateChallenge() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func VerifySignature(
	publicKeyBase64 string,
	challengeBase64 string,
	signatureBase64 string,
) error {

	pubKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return err
	}

	challenge, err := base64.StdEncoding.DecodeString(challengeBase64)
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return err
	}

	if !ed25519.Verify(
		pubKey,
		challenge,
		signature,
	) {
		return errors.New("invalid signature")
	}

	return nil
}
