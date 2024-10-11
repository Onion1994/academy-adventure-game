package model

import "fmt"

type Command interface {
	Execute(input PlayerInput, game *Game)
}


type LookCommand struct{}

func (l LookCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	game.player.ShowRoom(ConsoleDisplay{})
}

type ExitCommand struct{}

func (e ExitCommand) Execute(input PlayerInput, game *Game) {}

type CommandsCommand struct{}

func (c CommandsCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	ShowCommands(ConsoleDisplay{})
}

type TakeCommand struct{}

func (t TakeCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		game.player.Take(input.Args[0], ConsoleDisplay{})
	} else {
		fmt.Println("Specify an item to take.")
	}
}

type DropCommand struct{}

func (d DropCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		game.player.Drop(input.Args[0], ConsoleDisplay{})
	} else {
		fmt.Println("Specify an item to drop.")
	}
}

type InventoryCommand struct{}

func (i InventoryCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	game.player.ShowInventory(ConsoleDisplay{})
}

type ApproachCommand struct{}

func (a ApproachCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		game.player.Approach(input.Args[0], ConsoleDisplay{})

		if !game.unlockComputer.Triggered {
			if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "computer" {
				game.isAttemptingPassword = true
			}
		}
		if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "terminal" {
			game.isAttemptingTerminal = true
		}

	} else {
		fmt.Println("Specify an entity to approach.")
	}
}

type UseCommand struct{}

func (u UseCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		if game.player.CurrentEntity == nil {
			game.player.Use(input.Args[0], "unspecified_entity", ConsoleDisplay{})
		} else {
			game.player.Use(input.Args[0], game.player.CurrentEntity.Name, ConsoleDisplay{})
		}
	} else {
		fmt.Println("Specify an item to use.")
	}
}

type LeaveCommand struct{}

func (l LeaveCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	game.player.Leave()
}

type MoveCommand struct{}

func (m MoveCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	if _, ok := game.player.Inventory["lanyard"]; ok {
		if len(input.Args) > 0 {
			game.player.Move(input.Args[0], ConsoleDisplay{})
		} else {
			fmt.Println("Specify a direction to move (e.g., north).")
		}
	} else {
		fmt.Println("Doors are shut for you if you don't have a lanyard.")
	}
}

type MapCommand struct{}

func (m MapCommand) Execute(input PlayerInput, game *Game) {
	clearScreen()
	game.player.ShowMap(ConsoleDisplay{})
}
