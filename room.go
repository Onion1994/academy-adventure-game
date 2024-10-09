package main

func (r *Room) SetDescription(description string) {
	r.Description = description
}

func (r *Room) GetDescription() string {
	return r.Description
}
