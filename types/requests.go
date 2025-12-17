package types

type LoginBody struct {
	Username        string
	SignedChallenge string
	LoginId         string
}
