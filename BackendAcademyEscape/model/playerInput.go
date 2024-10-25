package model

type PlayerInput struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func (p *PlayerInput) ParseInput() {

}
