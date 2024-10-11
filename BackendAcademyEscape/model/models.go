package model

type Describable interface {
	SetDescription(description string)
	GetDescription() string
}

type GameResponse struct {
	Message  string `json:"message"`
	GameOver bool   `json:"game_over"`
}
