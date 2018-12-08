package main

import (
	"nbot/games"
	"fmt"
)

func main() {
	g := horse.Game{}
	fmt.Println(g)
	g.StartGame()
	fmt.Println(g)
	fmt.Println(g.GetGameDescribe())
	g.Run()
	fmt.Println(g.GetGameDescribe())
}
