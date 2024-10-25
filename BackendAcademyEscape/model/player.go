package model

import (
	"academy-adventure-game/global"
	"fmt"
	"strings"
)

type Player struct {
	CurrentRoom     *Room
	Inventory       map[string]*Item
	CurrentEntity   *Entity
	CarriedWeight   int
	AvailableWeight int
}

var plateOrder = []string{"first-plate", "second-plate", "third-plate", "fourth-plate", "fifth-plate", "sixth-plate"}
var currentPlateIndex = 0
var ValidInteractions = []*Interaction{}

func (p *Player) Move(direction string, display Display) string {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
	}
	if newRoom, ok := p.CurrentRoom.Exits[direction]; ok {
		p.CurrentRoom = newRoom

		return display.Show(fmt.Sprintf("You are in %s\n", p.CurrentRoom.Name))
	} else {
		return display.Show("You can't go that way!\n")
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

func (p *Player) Take(itemName string, display Display) string {
	item, ok := p.CurrentRoom.Items[itemName]
	switch {
	case !ok || item.Hidden:
		return display.Show(fmt.Sprintf("You can't take %s\n", itemName))

	case p.AvailableWeight < item.Weight:
		return display.Show(fmt.Sprintln("Weight limit reached! Please drop an item before taking more."))

	case isPlate(itemName):
		if itemName == plateOrder[currentPlateIndex] {
			currentPlateIndex++
			return p.AddToInventory(item, display)

		} else {
			global.GameOver = true
			return display.Show(fmt.Sprintln("As you attempt to grab the greasy plates without removing the ones stacked above them, they slip from your grasp and shatter, creating a chaotic mess.\n\nNow Rosie is very grumpy."))
		}

	default:
		return p.AddToInventory(item, display)
	}
}

func (p *Player) AddToInventory(item *Item, display Display) string {
	p.Inventory[item.Name] = item
	p.ChangeCarriedWeight(item, "increase")
	delete(p.CurrentRoom.Items, item.Name)
	if item.Name == "abandoned-lanyard" {
		return "Rosie caught you in the act of swiping a lanyard from a fellow student.\nYou have made Rosie grumpy and you've lost the game."
	}
	return display.Show(fmt.Sprintf("%s has been added to your inventory.\n", item.Name))
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

func (p *Player) Drop(itemName string, display Display) string {
	if item, ok := p.Inventory[itemName]; ok {
		if isPlate(itemName) {
			return display.Show("You can't just leave those plates lying around! It's time to load them into the dishwasher!")
		}

		delete(p.Inventory, item.Name)
		p.ChangeCarriedWeight(item, "decrease")
		p.CurrentRoom.Items[item.Name] = item

		return display.Show(fmt.Sprintf("You dropped %s.\n\n", item.Name))
	} else {
		return display.Show(fmt.Sprintf("You don't have %s.\n\n", itemName))
	}
}

func (p *Player) ShowInventory(display Display) string {
	if len(p.Inventory) == 0 {
		return display.Show(fmt.Sprintf("Your inventory is empty.\nAvailable space: %d\n", p.AvailableWeight))
	}
	var itemArray []string
	itemArray = append(itemArray, (fmt.Sprintf("Available space: %d\nYour inventory contains:\n", p.AvailableWeight)))
	for itemName, item := range p.Inventory {
		itemArray = append(itemArray, (fmt.Sprintf("- %s: %s Weight: %d\n", itemName, item.Description, item.Weight)))
	}
	return display.Show(strings.Join(itemArray, ""))
}

func (p *Player) ShowRoom(display Display) string {
	var returnValue []string
	returnValue = append(returnValue, display.Show(fmt.Sprintf("You are in %s\n\n%s\n", p.CurrentRoom.Name, p.CurrentRoom.Description)))

	if p.EntitiesArePresent() {
		returnValue = append(returnValue, display.Show("\nYou can approach:\n"))
		for _, entity := range p.CurrentRoom.Entities {
			switch {
			case p.PlayerIsEngaged():
				if entity.Name == p.CurrentEntity.Name {
					returnValue = append(returnValue, display.Show(fmt.Sprintf("- %s (currently approached)\n", entity.Name)))
				} else if !entity.Hidden {
					returnValue = append(returnValue, display.Show(fmt.Sprintf("- %s\n", entity.Name)))
				}
			default:
				if !entity.Hidden {
					returnValue = append(returnValue, display.Show(fmt.Sprintf("- %s\n", entity.Name)))
				}
			}
		}
	}

	if p.ItemsArePresent() {
		returnValue = append(returnValue, display.Show("\nThe room contains:"))
		for itemName, item := range p.CurrentRoom.Items {
			if !item.Hidden {
				returnValue = append(returnValue, display.Show(fmt.Sprintf("\n- %s: %s Weight: %d\n", itemName, item.Description, item.Weight)))
			}
		}
	}
	return strings.Join(returnValue, "")
}

func (p *Player) PlayerIsEngaged() bool {
	return p.CurrentEntity != nil
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

func (p *Player) Approach(entityName string, display Display) string {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
	}
	if entity, ok := p.CurrentRoom.Entities[entityName]; ok && !entity.Hidden {

		p.CurrentEntity = entity
		return display.Show(entity.Description)
	} else {
		return display.Show(fmt.Sprintf("You can't approach %s.\n", entityName))
	}
}

func (p *Player) Leave() string {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
		return p.ShowRoom(ConsoleDisplay{})
	} else {
		return ("You have not approached anything. If you wish to leave the game, use the exit command.")
	}
}

func (p *Player) ShowMap(display Display) string {
	var returnValue []string
	for direction, exit := range p.CurrentRoom.Exits {
		returnValue = append(returnValue, (fmt.Sprintf("%s: %s\n", direction, exit.Name)))
	}
	return display.Show(strings.Join(returnValue, ""))
}

func (p *Player) Use(itemName string, target string, display Display) string {

	if p.CurrentEntity == nil {
		return display.Show("Approach to use an item.\n")

	}

	if p.CurrentEntity.Name != target {
		return display.Show(fmt.Sprintf("%s not found.\n", target))
	}

	if itemIsNotInInventory(p, itemName) {
		return display.Show(fmt.Sprintf("You don't have %s.\n", itemName))

	}

	for _, interaction := range ValidInteractions {
		if interactionIsValid(interaction, itemName, target) {

			return handleInteraction(p, interaction, itemName)
		}
	}
	return display.Show(fmt.Sprintf("You can't use %s on %s.\n", itemName, target))
}

func itemIsNotInInventory(player *Player, itemName string) bool {
	_, ok := player.Inventory[itemName]
	return !ok
}
func interactionIsValid(interaction *Interaction, itemName string, target string) bool {
	return interaction.ItemName == itemName && interaction.EntityName == target
}

func handleInteraction(player *Player, interaction *Interaction, itemName string) string {
	
	player.ChangeCarriedWeight(player.Inventory[itemName], "decrease")
	delete(player.Inventory, itemName)
	return player.TriggerEvent(interaction.Event)
}

func (p *Player) TriggerEvent(event *Event) string {
	
	event.Triggered = true
	return event.Outcome
}
