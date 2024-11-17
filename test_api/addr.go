package test_api

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var serverAddr = func() string {
	var (
		isValid = true
		host    string
		port    int
	)
	portE := os.Getenv("PORT")
	if portE != "" {
		var err error
		port, err = strconv.Atoi(portE)
		if err != nil {
			isValid = false
			log.Println("Failed to parse PORT: ", err)
		}
	} else {
		isValid = false
		log.Println("PORT is not set")
	}

	hostE := os.Getenv("HOST")
	if hostE != "" {
		host = hostE
	} else {
		isValid = false
		log.Println("HOST is not set")
	}

	if !isValid {
		log.Fatal("Failed to get server address")
	}

	return fmt.Sprintf("http://%s:%d", host, port)
}()
