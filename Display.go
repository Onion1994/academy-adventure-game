package main

import "fmt"

type Display interface {
	Show(text string)
}

type ConsoleDisplay struct{}

func (c ConsoleDisplay) Show(text string) {
	fmt.Println(text)
}

