package env

import (
	"ciphera-api/types"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func load() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("env cannot be loaded.")
		}
	})
}

func get(key string) string {
	load()
	return os.Getenv(key)

}

func GetSystemPort() string {
	return get("PORT")
}

func GetDatabaseServerUri() string {
	return get("DATABASE_SERVER_URI")
}

func GetDatabaseName() string {
	return get("DATABASE_NAME")
}

func GetJWTPublicKey() string {
	return get("JWT_PUBLIC_KEY")
}

func GetJWTPrivateKey() string {
	return get("JWT_PRIVATE_KEY")
}

func GetStorageCreds() types.R2Config {
	var config types.R2Config
	config.AccessKeyID = get("STORAGE_ACCESS_KEY_ID")
	config.Region = get("STORAGE_REGION")
	config.Endpoint = get("STORAGE_ENDPOINT")
	config.SecretAccessKey = get("STORAGE_SECRET_ACCESS_KEY")
	return config
}
