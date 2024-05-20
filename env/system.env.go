package env

import "os"

func GetPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}
