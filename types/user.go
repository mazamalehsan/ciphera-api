package types

type User struct {
	ID           string
	IssuedTokens []string
	Username     string
	Groups       []string
	PublicKey    string
}
