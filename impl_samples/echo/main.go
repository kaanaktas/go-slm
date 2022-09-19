package main

import (
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/executor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := newEchoEngine()

	// Routes
	e.GET("/test", testGet)
	e.POST("/test", testPost)

	log.Printf("starting server at :%s", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("error while starting server at :%s, %v", port, err)
	}
}

func testPost(c echo.Context) error {
	respInByte, _ := ioutil.ReadAll(c.Request().Body)
	respBody := string(respInByte)
	responseBeforeHook(respBody, c, extractServiceId)

	return c.JSON(http.StatusOK, respBody)
}

func testGet(c echo.Context) error {
	p1 := c.QueryParam("param1")
	executor.Apply(p1, "test", config.Request)
	return c.JSON(http.StatusOK, "no_match")
}

func newEchoEngine() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(customRecover)
	e.Use(requestDump(extractServiceId))

	return e
}
