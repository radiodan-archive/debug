package utils

import (
	"log"
	"os"
)

func Hostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Could not determine hostname")
	}

	return
}
