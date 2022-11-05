package float

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/detunized/retool/util"
	t "github.com/rivo/tview"
)

const (
	labelRegular   = "Regular"
	labelHexLE     = "Hex/LE"
	labelHexBE     = "Hex/BE"
	labelHexUint64 = "Hex/uint64"
	labelAddSpaces = "Add Spaces"

	hexSeparator         = " "
	allowedHexSeparators = " "
)

type floatPane struct {
	view       *t.Form
	updating   bool
	lastInput  string
	lastSource *codec
	addSpaces  bool
}

type codec struct {
	name   string
	encode func(float64) (string, error)
	decode func(string) (float64, error)
	width  int
}

var instance = &floatPane{
	addSpaces: true,
}

var codecs = []*codec{
	{
		name: labelRegular,
		encode: func(f float64) (string, error) {
			return fmt.Sprintf("%f", f), nil
		},
		decode: func(s string) (float64, error) {
			return strconv.ParseFloat(s, 64)
		},
	},
	{
		name: labelHexLE,
		encode: func(f float64) (string, error) {
			return util.Float64ToHexLE(f, instance.getCurrentSeparator()), nil
		},
		decode: func(s string) (float64, error) {
			return util.HexToFloat64LE(s, allowedHexSeparators)
		},
	},
	{
		name: labelHexBE,
		encode: func(f float64) (string, error) {
			return util.Float64ToHexBE(f, instance.getCurrentSeparator()), nil
		},
		decode: func(s string) (float64, error) {
			u, err := util.HexToUint64LE(s, allowedHexSeparators)
			if err != nil {
				return 0.0, err
			}
			return math.Float64frombits(u), nil
		},
		width: 8*2 + 7,
	},
	{
		name: labelHexUint64,
		encode: func(f float64) (string, error) {
			return "0x" + util.Float64ToHexBE(f, ""), nil
		},
		decode: func(s string) (float64, error) {
			u, err := strconv.ParseUint(strings.TrimPrefix(strings.Trim(s, " "), "0x"), 16, 64)
			if err != nil {
				return 0.0, err
			}

			return math.Float64frombits(u), nil
		},
		width: 2 + 8*2,
	},
}

func (p *floatPane) GetName() string {
	return "Float"
}

func (p *floatPane) GetView() t.Primitive {
	return p.view
}

func MakePane() util.Pane {
	p := instance
	p.view = t.NewForm()

	for _, c := range codecs {
		c := c
		p.view.AddInputField(c.name, "", 0, nil, func(s string) {
			if p.updating {
				return
			}

			p.lastInput = s
			p.lastSource = c

			p.updateInput()
		})
	}

	p.view.AddCheckbox(labelAddSpaces, p.addSpaces, func(checked bool) {
		p.addSpaces = checked
		p.updateInput()
	})

	util.DecoratePane(p.view.Box, p.GetName())

	return instance
}

func (p *floatPane) updateInput() {
	if p.updating {
		return
	}

	p.updating = true
	defer func() {
		p.updating = false
	}()

	text := p.lastInput
	source := p.lastSource

	f, err := source.decode(text)
	if err == nil {
		for _, c := range codecs {
			if c.name == source.name {
				continue
			}

			s, _ := c.encode(f)
			util.SetFormField(p.view, c.name, s)
		}
		p.view.SetTitle(p.GetName())
	} else {
		p.view.SetTitle("Error")
	}
}

func (p *floatPane) getCurrentSeparator() string {
	if instance.addSpaces {
		return hexSeparator
	}
	return ""
}
