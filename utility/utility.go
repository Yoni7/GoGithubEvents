package utility

import (
	"os"
    "regexp"
)

func GetEnv(key string, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func IsEmailValid(e string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(e)
}