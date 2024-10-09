package main

import (
	"fmt"
	"strings"
	"testing"
)

func setUpValidInteractions() {
	validInteractions = []*Interaction{
		{
			ItemName:   "key",
			EntityName: "door",
			Event:      &Event{Description: "unlock_door", Outcome: "The door unlocks with a loud click.\n", Triggered: false},
		},
		{
			ItemName:   "water",
			EntityName: "plant",
			Event:      &Event{Description: "water_plant", Outcome: "The plant looks healthier after being watered.\n", Triggered: false},
		},
	}
}

type MockDisplay struct {
	Output []string
}

func (m *MockDisplay) Show(text string) {
	m.Output = append(m.Output, text)
}


func TestPlayerCanMoveToAvailableRoom(t *testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Description: "This is room 1.", Exits: make(map[string]*Room)}
    room2 := Room{Name: "Room 2", Description: "This is room 2.", Exits: make(map[string]*Room)}
    room1.Exits["north"] = &room2
    room2.Exits["south"] = &room1

    player := Player{CurrentRoom: &room1}

	mockDisplay := &MockDisplay{}

    // Act
    player.Move("north", mockDisplay)

    // Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedRoom := "Room 2"
	expectedOutput := fmt.Sprintf("You are in %s\n", player.CurrentRoom.Name)

    if player.CurrentRoom.Name != expectedRoom {
        t.Errorf("Expected %s, got %s", expectedRoom, player.CurrentRoom.Name)
    }
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestPlayerCannotMoveToUnavailableRoom(t *testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Description: "This is room 1.", Exits: make(map[string]*Room)}
	room2 := Room{Name: "Room 2", Description: "This is room 2.", Exits: make(map[string]*Room)}
	room1.Exits["north"] = &room2
	room2.Exits["south"] = &room1

	player := Player{CurrentRoom: &room1}

	mockDisplay := MockDisplay{}

	//Act
	player.Move("east", &mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintln("You can't go that way!")
	expectedRoom := "Room 1"

    if player.CurrentRoom.Name != expectedRoom {
        t.Errorf("Expected %s, got %s", expectedRoom, player.CurrentRoom.Name)
    }
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestTakeItemIfAvailable(t *testing.T) {
	//Arrange
	room := Room{Items: make(map[string]*Item)}
	
	item := Item{Name: "Item", Description: "This is an item."}

	room.Items[item.Name] = &item
	
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item)}

	mockDisplay := &MockDisplay{}
	
	//Act
	player.Take(item.Name, mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintf("%s has been added to your inventory.\n", item.Name)

	if _, ok := player.Inventory[item.Name]; !ok {
		t.Errorf("Expected true for item present in the inventory, got false")
	}
	
	if _, ok := room.Items[item.Name]; ok {
		t.Errorf("Expected false for item missing from the room, got true")
	}

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotTakeAbsentItem(t *testing.T) {
	//Arrange
	room1 := Room{Items: make(map[string]*Item)}
	room2 := Room{Items: make(map[string]*Item)}
	
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10}

	room1.Items[item.Name] = &item
	
	player := Player{CurrentRoom: &room2, Inventory: make(map[string]*Item),  CarriedWeight: 0, AvailableWeight: 30}
	
	mockDisplay := &MockDisplay{}
	
	//Act
	player.Take(item.Name, mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintf("You can't take %s\n", item.Name)

	if _, ok := player.Inventory[item.Name]; ok {
		t.Errorf("Expected false for picking up absent item, got true")
	}
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotTakeNonexistentItem(t *testing.T) {
	//Arrange
	room2 := Room{Items: make(map[string]*Item)}

	player := Player{CurrentRoom: &room2, Inventory: make(map[string]*Item), CarriedWeight: 0, AvailableWeight: 30}
	
	mockDisplay := &MockDisplay{}
	
	//Act
	player.Take("item", mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintln("You can't take item")

	if _, ok := player.Inventory["Item"]; ok {
		t.Errorf("Expected false for picking up nonexistent item, got true")
	}
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotTakeHiddenItem(t *testing.T) {
	//Arrange
	room := Room{Items: make(map[string]*Item)}
	
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10, Hidden: true}

	room.Items[item.Name] = &item
	
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item),  CarriedWeight: 0, AvailableWeight: 30}
	
	mockDisplay := &MockDisplay{}
	
	//Act
	player.Take(item.Name, mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintf("You can't take %s\n", item.Name)

	if _, ok := player.Inventory[item.Name]; ok {
		t.Errorf("Expected false for picking up hidden item, got true")
	}
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCanDropItemFromInventoryIntoRoom(t *testing.T) {
	//Arrange
	room := Room{Items: make(map[string]*Item)}
	
	item := Item{Name: "Item", Description: "This is an item."}

	room.Items[item.Name] = &item
	
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item)}
	player.Inventory[item.Name] = &item
	
	mockDisplay := &MockDisplay{}
	
	//Act

	player.Drop(item.Name, mockDisplay)

	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintf("You dropped %s.\n\n", item.Name)

	//Assert
	if _, ok := player.Inventory[item.Name]; ok {
		t.Errorf("Expected false for item absent from the inventory, got true")
	}
	if _, ok := room.Items[item.Name]; !ok {
		t.Errorf("Expected true for item present in the room, got false")
	}

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotDropAbsentItem(t *testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Description: "This is room 1.", Exits: make(map[string]*Room), Items: make(map[string]*Item)}
	room2 := Room{Name: "Room 2", Description: "This is room 2.", Exits: make(map[string]*Room), Items: make(map[string]*Item)}
	room1.Exits["north"] = &room2
	room2.Exits["south"] = &room1

	item := Item{Name: "Item", Description: "This is an item."}

	room1.Items[item.Name] = &item
	
	player := Player{CurrentRoom: &room2, Inventory: make(map[string]*Item)}

	mockDisplay := &MockDisplay{}
	
	//Act

	player.Drop(item.Name, mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintf("You don't have %s.\n\n", item.Name)

	if _, ok := room2.Items[item.Name]; ok {
		t.Errorf("Expected false for item absent from the room, got true")
	}
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotDropNonexistentItem(t *testing.T) {
	//Arrange
	room := Room{Items: make(map[string]*Item)}
	
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item)}

	mockDisplay := &MockDisplay{}
	
	//Act
	player.Drop("Item", mockDisplay)

	//Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := "You don't have Item.\n\n"

	if _, ok := player.Inventory["Item"]; ok {
		t.Errorf("Expected false for item absent from the inventory, got true")
	}
	if _, ok := room.Items["Item"]; ok {
		t.Errorf("Expected false for item absent from the room, got true")
	}
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestShowInventory(t *testing.T) {
	// Arrange
	room := Room{Items: make(map[string]*Item)}
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10}
	room.Items[item.Name] = &item
	
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item), AvailableWeight: 30}
	player.Inventory[item.Name] = &item
	mockDisplay := &MockDisplay{}

	// Act
	player.ShowInventory(mockDisplay)

	// Assert
	expectedOutput := fmt.Sprintf("Available space: %d\nYour inventory contains:\n- %s: %s Weight: %d\n", player.AvailableWeight, item.Name, item.Description, item.Weight)

	output := strings.Join(mockDisplay.Output, "")
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}


func TestShowInventoryIsEmpty(t *testing.T) {
	// Arrange
	room := Room{Items: make(map[string]*Item)}
	
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item), AvailableWeight: 30}
	mockDisplay := &MockDisplay{}

	// Act
	player.ShowInventory(mockDisplay)
	expectedOutput := fmt.Sprintf("Your inventory is empty.\nAvailable space: %d\n", player.AvailableWeight)

	output := strings.Join(mockDisplay.Output, "")
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestShowRoom(t *testing.T) {
	// Arrange
	room := Room{Name: "Room 1", Description: "This is room 1.", Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is Entity"}
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10}
	room.Items[item.Name] = &item
	room.Entities[entity.Name] = &entity

	player := Player{CurrentRoom: &room}

	mockDisplay := &MockDisplay{}

	// Act
	player.ShowRoom(mockDisplay)

	// Assert
	output := strings.Join(mockDisplay.Output, "")
	expectedOutput := fmt.Sprintf(
		"You are in %s\n\n%s\n\nYou can approach:\n- %s\n\nThe room contains:\n- %s: %s Weight: %d",
		room.Name,
		room.Description,
		entity.Name,
		item.Name,
		item.Description,
		item.Weight,
	)

	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestShowRoomEngagedEntity(t *testing.T) {
	// Arrange
	room := Room{Name: "Room 1", Description: "This is room 1.", Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is Entity"}
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10}
	room.Items[item.Name] = &item
	room.Entities[entity.Name] = &entity

	player := Player{CurrentRoom: &room, CurrentEntity: &entity}

	mockDisplay := &MockDisplay{}

	// Act
	player.ShowRoom(mockDisplay)

	// Assert
	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintf(
		"You are in %s\n\n%s\n\nYou can approach:\n- %s (currently approached)\n\nThe room contains:\n- %s: %s Weight: %d",
		room.Name,
		room.Description,
		entity.Name,
		item.Name,
		item.Description,
		item.Weight,
	)

	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestShouldNotShowHiddenItems(t *testing.T) {
	// Arrange
	room := Room{Name: "Room 1", Description: "This is room 1.", Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is Entity", Hidden: false}
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10, Hidden: true}
	room.Items[item.Name] = &item
	room.Entities[entity.Name] = &entity

	player := Player{CurrentRoom: &room}

	mockDisplay := &MockDisplay{}

	// Act
	player.ShowRoom(mockDisplay)

	// Assert
	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintf(
		"You are in %s\n\n%s\n\nYou can approach:\n- %s\n",
		room.Name,
		room.Description,
		entity.Name,
	)

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestShouldNotShowHiddenEntities(t *testing.T) {
	// Arrange
	room := Room{Name: "Room 1", Description: "This is room 1.", Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is Entity", Hidden: true}
	item := Item{Name: "Item", Description: "This is an item.", Weight: 10, Hidden: false}
	room.Items[item.Name] = &item
	room.Entities[entity.Name] = &entity

	player := Player{CurrentRoom: &room}

	mockDisplay := &MockDisplay{}

	// Act
	player.ShowRoom(mockDisplay)

	// Assert
	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintf(
		"You are in %s\n\n%s\n\nThe room contains:\n- %s: %s Weight: %d\n",
		room.Name,
		room.Description,
		item.Name,
		item.Description,
		item.Weight,
	)

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExpectedCarriedWeight(t *testing.T) {
	//Arrange
	room := Room{Items: make(map[string]*Item)}
	item1 := Item{Name: "Item", Weight: 5}
	item2 := Item{Name: "Item 2", Weight: 10}
	item3 := Item{Name: "Item 3", Weight: 15}
	room.Items[item1.Name] = &item1
	room.Items[item2.Name] = &item2
	room.Items[item3.Name] = &item3
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item), AvailableWeight: 30}
	
	mockDisplay := &MockDisplay{}
	
	//Act
	player.Take(item1.Name, mockDisplay)
	player.Take(item2.Name, mockDisplay)
	player.Drop(item2.Name, mockDisplay)
	player.Take(item3.Name, mockDisplay)

	//Assert
	expectedOutput := 20
	output := player.CarriedWeight
	
	if output != expectedOutput {
		t.Errorf("Expected output:\n%d\nGot:\n%d", expectedOutput, output)
	}
}

func TestExpectedAvailableAndCarriedWeight(t *testing.T) {
	//Arrange
	room := Room{Items: make(map[string]*Item)}
	item1 := Item{Name: "Item", Weight: 5}
	item2 := Item{Name: "Item 2", Weight: 16}
	item3 := Item{Name: "Item 3", Weight: 15}
	room.Items[item1.Name] = &item1
	room.Items[item2.Name] = &item2
	room.Items[item3.Name] = &item3
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item), AvailableWeight: 30}

	mockDisplay := &MockDisplay{}
	
	//Act
	player.Take(item1.Name, mockDisplay)
	player.Drop(item1.Name, mockDisplay)
	player.Take(item2.Name, mockDisplay)
	player.Take(item3.Name, mockDisplay)

	//Assert
	expectedCarriedWeight := 16
	actualCarriedWeight := player.CarriedWeight
	
	expectedAvailableWeight := 14
	actualAvailableWeight := player.AvailableWeight
	
	if expectedCarriedWeight != actualCarriedWeight {
		t.Errorf("Expected output:\n%d\nGot:\n%d", expectedCarriedWeight, actualCarriedWeight)
	}
	if expectedAvailableWeight != actualAvailableWeight {
		t.Errorf("Expected output:\n%d\nGot:\n%d", expectedAvailableWeight, actualAvailableWeight)
	}
}

func TestShouldApproachPresentEntity(t* testing.T) {
	//Arrange
	room := Room{Name: "Room", Description: "This is a room.", Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is an entity"}
	room.Entities[entity.Name] = &entity
	player := Player{CurrentRoom: &room}
	mockDisplay := &MockDisplay{}

	//Act
	player.Approach(entity.Name, mockDisplay)

	//Assert
	expectedCurrentEntity :=  entity.Name
	actualCurrentEntity := player.CurrentEntity.Name

	expectedOutput := player.CurrentEntity.Description
	output := strings.Join(mockDisplay.Output, "")

	if actualCurrentEntity != expectedCurrentEntity {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedCurrentEntity, actualCurrentEntity)
	}

	if expectedOutput != output {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}

func TestShouldNotApproachAbsentEntity(t* testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Entities: make(map[string]*Entity)}
	room2 := Room{Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is an entity"}
	room2.Entities[entity.Name] = &entity
	player := Player{CurrentRoom: &room1}

	mockDisplay := &MockDisplay{}

	//Act
	player.Approach(entity.Name, mockDisplay)

	//Assert
	expectedOutput := fmt.Sprintf("You can't approach %s.\n", entity.Name)
	output := strings.Join(mockDisplay.Output, "")

	if player.CurrentEntity != nil {
        t.Errorf("Expected CurrentEntity to be nil, but got a non-nil entity")
    }
	if expectedOutput != output {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}

func TestShouldNotApproachNonexistentEntity(t* testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Entities: make(map[string]*Entity)}
	player := Player{CurrentRoom: &room1}

	mockDisplay := &MockDisplay{}

	//Act
	player.Approach("Entity", mockDisplay)

	//Assert
	expectedOutput := "You can't approach Entity.\n"
	output := strings.Join(mockDisplay.Output, "")

	if player.CurrentEntity != nil {
        t.Errorf("Expected CurrentEntity to be nil, but got a non-nil entity")
    }
	if expectedOutput != output {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}

func TestShouldNotApproachHiddenEntity(t* testing.T) {
	//Arrange
	room := Room{Name: "Room 1", Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is an entity", Hidden: true}
	room.Entities[entity.Name] = &entity
	player := Player{CurrentRoom: &room}

	mockDisplay := &MockDisplay{}

	//Act
	player.Approach(entity.Name, mockDisplay)

	//Assert
	expectedOutput := fmt.Sprintf("You can't approach %s.\n", entity.Name)
	output := strings.Join(mockDisplay.Output, "")

	if player.CurrentEntity != nil {
        t.Errorf("Expected CurrentEntity to be nil, but got a non-nil entity")
    }
	if expectedOutput != output {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}

func TestUpdateDescription(t *testing.T) {
    //Arrange
    room := &Room{Name: "Room", Description: "This is the first description"}
    item := &Item{Name: "Item", Description: "This is the first description"}
    entity := &Entity{Name: "Entity", Description: "This is the first description"}
    newDescription := "This is the second description"
    
    //Act
    room.SetDescription(newDescription)
    item.SetDescription(newDescription)
    entity.SetDescription(newDescription)

    //Assert
    if room.GetDescription() != newDescription {
        t.Errorf("Expected description:\n%s\nGot:\n%s", newDescription, room.GetDescription())
    }
    if item.GetDescription() != newDescription {
        t.Errorf("Expected description:\n%s\nGot:\n%s", newDescription, item.GetDescription())
    }
    if entity.GetDescription() != newDescription {
        t.Errorf("Expected description:\n%s\nGot:\n%s", newDescription, entity.GetDescription())
    }
}

func TestShouldDisengageEntity(t *testing.T) {
	//Arrange
	room := Room{Name: "Room", Description: "This is a room.", Entities: make(map[string]*Entity)}
	entity := Entity{Name: "Entity", Description: "This is an entity"}
	room.Entities[entity.Name] = &entity
	player := Player{CurrentRoom: &room}

	mockDisplay := &MockDisplay{}

	//Act
	player.Approach(entity.Name, mockDisplay)
	player.Leave()

	if player.CurrentEntity != nil {
        t.Errorf("Expected CurrentEntity to be nil, but got a non-nil entity")
    }
}

func TestPlayerMoveShouldDisengageEntity(t *testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Description: "This is room 1.", Exits: make(map[string]*Room), Entities: make(map[string]*Entity)}
    room2 := Room{Name: "Room 2", Description: "This is room 2.", Exits: make(map[string]*Room), Entities: make(map[string]*Entity)}

    room1.Exits["north"] = &room2
    room2.Exits["south"] = &room1
	
	entity := Entity{Name: "Entity", Description: "This is an entity"}
	room1.Entities[entity.Name] = &entity
	player := Player{CurrentRoom: &room1, CurrentEntity: nil}

	mockDisplay := &MockDisplay{}

	//Act
	player.Approach(entity.Name, mockDisplay)
	player.Move("north", mockDisplay)

	//Assert
	if player.CurrentEntity != nil {
		t.Errorf("Expected player's current entity to be nil, got %s", player.CurrentEntity.Name)
	}
}

func TestNewEngagementShouldCancelFormer(t *testing.T) {
	//Arrange
	room := Room{Name: "Room", Description: "This is a room.", Exits: make(map[string]*Room), Entities: make(map[string]*Entity)}

	entity1 := Entity{Name: "Entity", Description: "This is an entity"}
	entity2 := Entity{Name: "Entity 2", Description: "This is an entity"}

	room.Entities[entity1.Name] = &entity1
	room.Entities[entity2.Name] = &entity2
	player := Player{CurrentRoom: &room}

	mockDisplay := &MockDisplay{}

	//Act
	player.Approach(entity1.Name, mockDisplay)
	player.Approach(entity2.Name, mockDisplay)

	//Assert
	if player.CurrentEntity.Name != entity2.Name {
		t.Errorf("Expected player's current entity to be %s, got %s", entity2.Name, player.CurrentEntity.Name)
	}
}

func TestShowMap(t * testing.T) {
	//Arrange
	room1 := Room{Name: "Room 1", Description: "This is room 1.", Exits: make(map[string]*Room), Entities: make(map[string]*Entity)}
    room2 := Room{Name: "Room 2", Description: "This is room 2.", Exits: make(map[string]*Room), Entities: make(map[string]*Entity)}

    room1.Exits["north"] = &room2
    room2.Exits["south"] = &room1

	player := Player{CurrentRoom: &room1}
	
	mockDisplay := &MockDisplay{}

	// Act
	player.ShowMap(mockDisplay)

	// Assert
	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintf("north: %s\n", player.CurrentRoom.Exits["north"].Name)

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestValidUseItem(t *testing.T) {
	//Arrange
	setUpValidInteractions()
	room := Room{Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	key := Item{Name: "key", Weight: 1}
	door := Entity{Name: "door"}
	room.Entities[door.Name] = &door
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item), AvailableWeight: 30, CurrentEntity: nil}
	player.Inventory["key"] = &key
	player.CurrentEntity = &door
	mockDisplay := &MockDisplay{}
	//Act

	player.Use("key", "door", mockDisplay)

	//Assert
	if !validInteractions[0].Event.Triggered {
		t.Errorf("Expected event to be true for triggered, got false")
	}
	if _, ok := player.Inventory["key"]; ok {
		t.Errorf("Expected used item to have been removed from inventory")
	}
	if _, ok := player.CurrentRoom.Items["key"]; ok {
		t.Errorf("Expected used item to not be present in the room")
	}
	if player.AvailableWeight < 30 {
		t.Errorf("Expected inventory to return to its original state after using item")
	}
}


func TestInvalidUseItem(t *testing.T) {
	//Arrange
	setUpValidInteractions()
	room := Room{Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	key := Item{Name: "key"}
	plant := Entity{Name: "plant"}
	room.Entities[plant.Name] = &plant
	room.Items[key.Name] = &key
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item)}
	player.Inventory[key.Name] = &key
	
	mockDisplay := &MockDisplay{}

	//Act
	player.Approach("plant", mockDisplay)
	player.Use("key", "plant", mockDisplay)

	//Assert
	for _, validInteraction := range validInteractions {
		if validInteraction.Event.Triggered {
			t.Errorf("Expected event to be false for triggered, got true")
		}
	}

	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintf("You can't use %s on %s.\n", key.Name, plant.Name)

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotUseAbsentItem(t *testing.T) {
	//Arrange
	setUpValidInteractions()
	room := Room{Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	key := Item{Name: "key"}
	door := Entity{Name: "door"}
	room.Entities[door.Name] = &door
	room.Items[key.Name] = &key
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item)}
	
	mockDisplay := &MockDisplay{}

	//Act
	player.Approach("door", mockDisplay)
	player.Use("key", "door", mockDisplay)

	//Assert
	if validInteractions[0].Event.Triggered {
		t.Errorf("Expected event to be false for triggered, got true")
	}

	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintf("You don't have %s.\n", key.Name)

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestCannotUseOnAbsentEntity(t *testing.T) {
	//Arrange
	setUpValidInteractions()
	room := Room{Items: make(map[string]*Item), Entities: make(map[string]*Entity)}
	key := Item{Name: "key"}
	door := Entity{Name: "door"}
	room.Entities[door.Name] = &door
	room.Items[key.Name] = &key
	player := Player{CurrentRoom: &room, Inventory: make(map[string]*Item), CurrentEntity: nil}
	player.Inventory["key"] = &key
	mockDisplay := &MockDisplay{}
	//Act

	player.Use("key", "door", mockDisplay)

	//Assert
	if validInteractions[0].Event.Triggered {
		t.Errorf("Expected event to be false for triggered, got true")
	}
	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := "Approach to use an item.\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestShowCommands(t *testing.T) {

	mockDisplay := &MockDisplay{}

	// Act
	showCommands(mockDisplay)

	// Assert
	output := strings.Join(mockDisplay.Output, "")

	expectedOutput := fmt.Sprintln("-exit -> quits the game\n\n-commands -> shows the commands\n\n-look -> shows the content of the room.\n\n-approach <entity> -> to approach an entity\n\n-leave -> to leave an entity\n\n-inventory -> shows items in the inventory\n\n-take <item> -> to take an item into your inventory\n\n-drop <item> -> to drop an item from your inventory and move it to the current room\n\n-use <item> -> to make use of a certain item when you approach an entity\n\n-move <direction> -> to move to a different room\n\n-map -> shows the directions you can take")

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}