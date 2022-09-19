package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
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

const (
	StatementSchedule = "schedule"
	StatementData     = "data"
)

var RootDirectory, _ = os.Getwd()

func MustReadFile(fileName string) []byte {
	content, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Sprintf("Error while reading %s. Error: %s", fileName, err))
	}
	return content
}

func MustUnmarshalYaml(path string, content []byte, decodedContent interface{}) {
	err := yaml.Unmarshal(content, decodedContent)
	if err != nil {
		panic(fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", path, err))
	}
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
