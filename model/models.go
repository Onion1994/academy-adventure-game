package model

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

type Room struct {
	Name string
	Description string
	Exits map[string]*Room
	Items map[string]*Item
	Entities map[string]*Entity
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
