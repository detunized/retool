package main

import (
	"fmt"
	"log"
	"math"

	"github.com/maja42/goval"
	t "github.com/rivo/tview"
)

var calcPane = struct {
	name string
	view *t.Form
	eval *goval.Evaluator
}{
	name: "Calc",
	eval: goval.NewEvaluator(),
}

func makeCalcPane() (t.Primitive, string) {
	p := &calcPane
	p.view = t.NewForm().
		AddInputField("Input", "", 0, nil, updateCalc).
		AddInputField("Result", "", 0, nil, nil).
		AddInputField("Hex LE", "", 8*2+7, nil, nil).
		AddInputField("Hex BE", "", 8*2+7, nil, nil)

	decoratePane(p.view.Box, p.name)
	return p.view, p.name
}

func updateCalc(text string) {
	p := &calcPane
	r, err := p.eval.Evaluate(text, nil, nil)

	result := ""
	if err == nil {
		result = fmt.Sprint(r)
	} else {
		result = err.Error()
	}

	setCalcFormField("Result", result)

	hexLE := ""
	hexBE := ""
	var n uint64

	switch v := r.(type) {
	case int:
		n = uint64(v)
	case float64:
		n = math.Float64bits(v)
	default:
		hexLE = fmt.Sprintf("%T?", v)
	}

	if hexLE == "" {
		var buf [8]byte
		buf[0] = byte(n >> 56)
		buf[1] = byte(n >> 48)
		buf[2] = byte(n >> 40)
		buf[3] = byte(n >> 32)
		buf[4] = byte(n >> 24)
		buf[5] = byte(n >> 16)
		buf[6] = byte(n >> 8)
		buf[7] = byte(n)

		hexLE = fmt.Sprintf(
			"%02x %02x %02x %02x %02x %02x %02x %02x",
			buf[7], buf[6], buf[5], buf[4], buf[3], buf[2], buf[1], buf[0],
		)
		hexBE = fmt.Sprintf(
			"%02x %02x %02x %02x %02x %02x %02x %02x",
			buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7],
		)
	}

	setCalcFormField("Hex LE", hexLE)
	setCalcFormField("Hex BE", hexBE)
}

func setCalcFormField(label string, text string) {
	i := calcPane.view.GetFormItemByLabel(label)
	if i == nil {
		log.Panicf("Form field '%s' not found", label)
	}

	if f, ok := i.(*t.InputField); ok {
		f.SetText(text)
	} else {
		log.Panicf("Form field '%s' is not an input field", label)
	}
}
