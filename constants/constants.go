package constants

func AllowedOrigins() []string {
	return []string{
		"http://localhost:3000",
		"https://privmail.com:2101",
		"http://localhost:3000/",
		"https://piehost.com",
		"https://privmail.com",
	}
}

func AllowedMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH"}
}

func AllowedHeaders() []string {
	return []string{"Authorization", "Content-Type"}
}

const (
	UserCollectionName    = "users"
	MessageCollectionName = "messages"
	GroupsCollectionName  = "groups"
)
