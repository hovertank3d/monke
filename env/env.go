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
)
