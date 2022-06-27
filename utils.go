package main

import (
	t "github.com/rivo/tview"
)

func NewFlex(direction int) *t.Flex {
	f := t.NewFlex()
	f.SetDirection(direction)
	return f
}

func decoratePane(pane *t.Box, name string) {
	pane.SetBorder(true)
	pane.SetTitle(" " + name + " ")
}

func NewFlexRow() *t.Flex {
	return NewFlex(t.FlexRowCSS)
}

func NewFlexColumn() *t.Flex {
	return NewFlex(t.FlexColumnCSS)
}

func getLabelWidth(label string) int {
	return len(label)
}

func wrapWithLabel(p t.Primitive, label string) t.Primitive {
	return NewFlexRow().
		AddItem(
			t.NewTextView().
				SetText(label).
				SetDynamicColors(true).
				SetTextColor(t.Styles.SecondaryTextColor),
			getLabelWidth(label),
			0,
			false,
		).
		AddItem(p, 0, 1, true)
}
