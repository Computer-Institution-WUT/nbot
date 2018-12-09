package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zzzlk123/nbot/games"
)

var h = horse.Game{}

func main() {
	e := echo.New()
	e.GET("/start", startGame)
	e.GET("/run", performRun)
	e.GET("/stop", performRun)
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func startGame(c echo.Context) error {
	ret := h.StartGame()
	return c.String(http.StatusOK, ret)
}

func performRun(c echo.Context) error {
	ret := h.Run()
	return c.String(http.StatusOK, ret)
}

func stopGame(c echo.Context) error {
	ret := h.FinishGame()
	return c.String(http.StatusOK, ret)
}
