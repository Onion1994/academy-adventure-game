package model

import (
	"academy-adventure-game/global"
	"fmt"
)

type Game struct {
	player *Player
	// validInteractions         []*Interaction
	// gameOver                  bool
	introduction              string
	introductionShown         bool
	dishwasherChallengeWon    *Event
	unlockComputer            *Event
	remainingPasswordAttempts int
	computerPassword          string
	isAttemptingPassword      bool
	isAttemptingTerminal      bool
	IsFirstCommand            bool
	staffRoom                 *Room
	codingLab                 *Room
	terminalRoom              *Room
}

var kettleApproachedFirst = false
var sofaApproachedFirst = false
var deskApproachedFirst = false
var lanyardEventCompleted = false

var Commands = map[string]Command{
	"look":      LookCommand{},
	"exit":      ExitCommand{},
	"commands":  CommandsCommand{},
	"take":      TakeCommand{},
	"drop":      DropCommand{},
	"inventory": InventoryCommand{},
	"approach":  ApproachCommand{},
	"use":       UseCommand{},
	"leave":     LeaveCommand{},
	"move":      MoveCommand{},
	"map":       MapCommand{},
}

func executeCommand(input PlayerInput, game *Game) string {

	command := input.Command


	cmd, exists := Commands[command]

	if command == game.computerPassword {
		return game.unlockComputer.Outcome
	}

	if !exists {
		return fmt.Sprintf("Unknown command: %s", command)
	}
	return cmd.Execute(input, game)

}

func (game *Game) RunGame(playerInput PlayerInput) GameResponse {
	abandonedLanyard := game.staffRoom.Items["abandoned-lanyard"]
	tea := game.staffRoom.Items["tea"]
	lanyard := game.staffRoom.Items["lanyard"]
	firstPlate := game.codingLab.Items["first-plate"]
	secondPlate := game.codingLab.Items["second-plate"]
	thirdPlate := game.codingLab.Items["third-plate"]
	fourthPlate := game.codingLab.Items["fourth-plate"]
	fifthPlate := game.codingLab.Items["fifth-plate"]
	sixthPlate := game.codingLab.Items["sixth-plate"]

	sofa := game.staffRoom.Entities["sofa"]
	terminal := game.terminalRoom.Entities["terminal"]
	computer := game.codingLab.Entities["computer"]
	kettle := game.staffRoom.Entities["kettle"]
	dishwasher := game.staffRoom.Entities["dishwasher"]
	desk := game.codingLab.Entities["desk"]
	alan := game.codingLab.Entities["alan"]
	dan := game.terminalRoom.Entities["dan"]
	rosie := game.staffRoom.Entities["rosie"]

	var response GameResponse
	response.GameOver = false

	if !global.GameOver {
		fmt.Println(game.player.Inventory)
		if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "sofa" && !sofaApproachedFirst{
			abandonedLanyard.Hidden = false
			sofa.SetDescription("Your fellow academy student continues to sleep on the sofa. Something tells you it's down to you to get stuff done today...")
			sofaApproachedFirst = true
		}

		if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "kettle" && !kettleApproachedFirst {
			fmt.Println("This is before nil pointer")
			fmt.Println(game.player.Inventory)
			tea.Hidden = false
			kettle.SetDescription("A kettle — essential for survival, impossible to function without one nearby.")
			kettleApproachedFirst = true
		}

		if game.player.CurrentEntity != nil && game.player.CurrentEntity.Name == "desk" && !deskApproachedFirst {
			firstPlate.Hidden = false
			secondPlate.Hidden = false
			thirdPlate.Hidden = false
			fourthPlate.Hidden = false
			fifthPlate.Hidden = false
			sixthPlate.Hidden = false
			desk.SetDescription("Despite the disarray, it's clear this desk sees frequent use, with just enough space left to get work done.")
			deskApproachedFirst = true
		}

		for _, validInteraction := range ValidInteractions {
			if validInteraction.Event.Description == "get-your-lanyard" && validInteraction.Event.Triggered && !lanyardEventCompleted{
				lanyard.Hidden = false
				rosie.SetDescription("Can I help with anything else?")
				lanyardEventCompleted = true
			}
		}

		dishwasherLoaded := true

		// Check if all plates have been loaded
		plates := []*Interaction{
			ValidInteractions[1],
			ValidInteractions[2],
			ValidInteractions[3],
			ValidInteractions[4],
			ValidInteractions[5],
			ValidInteractions[6],
		}
	
		for _, plate := range plates {
			fmt.Println(plate.Event.Triggered)
			if !plate.Event.Triggered {
				dishwasherLoaded = false
				break
			}
		}
	
		if !game.dishwasherChallengeWon.Triggered {
			if dishwasherLoaded {
				game.player.TriggerEvent(game.dishwasherChallengeWon)
				alan.SetDescription("Ah, so you've managed to load the dishwasher! Splendid work — consider this challenge complete.\nI could have done it myself instead of writing that clever recursive function, but where's the fun in that?\nAfter all, they pay me for my intellect, not for doing the heavy lifting!\nBut I digress. You're free to proceed to the terminal room and speak with Dan for your final challenge.\nYou're doing an excellent job; keep it up!")
				dan.Hidden = false
				terminal.Hidden = false

				return GameResponse{
					Message: game.dishwasherChallengeWon.Outcome,
					GameOver: false,
				}
			}
		}

		if _, ok := game.player.Inventory["abandoned-lanyard"]; ok {
			global.GameOver = true
			response.GameOver = true
			return response
		}

		if global.GameOver {
			response.Message = "Thank you for playing!"
			return response
		}

		if playerInput.Command == "start" {
			if !game.introductionShown {
				response.Message = game.introduction
				game.introductionShown = true
				return response
			}
		}

		input := playerInput.Command

		if input == "exit" {
			response.Message = "Thank you for playing!"
			response.GameOver = true
			global.GameOver = true
			return response
		}

		if game.isAttemptingPassword {
			if game.remainingPasswordAttempts == 1 && input != game.computerPassword {
				response.Message = "Alan's computer is locked. Thank you for playing!"
				response.GameOver = true
				global.GameOver = true
				return response
			}
			if input == game.computerPassword {
				game.player.TriggerEvent(game.unlockComputer)
				computer.SetDescription("function completeTask(pile)\n   if pile == 0:\n      return 'Task Complete'\n   else:\n      completeTask(pile - 1)\n")
				alan.SetDescription("You've cracked the password! Impressive work...")
				game.isAttemptingPassword = false
				desk.Hidden = false
				dishwasher.Hidden = false
			} else if input == "leave" {
				game.isAttemptingPassword = false
			} else {
				game.remainingPasswordAttempts--
				response.Message = fmt.Sprintf("Incorrect password. Remaining attempts: %d", game.remainingPasswordAttempts)
				computer.SetDescription(fmt.Sprintf("Alan's computer. Remaining attempts: %d.\nEnter the password:", game.remainingPasswordAttempts))
				return response
			}
		}

		if game.isAttemptingTerminal {
			if input == "leave" {
				game.isAttemptingTerminal = false
				game.player.Leave()
				return response
			}

			if !game.IsFirstCommand {
				if input == "cd /secret-files" {
					response.Message = "The terminal displays:\n\n/secret-files/\n\nEnter the final command to win the game!"
					game.IsFirstCommand = true
					terminal.SetDescription("A sleek terminal sits on the desk...")
				} else {
					response.Message = fmt.Sprintf("bash: %s: command not found", input)
				}
				return response
			} else {
				if input == "cat unlock-exits-instructions.txt" {
					response.Message = "Victory Achieved! The doors swing wide."
					response.GameOver = true
					global.GameOver = true
					return response
				} else {
					response.Message = fmt.Sprintf("bash: %s: command not found", input)
					return response
				}
			}
		}

		result := executeCommand(playerInput, game)
		response.Message = result
		response.GameOver = global.GameOver
		return response
	}
	return response
}

func (game *Game) SetupGame() {

	game.introduction = "It's the last day at the Academy, and you and your fellow graduates are ready to take on the final hack-day challenge.\nHowever, this time, it's different. Alan and Dan, your instructors, have prepared something more intense than ever before — a true test of your problem-solving and coding skills.\nThe doors to the academy are locked, the windows sealed. The only way out is to find and solve a series of riddles that lead to the terminal in a hidden room.\nThe challenge? Crack the code on the terminal to unlock the doors. But it's not that simple.\nYou'll need to gather items, approach Alan and Dan for cryptic tips, and outsmart the obstacles they've laid out for you.\nAs the tension rises, only your wits, teamwork, and knowledge can guide you to freedom.\nAre you ready to escape?\nOh and remember... You don't want to make Rosie grumpy! So don't do anything crazy.\n\nif at any point you feel lost, type 'commands' to display the list of all commands.\nThe command 'look' is always useful to get your bearings and see the options available to you.\nThe command 'exit' will make you quit the game at any time. Make sure you do mean to use it, or you will inadvertently lose all of your progress!"

	game.introductionShown = false

	ValidInteractions = []*Interaction{
		{
			ItemName:   "tea",
			EntityName: "rosie",
			Event:      &Event{Description: "get-your-lanyard", Outcome: "Cheers! I needed that... by the way, where is your lanyard? I must have forgotten to give it to you.\nYou'll need that to move between rooms, here it is.\n\n(lanyard can now be found in the room).\n", Triggered: false},
		},
		{
			ItemName:   "first-plate",
			EntityName: "dishwasher",
			Event:      &Event{Description: "first-plate-loaded", Outcome: "You loaded the first plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName:   "second-plate",
			EntityName: "dishwasher",
			Event:      &Event{Description: "second-plate-loaded", Outcome: "You loaded the second plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName:   "third-plate",
			EntityName: "dishwasher",
			Event:      &Event{Description: "third-plate-loaded", Outcome: "You loaded the third plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName:   "fourth-plate",
			EntityName: "dishwasher",
			Event:      &Event{Description: "fourth-plate-loaded", Outcome: "You loaded the fourth plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName:   "fifth-plate",
			EntityName: "dishwasher",
			Event:      &Event{Description: "fifth-plate-loaded", Outcome: "You loaded the fifth plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName:   "sixth-plate",
			EntityName: "dishwasher",
			Event:      &Event{Description: "sixth-plate-loaded", Outcome: "You loaded the sixth plate into the dishwasher.", Triggered: false},
		},
	}

	game.dishwasherChallengeWon = &Event{Description: "dishwasher-loaded", Outcome: "You load the dirty plates into the dishwasher and switch it on, a feeling of being used washing over you.\nThis challenge felt less like teamwork and more like being roped into someone else's mess.\nWith a sigh, you decide to head back to Alan to see if this effort has truly led you to victory...\n", Triggered: false}



	game.unlockComputer = &Event{Description: "computer-is-unlocked", Outcome: "You enter the password, holding your breath. Yes! The screen flickers to life.\nyou've unlocked the computer and now have full access.\n\nYou should approach Alan to find out what's next...\n", Triggered: false}

	game.computerPassword = "iiwsccrtc"

	game.remainingPasswordAttempts = 10

	game.staffRoom = &Room{
		Name:        "break-room",
		Description: "A cozy lounge designed for both academy students and tutors, offering a welcoming space to unwind and socialise.\nComfortable seating invites you to relax, while the warm ambiance encourages lively conversations and friendly exchanges.",
		Items:       make(map[string]*Item),
		Entities:    make(map[string]*Entity),
		Exits:       make(map[string]*Room),
	}

	game.codingLab = &Room{
		Name:        "coding-lab",
		Description: "A bright, tech-filled room with sleek workstations, whiteboards, and collaborative spaces.\nThe air buzzes with creativity as students code, share ideas, and tackle challenges together.",
		Items:       make(map[string]*Item),
		Entities:    make(map[string]*Entity),
		Exits:       make(map[string]*Room),
	}

	game.terminalRoom = &Room{
		Name:        "terminal-room",
		Description: "As you step into the terminal room, you're greeted by the soft hum of machines and the flickering glow of monitors lining the walls.\n\nThe air is charged with a sense of urgency, filled with the scent of freshly brewed coffee mingling with the faint odor of electrical components.\n\nIn the center of the room, a sleek, state-of-the-art terminal stands atop a polished wooden desk.",
		Items:       make(map[string]*Item),
		Entities:    make(map[string]*Entity),
		Exits:       make(map[string]*Room),
	}

	game.staffRoom.Exits["south"] = game.codingLab
	game.codingLab.Exits["north"] = game.staffRoom
	game.codingLab.Exits["east"] = game.terminalRoom
	game.terminalRoom.Exits["west"] = game.codingLab

	game.staffRoom.Items["tea"] = &Item{Name: "tea", Description: "A steaming cup of Yorkshire tea, rich and comforting.", Weight: 2, Hidden: true}
	game.staffRoom.Items["lanyard"] = &Item{Name: "lanyard", Description: "Your lanyard, a key to unlocking any door within the building.", Weight: 1, Hidden: true}
	game.staffRoom.Items["abandoned-lanyard"] = &Item{Name: "abandoned-lanyard", Description: "An abandoned lanyard, a key to unlocking any door within the building.", Weight: 1, Hidden: true}
	game.codingLab.Items["first-plate"] = &Item{Name: "first-plate", Description: "The plate on top of the stack.", Weight: 6, Hidden: true}
	game.codingLab.Items["second-plate"] = &Item{Name: "second-plate", Description: "The second plate of the stack.", Weight: 6, Hidden: true}
	game.codingLab.Items["third-plate"] = &Item{Name: "third-plate", Description: "The third plate of the stack.", Weight: 6, Hidden: true}
	game.codingLab.Items["fourth-plate"] = &Item{Name: "fourth-plate", Description: "The fourth plate of the stack.", Weight: 6, Hidden: true}
	game.codingLab.Items["fifth-plate"] = &Item{Name: "fifth-plate", Description: "The fifth plate of the stack.", Weight: 6, Hidden: true}
	game.codingLab.Items["sixth-plate"] = &Item{Name: "sixth-plate", Description: "The plate at the bottom of the stack.", Weight: 6, Hidden: true}
	game.codingLab.Items["cd"] = &Item{Name: "cd", Description: "A compact disc with '\\secret-files' written on it in bold letters.\nIt almost seems to call out to you, hinting at hidden knowledge.", Weight: 1, Hidden: false}

	game.staffRoom.Entities["rosie"] = &Entity{Name: "rosie", Description: "Ugh, what? Sorry, I can't think straight without a brew. Get me some tea, and then we'll talk...", Hidden: false}
	game.staffRoom.Entities["kettle"] = &Entity{Name: "kettle", Description: "You set the kettle to boil, brewing the strongest cup of tea you've ever made. A comforting aroma fills the room as the tea is now ready.\n\n(tea can now be found in the room)\n", Hidden: false}
	game.staffRoom.Entities["sofa"] = &Entity{Name: "sofa", Description: "You come across one of your fellow academy students fast asleep on the sofa. Next to them, their lanyard lies carelessly within reach.\nYou know you shouldn't take it, but the temptation lingers...\n\n(abandoned-lanyard can now be found in the room)\n", Hidden: false}
	game.staffRoom.Entities["dishwasher"] = &Entity{Name: "dishwasher", Description: "A stainless steel dishwasher sits quietly in the corner, its door slightly ajar.\nThe faint scent of soap lingers, and the racks inside are half-empty, waiting for the next load of dirty dishes to be placed inside.\nIt hums faintly, as if anticipating the task it was built for.", Hidden: true}
	game.staffRoom.Entities["cat"] = &Entity{Name: "cat", Description: "On one of the chairs, a fluffy cat lounges lazily, wearing a collar with a name tag that reads 'unlock-exits-instructions.txt'\n\nAn odd name for a cat. You get the feeling that this feline is more than it seems, possibly guarding crucial information", Hidden: false}
	game.codingLab.Entities["computer"] = &Entity{Name: "computer", Description: "Alan's computer. You need the password to get in.\n\nRemaining attempts: 10.\n\nType 'leave' to stop entering the password.\n\nEnter the password:\n", Hidden: false}
	game.codingLab.Entities["alan"] = &Entity{Name: "alan", Description: "Oh, you've finally made it... What are you waiting for, crack on with the code. The computer is right there...\nWhat's that? You don't know the password? Hmm... I seem to have forgotten it myself, but I do recall it's nine letters long.\nAnd for the love of all that's good, it's definitely not 'waterfall'!", Hidden: false}
	game.codingLab.Entities["agile-manifesto"] = &Entity{Name: "agile-manifesto", Description: "A large, framed document hangs prominently on the wall, its edges slightly frayed\nYou can almost feel the energy of past brainstorming sessions in the air as you read the four key values:\n\nIndividuals and Interactions over processes and tools.\n\nWorking Software over comprehensive documentation.\n\nCustomer Collaboration over contract negotiation.\n\nResponding To Change over following a plan.\n", Hidden: false}
	game.codingLab.Entities["desk"] = &Entity{Name: "desk", Description: "You approach the desk and spot a messy pile of dirty plates, stacked haphazardly. You think to yourself that somebody was too lazy to load the dishwasher.\nThe stack is too heavy to carry all the plates at once, and taking plates from the centre or bottom of the stack could pose a risk...\n\n(stack of plates can now be found in the room)\n\n", Hidden: true}
	game.terminalRoom.Entities["terminal"] = &Entity{Name: "terminal", Description: "A sleek terminal sits on the desk, its screen displaying lines of code and system commands.\nThe keyboard, slightly worn, hints at frequent use.\nThis device is essential for executing tasks and accessing the building's network.\n\nEnter your commands below or type 'leave' to exit the terminal.\n\n", Hidden: true}
	game.terminalRoom.Entities["dan"] = &Entity{Name: "dan", Description: "Congratulations on making it this far! I must say, I'm genuinely impressed. It appears I'm your final boss — muahahaha!\n...Oh, pardon my theatrics. Now, listen closely: the terminal holds the secret instructions to escape the building.\nYou only need two commands to access them.\nLook around the building to find some clues...\nYes, I know, this is actually the easiest task so far. If I am being totally honest, we just want to be done by 4pm...\nWhat are you standing there for? Get to it!\n", Hidden: true}

	game.isAttemptingPassword = false

	game.isAttemptingTerminal = false

	game.IsFirstCommand = false

	game.player = &Player{
		CurrentRoom:     game.staffRoom,
		Inventory:       make(map[string]*Item),
		AvailableWeight: 20,
		CurrentEntity:   nil,
	}

}
