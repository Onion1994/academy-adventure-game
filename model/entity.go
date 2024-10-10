package model

type Entity struct {
	Name string
	Description string
	Hidden bool
}


func (e *Entity) SetDescription(description string) {
	e.Description = description
}

func (e *Entity) GetDescription() string {
	return e.Description
}
