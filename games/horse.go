package horse

import (
	"fmt"
	"strings"
)

// Game status structure
type Game struct {
	gameOn     bool
	horses     map[string]*Horse
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
	return "500"
}

func (g *Game) initGame() {
	g.finishLine = 1000
	g.horses = make(map[string]*Horse)
	g.horses["小花"] = newHorse()
	g.horses["小明"] = newHorse()
}

func newHorse() *Horse {
	var h Horse
	h.position = 0
	h.alive = true
	h.buff = 0
	return &h
}

// GetGameDescribe will return the literal describtion of the current game
func (g *Game) GetGameDescribe() string {
	a := ""
	for k, v := range g.horses {
		a += strings.Repeat("=", 20)
		if v.alive {
			// a += fmt.Sprintf("\n%v 当前位置：%v  状态：活着  BUFF：%.1f\n", k, v.position, v.buff)
			a += fmt.Sprintf("\n" + strings.Repeat(" ", 19) + "🐎\n")
		} else {
			a += fmt.Sprintf("\n%v 当前位置：%v  状态：死了  BUFF：%.1f\n", k, v.position, v.buff)
		}
	}
	a += strings.Repeat("=", 20)
	return fmt.Sprintf("%v", a)
}

// Run will apply random event on each horse, thus make the game continue
// use the dispather to constantly call Run() until the game reaches end
func (g *Game) Run() string {
	g.horses["小花"].applyEvent()
	g.horses["小明"].applyEvent()
	return "1"
}

func (h *Horse) applyEvent() string {
	h.position += 100
	return "1"
}
