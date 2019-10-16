package main

import (
	"net/http"
	"os"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/labstack/echo"
)

type score struct {
	Name  string `form:"name"`
	Score int    `form:"score"`
}

func main() {
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	e.GET("/api/rank", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/api/rank", func(c echo.Context) error {
		s := new(score)
		if err = c.Bind(s); err != nil {
			e.Logger.Error("bind failed")
		}
		fmt.Printf("%v", s)
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/", "public/")
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
