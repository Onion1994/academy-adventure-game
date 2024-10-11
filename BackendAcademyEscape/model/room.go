package model

type Room struct {
	Name        string
	Description string
	Exits       map[string]*Room
	Items       map[string]*Item
	Entities    map[string]*Entity
}

func (r *Room) SetDescription(description string) {
	r.Description = description
}

func (r *Room) GetDescription() string {
	return r.Description
}
