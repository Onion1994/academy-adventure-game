package model

import "fmt"

type Command interface {
	Execute(input PlayerInput, game *Game)
}


type LookCommand struct{}

func (l LookCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	g.player.ShowRoom(ConsoleDisplay{})
}

type ExitCommand struct{}

func (e ExitCommand) Execute(input PlayerInput, g *Game) {}

type CommandsCommand struct{}

func (c CommandsCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	ShowCommands(ConsoleDisplay{})
}

type TakeCommand struct{}

func (t TakeCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		g.player.Take(input.Args[0], ConsoleDisplay{})
	} else {
		fmt.Println("Specify an item to take.")
	}
}

type DropCommand struct{}

func (d DropCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		g.player.Drop(input.Args[0], ConsoleDisplay{})
	} else {
		fmt.Println("Specify an item to drop.")
	}
}

type InventoryCommand struct{}

func (i InventoryCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	g.player.ShowInventory(ConsoleDisplay{})
}

type ApproachCommand struct{}

func (a ApproachCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		g.player.Approach(input.Args[0], ConsoleDisplay{})

		if !g.unlockComputer.Triggered {
			if g.player.CurrentEntity != nil && g.player.CurrentEntity.Name == "computer" {
				g.isAttemptingPassword = true
			}
		}
		if g.player.CurrentEntity != nil && g.player.CurrentEntity.Name == "terminal" {
			g.isAttemptingTerminal = true
		}

	} else {
		fmt.Println("Specify an entity to approach.")
	}
}

type UseCommand struct{}

func (u UseCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	if len(input.Args) > 0 {
		if g.player.CurrentEntity == nil {
			g.player.Use(input.Args[0], "unspecified_entity", ConsoleDisplay{})
		} else {
			g.player.Use(input.Args[0], g.player.CurrentEntity.Name, ConsoleDisplay{})
		}
	} else {
		fmt.Println("Specify an item to use.")
	}
}

type LeaveCommand struct{}

func (l LeaveCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	g.player.Leave()
}

type MoveCommand struct{}

func (m MoveCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	if _, ok := g.player.Inventory["lanyard"]; ok {
		if len(input.Args) > 0 {
			g.player.Move(input.Args[0], ConsoleDisplay{})
		} else {
			fmt.Println("Specify a direction to move (e.g., north).")
		}
	} else {
		fmt.Println("Doors are shut for you if you don't have a lanyard.")
	}
}

type MapCommand struct{}

func (m MapCommand) Execute(input PlayerInput, g *Game) {
	clearScreen()
	g.player.ShowMap(ConsoleDisplay{})
}
