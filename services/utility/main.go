package utility

import "github.com/joho/godotenv"

func LoadEnvironmentVariables(path string) {

	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}
}
