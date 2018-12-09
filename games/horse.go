package horse

import (
	"fmt"
	"math/rand"
	"strings"
)

// Game status structure
type Game struct {
	gameOn     bool
	horses     map[int]*Horse
	finishLine int
}

// Horse individual structure
type Horse struct {
	position int
	alive    bool
	buff     float64
	// health   float64
	// attack   float64
	// defence  float64
}

// StartGame init game and set the game status on
func (g *Game) StartGame() string {
	if g.gameOn == false {
		g.gameOn = true
		g.initGame()
		return "200"
	}
	return "游戏已经开始"
}

func (g *Game) initGame() {
	g.finishLine = 1000
	g.horses = make(map[int]*Horse)
	for i := 1; i < 5; i++ {
		g.horses[i] = newHorse()
	}
}

func newHorse() *Horse {
	var h Horse
	h.position = 0
	h.alive = true
	h.buff = 0
	return &h
}

func getRandom(min, max int) int {
	return rand.Intn(max-min) + min
}

// GetGameDescribe will return the literal describtion of the current game
func (g *Game) GetGameDescribe() string {
	laneLength := 36
	sideLength := 20
	minPosition := 1000
	for _, v := range g.horses {
		if v.position < minPosition {
			minPosition = v.position
		}
	}
	a := ""
	for _, v := range g.horses {
		a += strings.Repeat("-", sideLength)
		leadRange := (v.position - minPosition)
		if leadRange >= laneLength {
			leadRange = laneLength
		}
		leadRange = laneLength - leadRange
		if v.alive {
			// a += fmt.Sprintf("\n%v 当前位置：%v  状态：活着  BUFF：%.1f\n", k, v.position, v.buff)
			a += fmt.Sprintf("\n" + strings.Repeat(" ", leadRange) + "🐴\n")
		} else {
			a += fmt.Sprintf("\n" + strings.Repeat(" ", leadRange) + "☠️\n")
		}
	}
	a += strings.Repeat("-", 20)
	a += "\n"
	return fmt.Sprintf("%v", a)
}

// Run will apply random event on each horse, thus make the game continue
// use the dispather to constantly call Run() until the game reaches end
func (g *Game) Run() string {
	if g.gameOn {
		event := ""
		for i := range g.horses {
			event += g.horses[i].applyEvent()
			if g.horses[i].position > g.finishLine {
				g.FinishGame()
				event += fmt.Sprintf("\n游戏结束，%v号🐴获得了胜利\n", i)
				// event += g.GetGameDescribe()
				return event
			}
		}
		event += "\n"
		event += g.GetGameDescribe()
		return event
	}
	return "游戏尚未开始"
}

// FinishGame finishes a current running game
func (g *Game) FinishGame() string {
	if g.gameOn {
		g.gameOn = false
		for i := 1; i < 5; i++ {
			g.horses[i].position = 0
			g.horses[i].alive = true
			g.horses[i].buff = 0
		}
		return "游戏结束"
	}
	return "游戏尚未开始"
}

func (h *Horse) applyEvent() string {
	h.position += getRandom(100, 120)
	return ""
}
