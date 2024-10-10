package model

type Describable interface {
	SetDescription(description string)
	GetDescription() string
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
