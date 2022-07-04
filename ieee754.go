package main

import (
	"strconv"

	t "github.com/rivo/tview"
)

const (
	ieeeLabel_Input     = "Input"
	ieeeLabel_HexLE     = "HexLE"
	ieeeLabel_HexBE     = "HexBE"
	ieeeLabel_HexUInt   = "HexUInt"
	ieeeLabel_AddSpaces = "Add Spaces"
)

var ieeePane = struct {
	name      string
	view      *t.Form
	input     string
	addSpaces bool
	updating  bool
}{
	name:      "IEEE 754",
	addSpaces: true,
}

func makeIeeePane() (t.Primitive, string) {
	p := &ieeePane
	p.view = t.NewForm().
		AddInputField(ieeeLabel_Input, "", 0, nil, updateInput).
		AddInputField(ieeeLabel_HexLE, "", 8*2+7, nil, nil).
		AddInputField(ieeeLabel_HexBE, "", 8*2+7, nil, nil).
		AddInputField(ieeeLabel_HexUInt, "", 18, nil, nil).
		AddCheckbox(ieeeLabel_AddSpaces, p.addSpaces, updateHexSpaces)

	decoratePane(p.view.Box, p.name)
	return p.view, p.name
}

func updateInput(text string) {
	ieeePane.input = text
	updateIeee()
}

func updateHexSpaces(checked bool) {
	ieeePane.addSpaces = checked
	updateIeee()
}

func updateIeee() {
	p := &ieeePane

	if p.updating {
		return
	}

	p.updating = true
	defer func() {
		p.updating = false
	}()

	hexLE := ""
	hexBE := ""
	hexUInt := ""

	f, err := strconv.ParseFloat(p.input, 64)

	if err == nil {
		sep := ""
		if p.addSpaces {
			sep = " "
		}
		hexLE = float64ToHexLE(f, sep)
		hexBE = float64ToHexBE(f, sep)
		hexUInt = "0x" + float64ToHexBE(f, "")
	}

	setIeeeFormField(ieeeLabel_HexLE, hexLE)
	setIeeeFormField(ieeeLabel_HexBE, hexBE)
	setIeeeFormField(ieeeLabel_HexUInt, hexUInt)
}

func setIeeeFormField(label string, text string) {
	setFormField(ieeePane.view, label, text)
}
