package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func LoadEnvirolment() {
	projectDirName := "kala-be-upgrade" //edit
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + "/.env")
	if err != nil {
		panic(err)
	}
}
