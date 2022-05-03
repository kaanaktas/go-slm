package main

import (
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/datafilter"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.Println("Starting with number of worker", config.NumberOfWorker)
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func Benchmark(b *testing.B) {

	content, err := config.ReadFile("benchmark_load.json")
	if err != nil {
		panic(err)
	}

	data := string(content)
	serviceName := "test"

	for i := 0; i < b.N; i++ {
		datafilter.Execute(data, serviceName)
	}
}
