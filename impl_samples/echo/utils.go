package main

import (
	"bytes"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/executor"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
)

type ServiceIdExtractor func(c echo.Context) string

func requestDump(s ServiceIdExtractor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			var reqBody []byte
			if c.Request().Body != nil {
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

			serviceName := s(c)
			executor.Apply(string(reqBody), serviceName, config.Request)

			return next(c)
		}
	}
}

func extractServiceId(c echo.Context) string {
	url := c.Request().URL
	return strings.Split(url.Path, "/")[1]
}

func responseBeforeHook(respBody string, c echo.Context, s ServiceIdExtractor) {
	executor.Apply(respBody, s(c), config.Response)
}

func customRecover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err := c.JSON(http.StatusInternalServerError, r)
				if err != nil {
					return
				}
			}
		}()

		return next(c)
	}
}
