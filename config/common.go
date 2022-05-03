package config

import (
	"log"
	"os"
	"time"
)

const NumberOfWorker = 5

func ReadFile(fileName string) ([]byte, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("error reading the file %q: %v\n", fileName, err)
		return nil, err
	}

	return content, nil
}

func Elapsed(msg string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", msg, time.Since(start))
	}
}
