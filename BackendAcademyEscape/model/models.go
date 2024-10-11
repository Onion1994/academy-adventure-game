package model

type Describable interface {
	SetDescription(description string)
	GetDescription() string
}
