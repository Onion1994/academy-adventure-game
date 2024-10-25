package model

type Command interface {
	Execute(input PlayerInput, game *Game) string
}

type LookCommand struct{}

func (l LookCommand) Execute(input PlayerInput, game *Game) string {
	return game.player.ShowRoom(ConsoleDisplay{})
}

type ExitCommand struct{}

func (e ExitCommand) Execute(input PlayerInput, game *Game) string { return "" }

type CommandsCommand struct{}

func ShowCommands(d Display) string {
	return d.Show("-exit -> quits the game\n\n-commands -> shows the commands\n\n-look -> shows the content of the room.\n\n-approach <entity> -> to approach an entity\n\n-leave -> to leave an entity\n\n-inventory -> shows items in the inventory\n\n-take <item> -> to take an item into your inventory\n\n-drop <item> -> to drop an item from your inventory and move it to the current room\n\n-use <item> -> to make use of a certain item when you approach an entity\n\n-move <direction> -> to move to a different room\n\n-map -> shows the directions you can take\n")
}

func (c CommandsCommand) Execute(input PlayerInput, game *Game) string {

	return ShowCommands(ConsoleDisplay{})
}

type TakeCommand struct{}

func (t TakeCommand) Execute(input PlayerInput, game *Game) string {

	if len(input.Args) > 0 {
		return game.player.Take(input.Args[0], ConsoleDisplay{})
	} else {
		return "Specify an item to take."
	}
}

type DropCommand struct{}

func (d DropCommand) Execute(input PlayerInput, game *Game) string {

	if len(input.Args) > 0 {
		return game.player.Drop(input.Args[0], ConsoleDisplay{})
	} else {
		return "Specify an item to drop."
	}
}

type InventoryCommand struct{}

func (i InventoryCommand) Execute(input PlayerInput, game *Game) string {

	return game.player.ShowInventory(ConsoleDisplay{})
}

type ApproachCommand struct{}

func (a ApproachCommand) Execute(input PlayerInput, game *Game) string {

	if len(input.Args) > 0 {

		returnValue := game.player.Approach(input.Args[0], ConsoleDisplay{})

		if !game.unlockComputer.Triggered {
			if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "computer" {
				game.isAttemptingPassword = true
			}
		}
		if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "terminal" {
			game.isAttemptingTerminal = true
		}

		return returnValue

	} else {
		return "Specify an entity to approach."
	}
}

type UseCommand struct{}

func (u UseCommand) Execute(input PlayerInput, game *Game) string {

	if len(input.Args) > 0 {
		if game.player.CurrentEntity == nil {
			return game.player.Use(input.Args[0], "unspecified_entity", ConsoleDisplay{})
		} else {
			return game.player.Use(input.Args[0], game.player.CurrentEntity.Name, ConsoleDisplay{})
		}
	} else {
		return "Specify an item to use."
	}
}

type LeaveCommand struct{}

func (l LeaveCommand) Execute(input PlayerInput, game *Game) string {
	return game.player.Leave()
}

type MoveCommand struct{}

func (m MoveCommand) Execute(input PlayerInput, game *Game) string {

	if _, ok := game.player.Inventory["lanyard"]; ok {
		if len(input.Args) > 0 {
			return game.player.Move(input.Args[0], ConsoleDisplay{})
		} else {
			return "Specify a direction to move (e.g., north)."
		}
	} else {
		return "Doors are shut for you if you don't have a lanyard."
	}
}

type MapCommand struct{}

func (m MapCommand) Execute(input PlayerInput, game *Game) string {

	return game.player.ShowMap(ConsoleDisplay{})
}
