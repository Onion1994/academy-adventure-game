package model

import (
	"academy-adventure-game/global"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// TODO: Possibility:
// Would a map of Command be useful? If Command was an interface with one method Execute( input PlayerInput ).
// Maybe even Execute( w World, input PlayerInput ) World? Just thoughts top of head
// PlayerInput would be a struct with methods that wrapped the text the player typed in, and allowed it to be parsed/used (a great TDD candidate)

type Game struct {
	player *Player
	// validInteractions         []*Interaction
	// gameOver                  bool
	introduction              string
	introductionShown         bool
	dishwasherChallengeWon    *Event
	unlockComputer            *Event
	grumpyRosie               *Event
	remainingPasswordAttempts int
	computerPassword          string
	isAttemptingPassword      bool
	isAttemptingTerminal      bool
	IsFirstCommand            bool
	staffRoom                 *Room
	codingLab                 *Room
	terminalRoom              *Room
}

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

func executeCommand(command string, args []string, g *Game) bool {

	input := &PlayerInput{Command: command, Args: args}
	cmd, exists := Commands[command]

	if command == g.computerPassword {
		return true
	}
	if !exists {
		fmt.Printf("Unknown Command : %s", command)
		return false
	}
	cmd.Execute(*input, g)
	return true
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ShowCommands(d Display) {
	d.Show("-exit -> quits the game\n\n-commands -> shows the commands\n\n-look -> shows the content of the room.\n\n-approach <entity> -> to approach an entity\n\n-leave -> to leave an entity\n\n-inventory -> shows items in the inventory\n\n-take <item> -> to take an item into your inventory\n\n-drop <item> -> to drop an item from your inventory and move it to the current room\n\n-use <item> -> to make use of a certain item when you approach an entity\n\n-move <direction> -> to move to a different room\n\n-map -> shows the directions you can take\n")
}

func (g *Game) RunGame() {
	abandonedLanyard := g.staffRoom.Items["abandoned-lanyard"]
	tea := g.staffRoom.Items["tea"]
	lanyard := g.staffRoom.Items["lanyard"]
	firstPlate := g.codingLab.Items["first-plate"]
	secondPlate := g.codingLab.Items["second-plate"]
	thirdPlate := g.codingLab.Items["third-plate"]
	fourthPlate := g.codingLab.Items["fourth-plate"]
	fifthPlate := g.codingLab.Items["fifth-plate"]
	sixthPlate := g.codingLab.Items["sixth-plate"]

	sofa := g.staffRoom.Entities["sofa"]
	terminal := g.terminalRoom.Entities["terminal"]
	computer := g.codingLab.Entities["computer"]
	kettle := g.staffRoom.Entities["kettle"]
	dishwasher := g.staffRoom.Entities["dishwasher"]
	desk := g.codingLab.Entities["desk"]
	alan := g.codingLab.Entities["alan"]
	dan := g.terminalRoom.Entities["dan"]
	rosie := g.staffRoom.Entities["rosie"]
	// agileManifesto := codingLab.Entities["agile-manifesto"]
	// cd := codingLab.Items["cd"]
	// cat := staffRoom.Entities["cat"]

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if g.player.CurrentEntity != nil && g.player.CurrentEntity.Name == "sofa" {
			abandonedLanyard.Hidden = false
			sofa.SetDescription("Your fellow academy student continues to sleep on the sofa. Something tells you it's down to you to get stuff done today...")
		}

		if g.player.CurrentEntity != nil && g.player.CurrentEntity.Name == "kettle" {
			tea.Hidden = false
			kettle.SetDescription("A kettle — essential for survival, impossible to function without one nearby.")
		}

		if g.player.CurrentEntity != nil && g.player.CurrentEntity.Name == "desk" {
			firstPlate.Hidden = false
			secondPlate.Hidden = false
			thirdPlate.Hidden = false
			fourthPlate.Hidden = false
			fifthPlate.Hidden = false
			sixthPlate.Hidden = false
			desk.SetDescription("Despite the disarray, it's clear this desk sees frequent use, with just enough space left to get work done.")
		}

		for _, validInteraction := range ValidInteractions {
			if validInteraction.Event.Description == "get-your-lanyard" && validInteraction.Event.Triggered {
				lanyard.Hidden = false
				rosie.SetDescription("Can I help with anything else?")
			}
		}

		dishwasherLoaded := true
		for _, validInteraction := range ValidInteractions {
			if strings.HasSuffix(validInteraction.ItemName, "plate") && !validInteraction.Event.Triggered {
				dishwasherLoaded = false
				break
			}
		}

		if !g.dishwasherChallengeWon.Triggered {
			if dishwasherLoaded {
				g.player.TriggerEvent(g.dishwasherChallengeWon)
				alan.SetDescription("Ah, so you've managed to load the dishwasher! Splendid work — consider this challenge complete.\nI could have done it myself instead of writing that clever recursive function, but where's the fun in that?\nAfter all, they pay me for my intellect, not for doing the heavy lifting!\nBut I digress. You're free to proceed to the terminal room and speak with Dan for your final challenge.\nYou're doing an excellent job; keep it up!")
				dan.Hidden = false
				terminal.Hidden = false
			}
		}

		if _, ok := g.player.Inventory["abandoned-lanyard"]; ok {
			g.player.TriggerEvent(g.grumpyRosie)
			global.GameOver = true
		}

		if global.GameOver {
			fmt.Println("Thank you for playing!")
			break
		}

		if !g.introductionShown {
			clearScreen()
			fmt.Println(g.introduction)
			g.introductionShown = true
		}

		fmt.Print("Enter command: ")

		if scanner.Scan() {
			input := scanner.Text()
			input = strings.TrimSpace(input)
			input = strings.ToLower(input)

			if input == "exit" {
				clearScreen()
				fmt.Println("Thank you for playing!")
				break
			}

			if g.isAttemptingPassword {
				if g.remainingPasswordAttempts == 1 && input != g.computerPassword {
					clearScreen()
					fmt.Println("Alan's computer is locked, halting your progress in the challenge. To top it off, you've made Rosie grumpy, as she'll now have to take the computer to IT.\n\nThank you for playing!")
					break
				}
				if input == g.computerPassword {
					clearScreen()
					g.player.TriggerEvent(g.unlockComputer)
					computer.SetDescription("function completeTask(pile)\n   if pile == 0:\n      return 'Task Complete'\n   else:\n      completeTask(pile - 1)\n")
					alan.SetDescription("You've cracked the password! Impressive work... You should now see an open file containing a recursive function.\n\nFollow its instructions carefully, and you'll be one step closer to victory!\nBut, a word of caution: the task ahead is, well, a bit more hands-on than you might expect...")
					g.isAttemptingPassword = false
					desk.Hidden = false
					dishwasher.Hidden = false
				} else if input == "leave" {
					g.isAttemptingPassword = false
				} else {
					g.remainingPasswordAttempts--
					clearScreen()
					fmt.Printf("Incorrect password. Try again, or type 'leave' to stop entering the password.\n\nRemaining attempts: %d\n\n", g.remainingPasswordAttempts)
					computer.SetDescription(fmt.Sprintf("Alan's computer. You need the password to get in.\nRemaining attempts: %d.\nType 'leave' to stop entering the password.\n\nEnter the password:\n", g.remainingPasswordAttempts))
					continue
				}
			}

			if g.isAttemptingTerminal {
				if input == "leave" {
					g.isAttemptingTerminal = false
					clearScreen()
					g.player.Leave()
					continue
				}

				if !g.IsFirstCommand {
					if input == "cd /secret-files" {
						clearScreen()
						fmt.Println("The terminal displays:\n\n/secret-files/\n\nIt looks like you are on the right track.\nEnter the final command to win the game!\n\nType 'leave' to stop entering commands on the terminal.")
						g.IsFirstCommand = true
						terminal.SetDescription("A sleek terminal sits on the desk, its screen displaying lines of code and system commands.\nThe keyboard, slightly worn, hints at frequent use.\nThis device is essential for executing tasks and accessing the building's network.\n\nEnter your commands below or type 'leave' to exit the terminal.\n\nThe terminal displays:\n\n/secret-files/\n\nIt looks like you are on the right track.\nEnter the final command to win the game!\n\nType 'leave' to stop entering commands on the terminal.\n")
						continue
					} else {
						clearScreen()
						fmt.Printf("The terminal displays:\n\nbash: %s: command not found\n\nType 'leave' to stop entering commands on the terminal\n\n", input)
						continue
					}
				} else {
					if input == "cat unlock-exits-instructions.txt" {
						clearScreen()
						fmt.Println("As you execute the final command, the terminal whirs to life, and the screen fills with a flurry of colorful text.\nThe words 'Victory Achieved!' flash across the display, illuminating your face with a soft glow.\nYou feel a rush of adrenaline as the file containing the instructions to unlock the exits appears before you.\nFollowing the instructions carefully, you swiftly input the necessary commands, and with a satisfying beep, the locks on the exits click open.\nThe room is filled with the sound of machinery grinding to a halt as the doors swing wide.")
						global.GameOver = true
						continue
					} else {
						clearScreen()
						fmt.Printf("The terminal displays:\n\nbash: %s: command not found\n\nType 'leave' to stop entering commands on the terminal\n\n", input)
						continue
					}
				}
			}

			parts := strings.Fields(input)
			if len(parts) == 0 {
				continue
			}

			command := (parts[0])
			args := parts[1:]

			executionResult := executeCommand(command, args, g)

			if executionResult {
				continue
			}
		}
	}
}

func (g *Game) SetupGame() {

	g.introduction = "It's the last day at the Academy, and you and your fellow graduates are ready to take on the final hack-day challenge.\nHowever, this time, it's different. Alan and Dan, your instructors, have prepared something more intense than ever before — a true test of your problem-solving and coding skills.\nThe doors to the academy are locked, the windows sealed. The only way out is to find and solve a series of riddles that lead to the terminal in a hidden room.\nThe challenge? Crack the code on the terminal to unlock the doors. But it's not that simple.\nYou'll need to gather items, approach Alan and Dan for cryptic tips, and outsmart the obstacles they've laid out for you.\nAs the tension rises, only your wits, teamwork, and knowledge can guide you to freedom.\nAre you ready to escape?\nOh and remember... You don't want to make Rosie grumpy! So don't do anything crazy.\n\nif at any point you feel lost, type 'commands' to display the list of all commands.\nThe command 'look' is always useful to get your bearings and see the options available to you.\nThe command 'exit' will make you quit the game at any time. Make sure you do mean to use it, or you will inadvertently lose all of your progress!"

	g.introductionShown = false

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

	g.dishwasherChallengeWon = &Event{Description: "dishwasher-loaded", Outcome: "You load the dirty plates into the dishwasher and switch it on, a feeling of being used washing over you.\nThis challenge felt less like teamwork and more like being roped into someone else's mess.\nWith a sigh, you decide to head back to Alan to see if this effort has truly led you to victory...\n", Triggered: false}

	g.grumpyRosie = &Event{Description: "rosie-is-grumpy", Outcome: "Rosie caught you in the act of swiping a lanyard from a fellow student.\nYou have made Rosie grumpy and you've lost the game.\n", Triggered: false}

	g.unlockComputer = &Event{Description: "computer-is-unlocked", Outcome: "You enter the password, holding your breath. Yes! The screen flickers to life.\nyou've unlocked the computer and now have full access.\n\nYou should approach Alan to find out what's next...\n", Triggered: false}

	g.computerPassword = "iiwsccrtc"

	g.remainingPasswordAttempts = 10

	g.staffRoom = &Room{
		Name:        "break-room",
		Description: "A cozy lounge designed for both academy students and tutors, offering a welcoming space to unwind and socialise.\nComfortable seating invites you to relax, while the warm ambiance encourages lively conversations and friendly exchanges.",
		Items:       make(map[string]*Item),
		Entities:    make(map[string]*Entity),
		Exits:       make(map[string]*Room),
	}

	g.codingLab = &Room{
		Name:        "coding-lab",
		Description: "A bright, tech-filled room with sleek workstations, whiteboards, and collaborative spaces.\nThe air buzzes with creativity as students code, share ideas, and tackle challenges together.",
		Items:       make(map[string]*Item),
		Entities:    make(map[string]*Entity),
		Exits:       make(map[string]*Room),
	}

	g.terminalRoom = &Room{
		Name:        "terminal-room",
		Description: "As you step into the terminal room, you're greeted by the soft hum of machines and the flickering glow of monitors lining the walls.\n\nThe air is charged with a sense of urgency, filled with the scent of freshly brewed coffee mingling with the faint odor of electrical components.\n\nIn the center of the room, a sleek, state-of-the-art terminal stands atop a polished wooden desk.",
		Items:       make(map[string]*Item),
		Entities:    make(map[string]*Entity),
		Exits:       make(map[string]*Room),
	}

	g.staffRoom.Exits["south"] = g.codingLab
	g.codingLab.Exits["north"] = g.staffRoom
	g.codingLab.Exits["east"] = g.terminalRoom
	g.terminalRoom.Exits["west"] = g.codingLab

	g.staffRoom.Items["tea"] = &Item{Name: "tea", Description: "A steaming cup of Yorkshire tea, rich and comforting.", Weight: 2, Hidden: true}
	g.staffRoom.Items["lanyard"] = &Item{Name: "lanyard", Description: "Your lanyard, a key to unlocking any door within the building.", Weight: 1, Hidden: true}
	g.staffRoom.Items["abandoned-lanyard"] = &Item{Name: "abandoned-lanyard", Description: "An abandoned lanyard, a key to unlocking any door within the building.", Weight: 1, Hidden: true}
	g.codingLab.Items["first-plate"] = &Item{Name: "first-plate", Description: "The plate on top of the stack.", Weight: 6, Hidden: true}
	g.codingLab.Items["second-plate"] = &Item{Name: "second-plate", Description: "The second plate of the stack.", Weight: 6, Hidden: true}
	g.codingLab.Items["third-plate"] = &Item{Name: "third-plate", Description: "The third plate of the stack.", Weight: 6, Hidden: true}
	g.codingLab.Items["fourth-plate"] = &Item{Name: "fourth-plate", Description: "The fourth plate of the stack.", Weight: 6, Hidden: true}
	g.codingLab.Items["fifth-plate"] = &Item{Name: "fifth-plate", Description: "The fifth plate of the stack.", Weight: 6, Hidden: true}
	g.codingLab.Items["sixth-plate"] = &Item{Name: "sixth-plate", Description: "The plate at the bottom of the stack.", Weight: 6, Hidden: true}
	g.codingLab.Items["cd"] = &Item{Name: "cd", Description: "A compact disc with '\\secret-files' written on it in bold letters.\nIt almost seems to call out to you, hinting at hidden knowledge.", Weight: 1, Hidden: false}

	g.staffRoom.Entities["rosie"] = &Entity{Name: "rosie", Description: "Ugh, what? Sorry, I can't think straight without a brew. Get me some tea, and then we'll talk...", Hidden: false}
	g.staffRoom.Entities["kettle"] = &Entity{Name: "kettle", Description: "You set the kettle to boil, brewing the strongest cup of tea you've ever made. A comforting aroma fills the room as the tea is now ready.\n\n(tea can now be found in the room)\n", Hidden: false}
	g.staffRoom.Entities["sofa"] = &Entity{Name: "sofa", Description: "You come across one of your fellow academy students fast asleep on the sofa. Next to them, their lanyard lies carelessly within reach.\nYou know you shouldn't take it, but the temptation lingers...\n\n(abandoned-lanyard can now be found in the room)\n", Hidden: false}
	g.staffRoom.Entities["dishwasher"] = &Entity{Name: "dishwasher", Description: "A stainless steel dishwasher sits quietly in the corner, its door slightly ajar.\nThe faint scent of soap lingers, and the racks inside are half-empty, waiting for the next load of dirty dishes to be placed inside.\nIt hums faintly, as if anticipating the task it was built for.", Hidden: true}
	g.staffRoom.Entities["cat"] = &Entity{Name: "cat", Description: "On one of the chairs, a fluffy cat lounges lazily, wearing a collar with a name tag that reads 'unlock-exits-instructions.txt'\n\nAn odd name for a cat. You get the feeling that this feline is more than it seems, possibly guarding crucial information", Hidden: false}
	g.codingLab.Entities["computer"] = &Entity{Name: "computer", Description: "Alan's computer. You need the password to get in.\n\nRemaining attempts: 10.\n\nType 'leave' to stop entering the password.\n\nEnter the password:\n", Hidden: false}
	g.codingLab.Entities["alan"] = &Entity{Name: "alan", Description: "Oh, you've finally made it... What are you waiting for, crack on with the code. The computer is right there...\nWhat's that? You don't know the password? Hmm... I seem to have forgotten it myself, but I do recall it's nine letters long.\nAnd for the love of all that's good, it's definitely not 'waterfall'!", Hidden: false}
	g.codingLab.Entities["agile-manifesto"] = &Entity{Name: "agile-manifesto", Description: "A large, framed document hangs prominently on the wall, its edges slightly frayed\nYou can almost feel the energy of past brainstorming sessions in the air as you read the four key values:\n\nIndividuals and Interactions over processes and tools.\n\nWorking Software over comprehensive documentation.\n\nCustomer Collaboration over contract negotiation.\n\nResponding To Change over following a plan.\n", Hidden: false}
	g.codingLab.Entities["desk"] = &Entity{Name: "desk", Description: "You approach the desk and spot a messy pile of dirty plates, stacked haphazardly. You think to yourself that somebody was too lazy to load the dishwasher.\nThe stack is too heavy to carry all the plates at once, and taking plates from the centre or bottom of the stack could pose a risk...\n\n(stack of plates can now be found in the room)\n\n", Hidden: true}
	g.terminalRoom.Entities["terminal"] = &Entity{Name: "terminal", Description: "A sleek terminal sits on the desk, its screen displaying lines of code and system commands.\nThe keyboard, slightly worn, hints at frequent use.\nThis device is essential for executing tasks and accessing the building's network.\n\nEnter your commands below or type 'leave' to exit the terminal.\n\n", Hidden: true}
	g.terminalRoom.Entities["dan"] = &Entity{Name: "dan", Description: "Congratulations on making it this far! I must say, I'm genuinely impressed. It appears I'm your final boss — muahahaha!\n...Oh, pardon my theatrics. Now, listen closely: the terminal holds the secret instructions to escape the building.\nYou only need two commands to access them.\nLook around the building to find some clues...\nYes, I know, this is actually the easiest task so far. If I am being totally honest, we just want to be done by 4pm...\nWhat are you standing there for? Get to it!\n", Hidden: true}

	g.isAttemptingPassword = false

	g.isAttemptingTerminal = false

	g.IsFirstCommand = false

	g.player = &Player{
		CurrentRoom:     g.staffRoom,
		Inventory:       make(map[string]*Item),
		AvailableWeight: 20,
		CurrentEntity:   nil,
	}

}
