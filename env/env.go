package env

import (
	"fmt"
	"os"
)

var (
	DBUser = os.Getenv("DB_USER")
	DBName = os.Getenv("DB_NAME")
	DBPass = os.Getenv("DB_PASS")
	DBHost = os.Getenv("DB_HOST")
	DSL    = fmt.Sprintf("postgres://%s:%s@%s/%s", DBUser, DBPass, DBHost, DBName)

	MonkeHost   = os.Getenv("MONKE_HOST")
	MonkePlayer = (os.Getenv("MONKE_PLAYER") == "1")
)

func init() {
	if MonkeHost == "" {
		MonkeHost = "https://monke.oparysh.online"
	}
}
