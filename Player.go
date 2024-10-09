package main

import "fmt"

type Player struct {
	CurrentRoom     *Room
	Inventory       map[string]*Item
	CurrentEntity   *Entity
	CarriedWeight   int
	AvailableWeight int
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

func isPlate(itemName string) bool {
	for _, plate := range plateOrder {
		if itemName == plate {
			return true
		}
	}
	return false
}

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
			p.AddToInventory(item, d)
			currentPlateIndex++
		} else {
			d.Show(fmt.Sprintln("As you attempt to grab the greasy plates without removing the ones stacked above them, they slip from your grasp and shatter, creating a chaotic mess.\n\nNow Rosie is very grumpy."))
			gameOver = true
		}
		
	default:
		p.AddToInventory(item, d)
	}
}

func (p *Player) AddToInventory(item *Item, d Display) {
	p.Inventory[item.Name] = item
	p.ChangeCarriedWeight(item, "increase")
	delete(p.CurrentRoom.Items, item.Name)
	d.Show(fmt.Sprintf("%s has been added to your inventory.\n", item.Name))
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
				case p.PlayerIsEngaged():
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


func (p *Player) Approach(entityName string, d Display) {
	if p.CurrentEntity != nil {
		p.CurrentEntity = nil
	}
	if entity, ok := p.CurrentRoom.Entities[entityName]; ok && !entity.Hidden{

		p.CurrentEntity = entity
		d.Show(entity.Description)
	} else {
		d.Show(fmt.Sprintf("You can't approach %s.\n", entityName))
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

func (p *Player) Use(itemName string, target string, d Display) {

	if p.CurrentEntity == nil {
		d.Show("Approach to use an item.\n")
		return
	}

	if p.CurrentEntity.Name != target {
		d.Show(fmt.Sprintf("%s not found.\n", target))
		return
	}

	if itemIsNotInInventory(p, itemName) {
		d.Show(fmt.Sprintf("You don't have %s.\n", itemName))
		return
	}
	
	for _, interaction := range validInteractions {
		if interactionIsValid(interaction, itemName, target) {
			handleInteraction(p, interaction, itemName)
			return
		} 
	}
	d.Show(fmt.Sprintf("You can't use %s on %s.\n", itemName, target))
}


func itemIsNotInInventory(p *Player, itemName string) bool {
	_, ok := p.Inventory[itemName]
    return !ok
}
func interactionIsValid(interaction *Interaction, itemName string, target string) bool {
	return interaction.ItemName == itemName && interaction.EntityName == target
}

func handleInteraction(p *Player, interaction *Interaction, itemName string) {
	p.TriggerEvent(interaction.Event)
	p.ChangeCarriedWeight(p.Inventory[itemName], "decrease")
	delete(p.Inventory, itemName)
}

func (p *Player) TriggerEvent(event *Event) {
	fmt.Println(event.Outcome)
	event.Triggered = true
}