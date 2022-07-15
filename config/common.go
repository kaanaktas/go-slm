package config

import (
	"log"
	"os"
	"runtime/debug"
	"time"
)

const NumberOfWorker = 5

//policy rule set directions
const (
	Request  = "request"
	Response = "response"
)

var RootDirectory, _ = os.Getwd()

func ReadFile(fileName string) ([]byte, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
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

// PolicyKey policy key to return rules from policy rule set
func PolicyKey(serviceName, direction string) string {
	return serviceName + "_" + direction
}

func IsModuleImported(currentModuleName string) bool {
	if currentModuleName == "" {
		currentModuleName = "github.com/kaanaktas/go-slm"
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Println("Failed to read build info")
		return false
	}

	return !(currentModuleName == bi.Path)
}
