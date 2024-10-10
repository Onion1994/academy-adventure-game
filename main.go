package main

import (
	"academy-adventure-game/model"
)

func main() {

	game := &model.Game{}
	game.SetupGame()
	game.RunGame()

}
