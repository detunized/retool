package main

import (
	"fmt"

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
		AddInputField("Result", "", 0, nil, nil)

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

	p.view.GetFormItemByLabel("Result").(*t.InputField).SetText(result)
}
