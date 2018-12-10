package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/zzzlk123/nbot/games"
)

type safeHorseGame struct {
	game horse.Game
	mux  sync.Mutex
}

type betTable struct {
	betOn      bool
	scoreBoard map[string]int
	betPool    map[int]map[string]int
}

func (t *betTable) initTable() {
	t.scoreBoard = make(map[string]int)
	t.betPool = make(map[int]map[string]int)
	t.betOn = false
	for i := 1; i < 5; i++ {
		t.betPool[i] = make(map[string]int)
	}
}

func (t *betTable) changeScore(user string, score int) int {
	_, ok := t.scoreBoard[user]
	if !ok {
		fmt.Println("Register new user:", user)
		t.scoreBoard[user] = 1000
	}
	t.scoreBoard[user] += score
	return t.scoreBoard[user]
}

func (t *betTable) readScore() string {
	ret := ""
	for k, v := range t.scoreBoard {
		ret += fmt.Sprintf("%v\t%v\n", k, v)
	}
	return ret
}

func (t *betTable) placeBet(target int, user string, amount int) string {
	balance := t.changeScore(user, 0)
	if amount <= balance {
		t.betPool[target][user] += amount
		t.changeScore(user, amount*-1)
		return ""
	}
	t.changeScore(user, balance*-1)
	t.betPool[target][user] += balance
	return "余额不足，已 All in"
}

func (t *betTable) clearBet() {
	w := h.game.GetWinner()
	fmt.Println("Clearing bet... Winner is:", w)
	for k, v := range t.betPool {
		if k == w {
			for user, bet := range v {
				fmt.Println("YOU WIN, ", user)
				t.changeScore(user, bet*2)
			}
		}
		for user := range v {
			delete(v, user)
		}
	}
}

type safeBetTable struct {
	t   betTable
	mux sync.Mutex
}

var h = safeHorseGame{}
var bet = safeBetTable{}

func main() {
	e := echo.New()
	bet.t.initTable()
	e.GET("/start", startGame)
	e.GET("/run", performRun)
	e.GET("/start", startGame)
	e.POST("/bet", placeBet)
	e.GET("/clearbet", clearBet)
	e.GET("/rank", getRank)
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func placeBet(c echo.Context) error {
	bet.mux.Lock()
	defer bet.mux.Unlock()
	if bet.t.betOn {
		target := c.FormValue("target")
		targetInt, err := strconv.Atoi(target)
		if err != nil {
			return err
		}
		if targetInt < 1 || targetInt >= 5 {
			return c.String(http.StatusOK, "没有这匹马儿")
		}
		user := c.FormValue("user")
		amount := c.FormValue("amount")
		amountInt, err := strconv.Atoi(amount)
		if err != nil {
			return err
		}
		ret := bet.t.placeBet(targetInt, user, amountInt)
		return c.String(http.StatusOK, ret)
	}
	return c.String(http.StatusOK, "已封盘")
}

func clearBet(c echo.Context) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	bet.mux.Lock()
	defer bet.mux.Unlock()
	if !bet.t.betOn && !h.game.GetGameStatus() {
		bet.t.clearBet()
	}
	return c.String(http.StatusOK, "")
}

func getRank(c echo.Context) error {
	bet.mux.Lock()
	defer bet.mux.Unlock()
	ret := bet.t.readScore()
	return c.String(http.StatusOK, ret)
}

func startGame(c echo.Context) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	bet.mux.Lock()
	defer bet.mux.Unlock()
	bet.t.betOn = true
	fmt.Println("Starting the table...")
	ret := h.game.StartGame()
	return c.String(http.StatusOK, ret)
}

func performRun(c echo.Context) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	bet.mux.Lock()
	defer bet.mux.Unlock()
	if bet.t.betOn {
		bet.t.betOn = false
		fmt.Println("Closing table...")
		fmt.Println(bet.t.betPool)
	}
	ret := h.game.Run()
	return c.String(http.StatusOK, ret)
}

func stopGame(c echo.Context) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	ret := h.game.FinishGame(0)
	return c.String(http.StatusOK, ret)
}
