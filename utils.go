package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"

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

func removeQuotes(s string) string {
	return s[1 : len(s)-1]
}

func setFormField(form *t.Form, label string, text string) {
	i := form.GetFormItemByLabel(label)
	if i == nil {
		log.Panicf("Form field '%s' not found", label)
	}

	if f, ok := i.(*t.InputField); ok {
		f.SetText(text)
	} else {
		log.Panicf("Form field '%s' is not an input field", label)
	}
}

func float64ToHexLE(v float64, separator string) string {
	return uint64ToHexLE(math.Float64bits(v), separator)
}

func float64ToHexBE(v float64, separator string) string {
	return uint64ToHexBE(math.Float64bits(v), separator)
}

func uint64ToHexLE(v uint64, separator string) string {
	return uint64ToHex(v, separator, binary.LittleEndian)
}

func uint64ToHexBE(v uint64, separator string) string {
	return uint64ToHex(v, separator, binary.BigEndian)
}

func uint64ToHex(v uint64, separator string, bo binary.ByteOrder) string {
	var buf [8]byte
	bo.PutUint64(buf[:], v)
	return bytesToHex(buf[:], separator)
}

func bytesToHex(v []byte, separator string) string {
	var s strings.Builder
	for i, x := range v {
		if i > 0 {
			s.WriteString(separator)
		}
		s.WriteString(fmt.Sprintf("%02x", x))
	}
	return s.String()
}
