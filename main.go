package main

import (
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/datafilter"
	"log"
	"runtime"
)

func main() {
	defer config.Elapsed("Execution")()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in Execute", r)
		}

		log.Println("All Channels were closed successfully. Number of goroutine:", runtime.NumGoroutine())
	}()

	serviceName := "test"

	data := "clear data with no match"
	//data := "admin' AND 1=1 --"
	//data := "http://testing.com/book.html?default=<script>alert(document.cookie)</script>"
	//data := "44044333322221111deded AND 1=1 --ede4444333322221111dededede44044333322221111dededede4442333322221111dededede"

	datafilter.Execute(data, serviceName)
}
