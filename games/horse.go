package horse

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Game status structure
type Game struct {
	gameOn     bool
	horses     map[int]*Horse
	finishLine int
	winner     int
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
	rand.Seed(time.Now().UnixNano())
	if g.gameOn == false {
		g.gameOn = true
		g.initGame()
		return "200"
	}
	return "游戏已经开始"
}

// GetWinner returns the winner
func (g *Game) GetWinner() int {
	return g.winner
}

// GetGameStatus returns the state of game
func (g *Game) GetGameStatus() bool {
	return g.gameOn
}

func (g *Game) initGame() {
	g.finishLine = 1500
	g.winner = 0
	if g.horses == nil {
		fmt.Println("Initializing horse...")
		g.horses = make(map[int]*Horse)
		for i := 1; i < 5; i++ {
			g.horses[i] = newHorse()
		}
	} else {
		for i := 1; i < 5; i++ {
			g.horses[i].alive = true
			g.horses[i].position = 0
			g.horses[i].buff = 1.0
		}
	}
}

func newHorse() *Horse {
	var h Horse
	h.position = 0
	h.alive = true
	h.buff = 1.0
	return &h
}

func getRandom(min, max int) int {
	return rand.Intn(max-min) + min
}

// Run will apply random event on each horse, thus make the game continue
// use the dispather to constantly call Run() until the game reaches end
func (g *Game) Run() string {
	if g.gameOn {
		horseEvent := make(map[int]string)
		laneLength := 36
		sideLength := 20
		minPosition := 1000
		// fmt.Println(minPosition)
		event := ""
		for i := range g.horses {
			eventStr := g.horses[i].applyEvent()
			switch eventStr {
			case "death":
				event += fmt.Sprintf("%v号🐴突然去世了\n", i)
				horseEvent[i] = ""
			case "revive":
				event += fmt.Sprintf("%v号🐴获得在👼姐姐的祝福下重生\n", i)
				horseEvent[i] = ""
			case "buff":
				event += fmt.Sprintf("%v号🐴放飞了自我\n", i)
				horseEvent[i] = "💨"
			case "buff2":
				event += fmt.Sprintf("%v号🐴开始冲刺\n", i)
				horseEvent[i] = "💨💨"
			case "buffgone":
				event += fmt.Sprintf("%v号🐴的加速效果已消失\n", i)
				horseEvent[i] = ""
			default:
				if g.horses[i].buff != 1.0 {
					horseEvent[i] = "💨"
				} else {
					horseEvent[i] = ""
				}
			}
			switch i {
			case 1:
				horseEvent[i] += "①"
			case 2:
				horseEvent[i] += "②"
			case 3:
				horseEvent[i] += "③"
			case 4:
				horseEvent[i] += "④"
			}
			if g.horses[i].position > g.finishLine {
				g.FinishGame(i)
				event += fmt.Sprintf("游戏结束，%v号🐴获得了胜利\n", i)
				// event += g.GetGameDescribe()
				return event
			}
		}
		for _, v := range g.horses {
			if v.position < minPosition && v.alive {
				minPosition = v.position
			}
		}
		for i := 1; i < 5; i++ {
			// fmt.Printf("%v horse position: %v, alive: %v, buff: %v\n", i, g.horses[i].position, g.horses[i].alive, g.horses[i].buff)
			if i != 1 && g.horses[i].alive {
				event += fmt.Sprintf("%s\n", strings.Repeat("-", sideLength))
			}
			leadRange := (g.horses[i].position - minPosition)
			if leadRange >= laneLength {
				leadRange = laneLength
			}
			leadRange = laneLength - leadRange
			if g.horses[i].alive {
				event += fmt.Sprintf("%s%s%s\n", strings.Repeat(" ", leadRange), "🐴", horseEvent[i])
			}
			//  else {
			// event += fmt.Sprintf("%s%s%s\n", strings.Repeat(" ", leadRange), "☠️", horseEvent[i])
			// }
		}
		return event
	}
	return "游戏尚未开始"
}

// FinishGame finishes a current running game
func (g *Game) FinishGame(winner int) string {
	if g.gameOn {
		g.gameOn = false
		g.winner = winner
		return "游戏结束"
	}
	return "游戏尚未开始"
}

func (h *Horse) applyEvent() string {
	ret := ""

	deathDice := getRandom(0, 100)
	buffDice := getRandom(0, 100)
	if h.alive && deathDice < 5 {
		h.alive = false
		return "death"
	}
	if !h.alive && deathDice < 30 {
		h.alive = true
		h.position += 300
		return "revive"
	}
	if h.alive && h.buff == 1.0 {
		switch {
		case buffDice < 10:
			h.buff = 2.2
			ret += "buff2"
		case buffDice < 15:
			h.buff = 1.4
			ret += "buff"
		}
		h.position += int(float64(getRandom(100, 120)) * h.buff)
	} else if h.alive && h.buff != 1.0 {
		if buffDice < 70 {
			h.buff = 1.0
			ret += "buffgone"
		}
		h.position += int(float64(getRandom(60, 100)) * h.buff)
	}
	return ret
}
