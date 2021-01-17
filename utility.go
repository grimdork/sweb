package sweb

import (
	"os"
)

// Getenv returns a value for key k, or value v if the environment variable is undefined.
func Getenv(k, v string) string {
	x := os.Getenv(k)
	if x == "" {
		return v
	}

	return x
}
