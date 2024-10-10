package model

type Item struct {
	Name        string
	Description string
	Weight      int
	Hidden      bool
}

func (i *Item) SetDescription(description string) {
	i.Description = description
}

func (i *Item) GetDescription() string {
	return i.Description
}
