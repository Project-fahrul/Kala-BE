package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func LoadEnvirolment() {
	projectDirName := "Upgrade" //edit
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	godotenv.Load(string(rootPath) + "/.env")
}
