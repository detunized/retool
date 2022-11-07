package shape

import (
	"fmt"
	"github.com/detunized/retool/util"
	t "github.com/rivo/tview"
	"strconv"
)

const (
	labelHex = "Hex"
	labelInt = "Int"

	hexSeparator         = " "
	allowedHexSeparators = " "
)

type shapePane struct {
	view     *t.Form
	updating bool
}

var instance = &shapePane{}

func (p *shapePane) GetName() string {
	return "Shape"
}

func (p *shapePane) GetView() t.Primitive {
	return p.view
}

func MakePane() util.Pane {
	p := instance
	p.view = t.NewForm()

	p.view.AddInputField(labelHex, "7c ad dd f6 03", 0, nil, func(text string) {
		p.updateInput(labelHex, text)
	})

	p.view.AddInputField(labelInt, "", 0, nil, func(text string) {
		p.updateInput(labelInt, text)
	})

	util.DecoratePane(p.view.Box, p.GetName())

	return instance
}

func (p *shapePane) updateInput(sourceName string, input string) {
	if p.updating {
		return
	}

	p.updating = true
	defer func() {
		p.updating = false
	}()

	switch sourceName {
	case labelHex:
		b, err := util.HexToBytes(input, allowedHexSeparators)
		if err == nil {
			varInt, err := decodeVarInt(b)
			if err == nil {
				util.SetFormField(p.view, labelInt, strconv.Itoa(varInt))
			}
		}
		break
	case labelInt:
		i, err := strconv.ParseInt(input, 0, 32)
		if err == nil {
			util.SetFormField(p.view, labelHex, util.BytesToHex(encodeVarInt(int(i)), hexSeparator))
		}
		break
	}
}

func decodeVarInt(bytes []byte) (int, error) {
	if len(bytes) == 0 || len(bytes) > 5 {
		return 0, fmt.Errorf("invalid length: %v", len(bytes))
	}

	negative := (bytes[0] & 64) != 0
	n := int(bytes[0] & 31)

	if (bytes[0] & 32) != 0 {
		shift := 5
		for i := 1; i < len(bytes); i++ {
			n |= int(bytes[i]&127) << shift
			shift += 7
			if bytes[i] < 128 {
				break
			}
		}
	}

	if negative {
		n = -n
	}

	return n, nil
}

func encodeVarInt(n int) []byte {
	bytes := make([]byte, 1)

	if n < 0 {
		bytes[0] |= 64
		n = -n
	}

	bytes[0] |= byte(n & 31)
	if n > 31 {
		bytes[0] |= 32
	}

	n = int(uint(n) >> 5)
	if n > 0 {
		for {
			bytes = append(bytes, byte(n&127))
			n = int(uint(n) >> 7)
			if n == 0 {
				break
			}
			bytes[len(bytes)-1] |= 128
		}
	}

	return bytes
}
