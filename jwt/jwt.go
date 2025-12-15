package jwt

import (
	"ciphera-api/types"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(privateKey *rsa.PrivateKey, userId string, expiresAt time.Time) (string, error) {
	now := time.Now()

	claims := types.JWTClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func VerifyAccessToken(
	tokenStr string,
	publicKey *rsa.PublicKey,
) (*types.JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&types.JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodRS256 {
				return nil, errors.New("unexpected signing method")
			}
			return publicKey, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
