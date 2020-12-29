package sweb

import (
	"os"
)

func getenv(k, v string) string {
	x := os.Getenv(k)
	if x == "" {
		return v
	}

	return x
}
