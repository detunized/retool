package main

import t "github.com/rivo/tview"

func NewFlex(direction int) *t.Flex {
	f := t.NewFlex()
	f.SetDirection(direction)
	return f
}

func NewFlexRow() *t.Flex {
	return NewFlex(t.FlexRowCSS)
}

func NewFlexColumn() *t.Flex {
	return NewFlex(t.FlexColumnCSS)
}
