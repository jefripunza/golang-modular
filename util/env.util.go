package util

import (
	"fmt"
	"path/filepath"
	"project/env"

	"github.com/joho/godotenv"
)

type Env struct{}

func (ref Env) Load() {
	pwd := env.GetPwd()
	envFilePath := filepath.Join(pwd, ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Println("file .env not found")
	}
}
