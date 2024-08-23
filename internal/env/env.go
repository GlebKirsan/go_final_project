package env

import "os"

func GetEnvOrDefault(key string, def string) string {
	if variable := os.Getenv(key); len(variable) != 0 {
		return variable
	}
	return def
}
