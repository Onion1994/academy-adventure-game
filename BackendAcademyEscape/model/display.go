package model

type Display interface {
	Show(text string) string
}

type ConsoleDisplay struct{}

func (c ConsoleDisplay) Show(text string) string {
	return text
}
