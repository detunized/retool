package util

import t "github.com/rivo/tview"

type CodecGroup[T any] struct {
	View       *t.Form
	Codecs     []Codec[T]
	SetError   func(string)
	ClearError func()

	updating        bool
	lastInput       string
	lastSourceIndex int
}

type Codec[T any] struct {
	Name   string
	Encode func(T) (string, error)
	Decode func(string) (T, error)
	Width  int
}

func (cg *CodecGroup[T]) InitView() {
	for i, c := range cg.Codecs {
		i := i
		c := c
		cg.View.AddInputField(c.Name, "", 0, nil, func(s string) {
			if cg.updating {
				return
			}

			cg.lastInput = s
			cg.lastSourceIndex = i

			cg.update()
		})
	}
}

func (cg *CodecGroup[T]) update() {
	if cg.updating {
		return
	}

	cg.updating = true
	defer func() {
		cg.updating = false
	}()

	text := cg.lastInput
	source := cg.Codecs[cg.lastSourceIndex]

	f, err := source.Decode(text)
	if err == nil {
		if cg.ClearError != nil {
			cg.ClearError()
		}

		for i, c := range cg.Codecs {
			if i == cg.lastSourceIndex {
				continue
			}

			s, err := c.Encode(f)
			if err != nil {
				s = "???"
			}
			SetFormField(cg.View, c.Name, s)
		}
	} else if cg.SetError != nil {
		cg.SetError(err.Error())
	}
}
