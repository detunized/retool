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
	label_Regular   = "Regular"
	label_HexLE     = "Hex/LE"
	label_HexBE     = "Hex/BE"
	label_HexUint64 = "Hex/uint64"
	label_AddSpaces = "Add Spaces"

	hexSeparator         = " "
	allowedHexSeparators = " "
)

type floatPane struct {
	view       *t.Form
	updating   bool
	lastInput  string
	lastSource *codec
}

type codec struct {
	name   string
	encode func(float64) (string, error)
	decode func(string) (float64, error)
	width  int
}

var instance = &floatPane{}

var codecs = []*codec{
	{
		name: label_Regular,
		encode: func(f float64) (string, error) {
			return fmt.Sprintf("%f", f), nil
		},
		decode: func(s string) (float64, error) {
			return strconv.ParseFloat(s, 64)
		},
	},
	{
		name: label_HexLE,
		encode: func(f float64) (string, error) {
			return util.Float64ToHexLE(f, hexSeparator), nil
		},
		decode: func(s string) (float64, error) {
			return util.HexToFloat64LE(s, allowedHexSeparators)
		},
	},
	{
		name: label_HexBE,
		encode: func(f float64) (string, error) {
			return util.Float64ToHexBE(f, hexSeparator), nil
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
		name: label_HexUint64,
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
