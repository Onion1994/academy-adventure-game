package model

func (i *Item) SetDescription(description string) {
	i.Description = description
}

func (i *Item) GetDescription() string {
	return i.Description
}
