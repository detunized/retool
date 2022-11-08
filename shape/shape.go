package shape

import (
	"encoding/binary"
	"fmt"
	"github.com/detunized/retool/util"
	t "github.com/rivo/tview"
	"math"
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

	p.view.AddInputField(labelHex, "", 0, nil, func(text string) {
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
		{
			b, err := util.HexToBytes(input, allowedHexSeparators)
			s := ""
			if err == nil {
				varInt, err := decodeVarInt(b)
				if err == nil {
					s = strconv.FormatInt(varInt, 10)
				} else {
					s = err.Error()
				}
			} else {
				s = err.Error()
			}
			util.SetFormField(p.view, labelInt, s)
		}
		break
	case labelInt:
		{
			i, err := strconv.ParseInt(input, 0, 64)
			s := ""
			if err == nil {
				s = util.BytesToHex(encodeVarInt(i), hexSeparator)
			} else {
				s = err.Error()
			}
			util.SetFormField(p.view, labelHex, s)
		}
		break
	}
}

func decodeVarInt(bytes []byte) (int64, error) {
	if len(bytes) == 0 || len(bytes) > 9 {
		return 0, fmt.Errorf("invalid length: %v", len(bytes))
	}

	// Double
	if bytes[0] == 0x80 {
		u := binary.LittleEndian.Uint64(bytes[1:9])
		f := math.Float64frombits(u)
		return int64(f), nil
	}

	negative := (bytes[0] & 64) != 0
	n := int64(bytes[0] & 31)

	if (bytes[0] & 32) != 0 {
		shift := 5
		for i := 1; i < len(bytes); i++ {
			n |= int64(bytes[i]&127) << shift
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

func encodeVarInt(n64 int64) []byte {
	bytes := make([]byte, 1)

	// Store as double
	if (uint64(n64) & 0xffff_ffff_0000_0000) != 0 {
		bytes := make([]byte, 9)
		bytes[0] = 0x80

		f64 := math.Float64bits(float64(n64))
		binary.LittleEndian.PutUint64(bytes[1:9], f64)

		return bytes
	}

	n := int32(n64)
	if n < 0 {
		bytes[0] |= 64
		n = -n
	}

	bytes[0] |= byte(n & 31)
	if n > 31 {
		bytes[0] |= 32
	}

	n = int32(uint32(n) >> 5)
	if n > 0 {
		for {
			bytes = append(bytes, byte(n&127))
			n = int32(uint32(n) >> 7)
			if n == 0 {
				break
			}
			bytes[len(bytes)-1] |= 128
		}
	}

	return bytes
}
