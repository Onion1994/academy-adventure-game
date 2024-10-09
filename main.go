package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Display interface {
	Show( text string )
 }

 type ConsoleDisplay struct {}

 func (c ConsoleDisplay) Show(text string) {
	fmt.Println(text)
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

type Describable interface {
	SetDescription(description string)
	GetDescription() string
}

type Item struct {
	Name string
	Description string
	Weight int
	Hidden bool
}

func (i *Item) SetDescription(description string) {
    i.Description = description
}

func (i *Item) GetDescription() string {
    return i.Description
}

type Room struct {
	Name string
	Description string
	Exits map[string]*Room
	Items map[string]*Item
	Entities map[string]*Entity
}

func (r *Room) SetDescription(description string) {
    r.Description = description
}

func (r *Room) GetDescription() string {
    return r.Description
}

type Player struct {
	CurrentRoom *Room
	Inventory   map[string]*Item
	CurrentEntity *Entity
	CarriedWeight int
	AvailableWeight int
}

type Entity struct {
	Name string
	Description string
	Hidden bool
}

type Event struct {
	Description string
	Outcome string
	Triggered bool
}

type Interaction struct {
	ItemName string
	EntityName string
	Event *Event
}

var validInteractions = []*Interaction{}

func (e *Entity) SetDescription(description string) {
    e.Description = description
}

func (e *Entity) GetDescription() string {
    return e.Description
}

func (p *Player) Move(direction string, d Display) {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
	}
	if newRoom, ok := p.CurrentRoom.Exits[direction]; ok {
		p.CurrentRoom = newRoom

		d.Show(fmt.Sprintf("You are in %s\n", p.CurrentRoom.Name))
	} else {
		d.Show("You can't go that way!\n")
	}
}

var plateOrder = []string{"first-plate", "second-plate", "third-plate", "fourth-plate", "fifth-plate", "sixth-plate"}
var currentPlateIndex = 0
var gameOver = false


func (p *Player) Take(itemName string, d Display) {
	item, ok := p.CurrentRoom.Items[itemName]
	switch {
	case !ok || item.Hidden:
		d.Show(fmt.Sprintf("You can't take %s\n", itemName))
		return
	case p.AvailableWeight < item.Weight:
		d.Show(fmt.Sprintln("Weight limit reached! Please drop an item before taking more."))
		return
	case isPlate(itemName):
		if itemName == plateOrder[currentPlateIndex] {
			p.Inventory[item.Name] = item
			p.ChangeCarriedWeight(item, "increase")
			delete(p.CurrentRoom.Items, item.Name)
			currentPlateIndex++

			d.Show(fmt.Sprintf("%s has been added to your inventory.\n", item.Name))
		} else {
			d.Show(fmt.Sprintln("As you attempt to grab the greasy plates without removing the ones stacked above them, they slip from your grasp and shatter, creating a chaotic mess.\n\nNow Rosie is very grumpy."))
			gameOver = true
		}
		
	default:
		p.Inventory[item.Name] = item
		p.ChangeCarriedWeight(item, "increase")
		delete(p.CurrentRoom.Items, item.Name)

		d.Show(fmt.Sprintf("%s has been added to your inventory.\n", item.Name))
	}
}

func isPlate(itemName string) bool {
    for _, plate := range plateOrder {
        if itemName == plate {
            return true
        }
    }
    return false
}

func (p *Player) ChangeCarriedWeight(item *Item, operation string) {
	switch {
	case operation == "increase":
		p.CarriedWeight += item.Weight
		p.AvailableWeight -= item.Weight
		return
	case operation == "decrease":
		p.CarriedWeight -= item.Weight
		p.AvailableWeight += item.Weight
		return
	}
}

func (p *Player) Drop(itemName string, d Display) {
	if item, ok := p.Inventory[itemName]; ok {
		if isPlate(itemName) {
			d.Show("You can't just leave those plates lying around! It's time to load them into the dishwasher!")
			return
		}

		delete(p.Inventory, item.Name)
		p.ChangeCarriedWeight(item, "decrease")
		p.CurrentRoom.Items[item.Name] = item

		d.Show(fmt.Sprintf("You dropped %s.\n\n", item.Name))
	} else {
		d.Show(fmt.Sprintf("You don't have %s.\n\n", itemName))
	}
}

func (p *Player) ShowInventory(d Display) {
	if len(p.Inventory) == 0 {
		d.Show(fmt.Sprintf("Your inventory is empty.\nAvailable space: %d\n", p.AvailableWeight))
		return
	}
	d.Show(fmt.Sprintf("Available space: %d\nYour inventory contains:\n", p.AvailableWeight))
	for itemName, item := range p.Inventory {
		d.Show(fmt.Sprintf("- %s: %s Weight: %d\n", itemName, item.Description, item.Weight))
	}
}

func (p *Player) ShowRoom(d Display) {
    d.Show(fmt.Sprintf("You are in %s\n\n%s\n", p.CurrentRoom.Name, p.CurrentRoom.Description))

	if p.EntitiesArePresent() {
			d.Show("\nYou can approach:\n")
			for _, entity := range p.CurrentRoom.Entities {
				switch {
				case p.CurrentEntity != nil:
					if entity.Name == p.CurrentEntity.Name {
						d.Show(fmt.Sprintf("- %s (currently approached)\n", entity.Name))
					} else if !entity.Hidden{
						d.Show(fmt.Sprintf("- %s\n", entity.Name))
					}
				default:
					if !entity.Hidden{
						d.Show(fmt.Sprintf("- %s\n", entity.Name))
					}
				}
			}
		}
		
	if p.ItemsArePresent() {
			d.Show("\nThe room contains:")
			for itemName, item := range p.CurrentRoom.Items {
				if !item.Hidden {
					d.Show(fmt.Sprintf("\n- %s: %s Weight: %d\n", itemName, item.Description, item.Weight))
				}
			}
		}
	}

	

func (p *Player) ItemsArePresent() bool {
	if len(p.CurrentRoom.Items) != 0 {
	for _, item := range p.CurrentRoom.Items {
		if !item.Hidden {
			return true
		}
	}
}
	return false
}

func (p *Player) EntitiesArePresent() bool {
	if len(p.CurrentRoom.Entities) != 0 {
	for _, entity := range p.CurrentRoom.Entities {
		if !entity.Hidden {
			return true
		}
	}
}
	return false
}


func (p *Player) Approach(entityName string) {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
	}
	if entity, ok := p.CurrentRoom.Entities[entityName]; ok && !entity.Hidden{

		p.CurrentEntity = entity
		fmt.Println(entity.Description)
	} else {
		fmt.Printf("You can't approach %s.\n", entityName)
	}
}

func (p *Player) Leave() {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
		p.ShowRoom(ConsoleDisplay{})
	} else {
		fmt.Println("You have not approached anything. If you wish to leave the game, use the exit command.")
	}
}

func (p *Player) ShowMap(d Display) {
	for direction, exit := range p.CurrentRoom.Exits {
		d.Show(fmt.Sprintf("%s: %s\n", direction, exit.Name))
	}
}

func (p *Player) Use(itemName string, target string) {
	if p.CurrentEntity == nil {
		fmt.Println("Approach to use an item.")
		return
	}	
	if p.CurrentEntity.Name == target {
		if _, ok := p.Inventory[itemName]; ok {
				for _, interaction := range validInteractions {
					if interaction.ItemName == itemName && interaction.EntityName == target {
						p.TriggerEvent(interaction.Event)
						p.ChangeCarriedWeight(p.Inventory[itemName], "decrease")
						delete(p.Inventory, itemName)
						return
					}
				}
		} else {
			fmt.Printf("You don't have %s.\n", itemName)
			return
		}
	} else {
		fmt.Printf("%s not found.\n", target)
		return
	}
	fmt.Printf("You can't use %s on %s.\n", itemName, target)
}

func (p *Player) TriggerEvent(event *Event) {
	fmt.Println(event.Outcome)
	event.Triggered = true
}

func updateDescription(d Describable, newDescription string) {
	d.SetDescription(newDescription)
}

func showCommands(d Display) {
	d.Show("-exit -> quits the game\n\n-commands -> shows the commands\n\n-look -> shows the content of the room.\n\n-approach <entity> -> to approach an entity\n\n-leave -> to leave an entity\n\n-inventory -> shows items in the inventory\n\n-take <item> -> to take an item into your inventory\n\n-drop <item> -> to drop an item from your inventory and move it to the current room\n\n-use <item> -> to make use of a certain item when you approach an entity\n\n-move <direction> -> to move to a different room\n\n-map -> shows the directions you can take\n")
}

func main() {
	introduction := "It's the last day at the Academy, and you and your fellow graduates are ready to take on the final hack-day challenge.\nHowever, this time, it's different. Alan and Dan, your instructors, have prepared something more intense than ever before — a true test of your problem-solving and coding skills.\nThe doors to the academy are locked, the windows sealed. The only way out is to find and solve a series of riddles that lead to the terminal in a hidden room.\nThe challenge? Crack the code on the terminal to unlock the doors. But it's not that simple.\nYou'll need to gather items, approach Alan and Dan for cryptic tips, and outsmart the obstacles they've laid out for you.\nAs the tension rises, only your wits, teamwork, and knowledge can guide you to freedom.\nAre you ready to escape?\nOh and remember... You don't want to make Rosie grumpy! So don't do anything crazy.\n\nif at any point you feel lost, type 'commands' to display the list of all commands.\nThe command 'look' is always useful to get your bearings and see the options available to you.\nThe command 'exit' will make you quit the game at any time. Make sure you do mean to use it, or you will inadvertently lose all of your progress!"

	introductionShown:= false

	validInteractions = []*Interaction{
		{
			ItemName:   "tea",
			EntityName: "rosie",
			Event:      &Event{Description: "get-your-lanyard", Outcome: "Cheers! I needed that... by the way, where is your lanyard? I must have forgotten to give it to you.\nYou'll need that to move between rooms, here it is.\n\n(lanyard can now be found in the room).\n", Triggered: false},
		},
		{
			ItemName: "first-plate",
			EntityName: "dishwasher",
			Event: &Event{Description: "first-plate-loaded", Outcome: "You loaded the first plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName: "second-plate",
			EntityName: "dishwasher",
			Event: &Event{Description: "second-plate-loaded", Outcome: "You loaded the second plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName: "third-plate",
			EntityName: "dishwasher",
			Event: &Event{Description: "third-plate-loaded", Outcome: "You loaded the third plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName: "fourth-plate",
			EntityName: "dishwasher",
			Event: &Event{Description: "fourth-plate-loaded", Outcome: "You loaded the fourth plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName: "fifth-plate",
			EntityName: "dishwasher",
			Event: &Event{Description: "fifth-plate-loaded", Outcome: "You loaded the fifth plate into the dishwasher.", Triggered: false},
		},
		{
			ItemName: "sixth-plate",
			EntityName: "dishwasher",
			Event: &Event{Description: "sixth-plate-loaded", Outcome: "You loaded the sixth plate into the dishwasher.", Triggered: false},
		},
	}

	dishwasherChallengeWon := &Event{Description: "dishwasher-loaded", Outcome: "You load the dirty plates into the dishwasher and switch it on, a feeling of being used washing over you.\nThis challenge felt less like teamwork and more like being roped into someone else's mess.\nWith a sigh, you decide to head back to Alan to see if this effort has truly led you to victory...\n", Triggered: false}

	grumpyRosie := &Event{Description: "rosie-is-grumpy", Outcome: "Rosie caught you in the act of swiping a lanyard from a fellow student.\nYou have made Rosie grumpy and you've lost the game.\n", Triggered: false}

	unlockComputer := &Event{Description: "computer-is-unlocked", Outcome: "You enter the password, holding your breath. Yes! The screen flickers to life.\nyou've unlocked the computer and now have full access.\n\nYou should approach Alan to find out what's next...\n", Triggered: false}

	computerPassword := "iiwsccrtc"

	remainingPasswordAttempts := 10

	staffRoom := Room{
		Name:        "break-room",
		Description: "A cozy lounge designed for both academy students and tutors, offering a welcoming space to unwind and socialise.\nComfortable seating invites you to relax, while the warm ambiance encourages lively conversations and friendly exchanges.",
		Items:      make(map[string]*Item),
		Entities:   make(map[string]*Entity),
		Exits:      make(map[string]*Room),
	}

	codingLab := Room{
		Name:        "coding-lab",
		Description: "A bright, tech-filled room with sleek workstations, whiteboards, and collaborative spaces.\nThe air buzzes with creativity as students code, share ideas, and tackle challenges together.",
		Items:      make(map[string]*Item),
		Entities:   make(map[string]*Entity),
		Exits:      make(map[string]*Room),
	}

	terminalRoom := Room{
		Name:        "terminal-room",
		Description: "As you step into the terminal room, you're greeted by the soft hum of machines and the flickering glow of monitors lining the walls.\n\nThe air is charged with a sense of urgency, filled with the scent of freshly brewed coffee mingling with the faint odor of electrical components.\n\nIn the center of the room, a sleek, state-of-the-art terminal stands atop a polished wooden desk.",
		Items:      make(map[string]*Item),
		Entities:   make(map[string]*Entity),
		Exits:      make(map[string]*Room),
	}

	staffRoom.Exits["south"] = &codingLab
	codingLab.Exits["north"] = &staffRoom
	codingLab.Exits["east"] = &terminalRoom
	terminalRoom.Exits["west"] = &codingLab

	rosie := Entity{Name: "rosie", Description: "Ugh, what? Sorry, I can't think straight without a brew. Get me some tea, and then we'll talk...", Hidden: false}
	kettle := Entity{Name: "kettle", Description: "You set the kettle to boil, brewing the strongest cup of tea you've ever made. A comforting aroma fills the room as the tea is now ready.\n\n(tea can now be found in the room)\n", Hidden: false}
	sofa := Entity{Name: "sofa", Description: "You come across one of your fellow academy students fast asleep on the sofa. Next to them, their lanyard lies carelessly within reach.\nYou know you shouldn't take it, but the temptation lingers...\n\n(abandoned-lanyard can now be found in the room)\n", Hidden: false}
	tea := Item{Name: "tea", Description: "A steaming cup of Yorkshire tea, rich and comforting.", Weight: 2, Hidden: true}
	lanyard := Item{Name: "lanyard", Description: "Your lanyard, a key to unlocking any door within the building.", Weight: 1, Hidden: true}
	abandonedLanyard := Item{Name: "abandoned-lanyard", Description: "An abandoned lanyard, a key to unlocking any door within the building.", Weight: 1, Hidden: true}
	computer := Entity{Name: "computer", Description: "Alan's computer. You need the password to get in.\n\nRemaining attempts: 10.\n\nType 'leave' to stop entering the password.\n\nEnter the password:\n", Hidden: false}
	alan := Entity{Name: "alan", Description: "Oh, you've finally made it... What are you waiting for, crack on with the code. The computer is right there...\nWhat's that? You don't know the password? Hmm... I seem to have forgotten it myself, but I do recall it's nine letters long.\nAnd for the love of all that's good, it's definitely not 'waterfall'!", Hidden: false}
	agileManifesto := Entity{Name: "agile-manifesto", Description: "A large, framed document hangs prominently on the wall, its edges slightly frayed\nYou can almost feel the energy of past brainstorming sessions in the air as you read the four key values:\n\nIndividuals and Interactions over processes and tools.\n\nWorking Software over comprehensive documentation.\n\nCustomer Collaboration over contract negotiation.\n\nResponding To Change over following a plan.\n", Hidden: false}
	desk := Entity{Name: "desk", Description: "You approach the desk and spot a messy pile of dirty plates, stacked haphazardly. You think to yourself that somebody was too lazy to load the dishwasher.\nThe stack is too heavy to carry all the plates at once, and taking plates from the centre or bottom of the stack could pose a risk...\n\n(stack of plates can now be found in the room)\n\n", Hidden: true}
	dishwasher := Entity{Name: "dishwasher", Description: "A stainless steel dishwasher sits quietly in the corner, its door slightly ajar.\nThe faint scent of soap lingers, and the racks inside are half-empty, waiting for the next load of dirty dishes to be placed inside.\nIt hums faintly, as if anticipating the task it was built for.", Hidden: true}
	firstPlate := Item{Name: "first-plate", Description: "The plate on top of the stack.", Weight: 6, Hidden: true}
	secondPlate := Item{Name: "second-plate", Description: "The second plate of the stack.", Weight: 6, Hidden: true}
	thirdPlate := Item{Name: "third-plate", Description: "The third plate of the stack.", Weight: 6, Hidden: true}
	fourthPlate := Item{Name: "fourth-plate", Description: "The fourth plate of the stack.", Weight: 6, Hidden: true}
	fifthPlate := Item{Name: "fifth-plate", Description: "The fifth plate of the stack.", Weight: 6, Hidden: true}
	sixthPlate := Item{Name: "sixth-plate", Description: "The plate at the bottom of the stack.", Weight: 6, Hidden: true}
	terminal := Entity{Name: "terminal", Description: "A sleek terminal sits on the desk, its screen displaying lines of code and system commands.\nThe keyboard, slightly worn, hints at frequent use.\nThis device is essential for executing tasks and accessing the building's network.\n\nEnter your commands below or type 'leave' to exit the terminal.\n\n", Hidden: true}
	dan := Entity{Name: "dan", Description: "Congratulations on making it this far! I must say, I'm genuinely impressed. It appears I'm your final boss — muahahaha!\n...Oh, pardon my theatrics. Now, listen closely: the terminal holds the secret instructions to escape the building.\nYou only need two commands to access them.\nLook around the building to find some clues...\nYes, I know, this is actually the easiest task so far. If I am being totally honest, we just want to be done by 4pm...\nWhat are you standing there for? Get to it!\n", Hidden: true}
	cd := Item{Name: "cd", Description: "A compact disc with '\\secret-files' written on it in bold letters.\nIt almost seems to call out to you, hinting at hidden knowledge.", Weight: 1, Hidden: false}
	cat := Entity{Name: "cat", Description: "On one of the chairs, a fluffy cat lounges lazily, wearing a collar with a name tag that reads 'unlock-exits-instructions.txt'\n\nAn odd name for a cat. You get the feeling that this feline is more than it seems, possibly guarding crucial information", Hidden: false}

	staffRoom.Items[tea.Name] = &tea
	staffRoom.Items[lanyard.Name] = &lanyard
	staffRoom.Items[abandonedLanyard.Name] = &abandonedLanyard
	staffRoom.Entities[rosie.Name] = &rosie
	staffRoom.Entities[kettle.Name] = &kettle
	staffRoom.Entities[sofa.Name] = &sofa
	staffRoom.Entities[dishwasher.Name] = &dishwasher
	staffRoom.Entities[cat.Name] = &cat
	codingLab.Items[cd.Name] = &cd
	codingLab.Entities[computer.Name] = &computer
	codingLab.Entities[alan.Name] = &alan
	codingLab.Entities[agileManifesto.Name] = &agileManifesto
	codingLab.Entities[desk.Name] = &desk
	codingLab.Items[firstPlate.Name] = &firstPlate
	codingLab.Items[secondPlate.Name] = &secondPlate
	codingLab.Items[thirdPlate.Name] = &thirdPlate
	codingLab.Items[fourthPlate.Name] = &fourthPlate
	codingLab.Items[fifthPlate.Name] = &fifthPlate
	codingLab.Items[sixthPlate.Name] = &sixthPlate
	terminalRoom.Entities[terminal.Name] = &terminal
	terminalRoom.Entities[dan.Name] = &dan

	isAttemptingPassword := false

	isAttemptingTerminal := false

	IsFirstCommand := false

	player := Player{
		CurrentRoom:     &staffRoom,
		Inventory:       make(map[string]*Item),
		AvailableWeight: 20,
		CurrentEntity:   nil,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {

		if player.CurrentEntity != nil && player.CurrentEntity.Name == "sofa" {
			abandonedLanyard.Hidden = false
			sofa.SetDescription("Your fellow academy student continues to sleep on the sofa. Something tells you it's down to you to get stuff done today...")
		}

		if player.CurrentEntity != nil && player.CurrentEntity.Name == "kettle" {
			tea.Hidden = false
			kettle.SetDescription("A kettle — essential for survival, impossible to function without one nearby.");
		}

		if player.CurrentEntity != nil && player.CurrentEntity.Name == "desk" {
			firstPlate.Hidden = false
			secondPlate.Hidden = false
			thirdPlate.Hidden = false
			fourthPlate.Hidden = false
			fifthPlate.Hidden = false
			sixthPlate.Hidden = false
			desk.SetDescription("Despite the disarray, it's clear this desk sees frequent use, with just enough space left to get work done.")
		}

		for _, validInteraction := range validInteractions {
			if validInteraction.Event.Description == "get-your-lanyard" && validInteraction.Event.Triggered {
				lanyard.Hidden = false
				rosie.SetDescription("Can I help with anything else?")
			}
		}

		dishwasherLoaded := true
		for _, validInteraction := range validInteractions {
			if strings.HasSuffix(validInteraction.ItemName, "plate") && !validInteraction.Event.Triggered {
				dishwasherLoaded = false
				break
			}
		}

		if !dishwasherChallengeWon.Triggered {
			if dishwasherLoaded {
				player.TriggerEvent(dishwasherChallengeWon)
				alan.SetDescription("Ah, so you've managed to load the dishwasher! Splendid work — consider this challenge complete.\nI could have done it myself instead of writing that clever recursive function, but where's the fun in that?\nAfter all, they pay me for my intellect, not for doing the heavy lifting!\nBut I digress. You're free to proceed to the terminal room and speak with Dan for your final challenge.\nYou're doing an excellent job; keep it up!")
				dan.Hidden = false
				terminal.Hidden = false
			}
		}

		if _, ok := player.Inventory["abandoned-lanyard"]; ok {
			player.TriggerEvent(grumpyRosie)
			gameOver = true
		}

		if gameOver {
			fmt.Println("Thank you for playing!")
			break
		}

		if !introductionShown {
			clearScreen()
			fmt.Println(introduction)
			introductionShown = true
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

			

			if isAttemptingPassword {
				if remainingPasswordAttempts == 1 && input != computerPassword{
					clearScreen()
					fmt.Println("Alan's computer is locked, halting your progress in the challenge. To top it off, you've made Rosie grumpy, as she'll now have to take the computer to IT.\n\nThank you for playing!")
					break
				}
				if input == computerPassword {
					clearScreen()
					player.TriggerEvent(unlockComputer)
					computer.SetDescription("function completeTask(pile)\n   if pile == 0:\n      return 'Task Complete'\n   else:\n      completeTask(pile - 1)\n")
					alan.SetDescription("You've cracked the password! Impressive work... You should now see an open file containing a recursive function.\n\nFollow its instructions carefully, and you'll be one step closer to victory!\nBut, a word of caution: the task ahead is, well, a bit more hands-on than you might expect...")
					isAttemptingPassword = false
					desk.Hidden = false
					dishwasher.Hidden = false
				} else if input == "leave" {
					isAttemptingPassword = false
					} else {
					remainingPasswordAttempts--
					clearScreen()
					fmt.Printf("Incorrect password. Try again, or type 'leave' to stop entering the password.\n\nRemaining attempts: %d\n\n", remainingPasswordAttempts)
					computer.SetDescription(fmt.Sprintf("Alan's computer. You need the password to get in.\nRemaining attempts: %d.\nType 'leave' to stop entering the password.\n\nEnter the password:\n", remainingPasswordAttempts))
					continue
				}
			}

			if isAttemptingTerminal {
				if input == "leave" {
					isAttemptingTerminal = false
					clearScreen()
					player.Leave()
					continue
				}
			
				if !IsFirstCommand {
					if input == "cd /secret-files" {
						clearScreen()
						fmt.Println("The terminal displays:\n\n/secret-files/\n\nIt looks like you are on the right track.\nEnter the final command to win the game!\n\nType 'leave' to stop entering commands on the terminal.")
						IsFirstCommand = true
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
						gameOver = true
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

			switch command {
			case "commands":
				clearScreen()
				showCommands(ConsoleDisplay{})
			case "look":
				clearScreen()
				player.ShowRoom(ConsoleDisplay{})
			case "take":
				clearScreen()
				if len(args) > 0 {
					player.Take(args[0], ConsoleDisplay{})
				} else {
					fmt.Println("Specify an item to take.")
				}
			case "drop":
				clearScreen()
				if len(args) > 0 {
					player.Drop(args[0], ConsoleDisplay{})
				} else {
					fmt.Println("Specify an item to drop.")
				}
			case "inventory":
				clearScreen()
				player.ShowInventory(ConsoleDisplay{})
			case "approach":
				clearScreen()
				if len(args) > 0 {
					player.Approach(args[0])

					if !unlockComputer.Triggered {
						if player.CurrentEntity != nil && player.CurrentEntity.Name == "computer" {
							isAttemptingPassword = true
					}
                }
				if player.CurrentEntity != nil && player.CurrentEntity.Name == "terminal" {
					isAttemptingTerminal = true
				}
                
				} else {
					fmt.Println("Specify an entity to approach.")
				}
			case "use":
				clearScreen()
				if len(args) > 0 {
					if player.CurrentEntity == nil {
						player.Use(args[0], "unspecified_entity")
					} else {
						player.Use(args[0], player.CurrentEntity.Name)
					}
				} else {
					fmt.Println("Specify an item to use.")
				}
			case "leave":
				clearScreen()
				player.Leave()
			case "move":
				clearScreen()
				if _, ok := player.Inventory["lanyard"]; ok {
					if len(args) > 0 {
						player.Move(args[0], ConsoleDisplay{})
					} else {
						fmt.Println("Specify a direction to move (e.g., north).")
					}
				} else {
					fmt.Println("Doors are shut for you if you don't have a lanyard.")
				}
			case "map":
				clearScreen()
				player.ShowMap(ConsoleDisplay{})
			case computerPassword:
				continue
			default:
					clearScreen()
					fmt.Println("Unknown command:", command)
			}
		}
	}
}

