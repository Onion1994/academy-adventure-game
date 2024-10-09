package model

import (
	"fmt"
	"academy-adventure-game/global"
)

var plateOrder = []string{"first-plate", "second-plate", "third-plate", "fourth-plate", "fifth-plate", "sixth-plate"}
var currentPlateIndex = 0
var ValidInteractions = []*Interaction{}

func (p *Player) Move(direction string) {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
	}
	if newRoom, ok := p.CurrentRoom.Exits[direction]; ok {
		p.CurrentRoom = newRoom

		fmt.Printf("You are in %s\n", p.CurrentRoom.Name)
	} else {
		fmt.Println("You can't go that way!")
	}
}

func (p *Player) Take(itemName string) {
	item, ok := p.CurrentRoom.Items[itemName]
	switch {
	case !ok || item.Hidden:
		fmt.Printf("You can't take %s\n", itemName)
		return
	case p.AvailableWeight < item.Weight:
		fmt.Println("Weight limit reached! Please drop an item before taking more.")
		return
	case isPlate(itemName):
		if itemName == plateOrder[currentPlateIndex] {
			p.Inventory[item.Name] = item
			p.ChangeCarriedWeight(item, "increase")
			delete(p.CurrentRoom.Items, item.Name)
			currentPlateIndex++

			fmt.Printf("%s has been added to your inventory.\n", item.Name)
		} else {
			fmt.Println("As you attempt to grab the greasy plates without removing the ones stacked above them, they slip from your grasp and shatter, creating a chaotic mess.\n\nNow Rosie is very grumpy.")
			global.GameOver = true
		}

	default:
		p.Inventory[item.Name] = item
		p.ChangeCarriedWeight(item, "increase")
		delete(p.CurrentRoom.Items, item.Name)

		fmt.Printf("%s has been added to your inventory.\n", item.Name)
	}
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

func (p *Player) Drop(itemName string) {
	if item, ok := p.Inventory[itemName]; ok {
		if isPlate(itemName) {
			println("You can't just leave those plates lying around! It's time to load them into the dishwasher!")
			return
		}

		delete(p.Inventory, item.Name)
		p.ChangeCarriedWeight(item, "decrease")
		p.CurrentRoom.Items[item.Name] = item

		fmt.Printf("You dropped %s.\n", item.Name)
	} else {
		fmt.Printf("You don't have %s.\n", itemName)
	}
}

func (p *Player) ShowInventory() {
	if len(p.Inventory) == 0 {
		fmt.Printf("Your inventory is empty.\nAvailable space: %d\n", p.AvailableWeight)
		return
	}
	fmt.Printf("Available space: %d\nYour inventory contains:\n", p.AvailableWeight)
	for itemName, item := range p.Inventory {
		fmt.Printf("- %s: %s Weight: %d\n", itemName, item.Description, item.Weight)
	}
}

func (p *Player) ShowRoom() {
	fmt.Printf("You are in %s\n\n%s\n", p.CurrentRoom.Name, p.CurrentRoom.Description)

	if p.EntitiesArePresent() {
		fmt.Println("\nYou can approach:")
		for _, entity := range p.CurrentRoom.Entities {
			switch {
			case p.CurrentEntity != nil:
				if entity.Name == p.CurrentEntity.Name {
					fmt.Printf("- %s (currently approached)\n", entity.Name)
				} else if !entity.Hidden {
					fmt.Printf("- %s\n", entity.Name)
				}
			default:
				if !entity.Hidden {
					fmt.Printf("- %s\n", entity.Name)
				}
			}
		}
	}

	if p.ItemsArePresent() {
		fmt.Println("\nThe room contains:")
		for itemName, item := range p.CurrentRoom.Items {
			if !item.Hidden {
				fmt.Printf("- %s: %s Weight: %d\n", itemName, item.Description, item.Weight)
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
	if entity, ok := p.CurrentRoom.Entities[entityName]; ok && !entity.Hidden {

		p.CurrentEntity = entity
		fmt.Println(entity.Description)
	} else {
		fmt.Printf("You can't approach %s.\n", entityName)
	}
}

func (p *Player) Leave() {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
		p.ShowRoom()
	} else {
		fmt.Println("You have not approached anything. If you wish to leave the game, use the exit command.")
	}
}

func (p *Player) ShowMap() {
	for direction, exit := range p.CurrentRoom.Exits {
		fmt.Printf("%s: %s\n", direction, exit.Name)
	}
}

func (p *Player) Use(itemName string, target string) {
	if p.CurrentEntity == nil {
		fmt.Println("Approach to use an item.")
		return
	}
	if p.CurrentEntity.Name == target {
		if _, ok := p.Inventory[itemName]; ok {
			for _, interaction := range ValidInteractions {
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

func isPlate(itemName string) bool {
	for _, plate := range plateOrder {
		if itemName == plate {
			return true
		}
	}
	return false
}