package shape

import (
	"encoding/binary"
	"fmt"
	"github.com/detunized/retool/util"
	t "github.com/rivo/tview"
	"math"
	"strconv"
	"time"
)

const (
	labelHex    = "Hex"
	labelInt    = "Int"
	labelString = "String"
	labelHash   = "Hash"

	hexSeparator         = " "
	allowedHexSeparators = " "
)

type shapePane struct {
	view        *t.Form
	varIntGroup *util.CodecGroup[int64]
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

	p.varIntGroup = &util.CodecGroup[int64]{
		View: p.view,
		Codecs: []util.Codec[int64]{
			{
				Name:   labelHex,
				Encode: encodeToHex,
				Decode: decodeFromHex,
			},
			{
				Name:   labelInt,
				Encode: encodeToInt,
				Decode: decodeFromInt,
			},
		},
		SetError: func(message string) {
			util.DecoratePane(p.view.Box, "Error: "+message)
		},
		ClearError: func() {
			util.DecoratePane(p.view.Box, p.GetName())
		},
	}

	p.varIntGroup.InitView()
	p.varIntGroup.ClearError()

	hashTimer := time.NewTimer(time.Hour * 1_000_000)
	hashString := ""

	counter := 0
	go func() {
		for {
			<-hashTimer.C
			counter++
			h := CalcStringHashInt([]byte(hashString))
			s := fmt.Sprintf("[%v] S: %v | U: %v | HEX: %x", counter, int32(h), h, h)
			util.SetFormField(p.view, labelHash, s)
		}
	}()

	p.view.AddInputField(labelString, "", 0, nil, func(text string) {
		hashTimer.Stop()
		hashTimer.Reset(1000 * time.Millisecond)
		hashString = text
	})

	p.view.AddInputField(labelHash, "", 0, func(textToCheck string, lastChar rune) bool {
		return false
	}, nil)

	return instance
}

func encodeToHex(n int64) (string, error) {
	bytes := encodeVarInt(n)
	return util.BytesToHex(bytes, hexSeparator), nil
}

func decodeFromHex(s string) (int64, error) {
	bytes, err := util.HexToBytes(s, allowedHexSeparators)
	if err != nil {
		return 0, err
	}

	return decodeVarInt(bytes)
}

func encodeToInt(n int64) (string, error) {
	return strconv.FormatInt(n, 10), nil
}

func decodeFromInt(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

func decodeVarInt(bytes []byte) (int64, error) {
	if len(bytes) == 0 {
		return 0, fmt.Errorf("bytes are empty")
	}

	// Double
	if bytes[0] == 0x80 {
		if len(bytes) != 9 {
			return 0, fmt.Errorf("must be exactly 9 bytes, got %v", len(bytes))
		}
		u := binary.LittleEndian.Uint64(bytes[1:9])
		f := math.Float64frombits(u)
		return int64(f), nil
	}

	negative := (bytes[0] & 64) != 0
	n := int64(bytes[0] & 31)
	index := 0

	if (bytes[0] & 32) != 0 {
		shift := 5
		for index = 1; index < len(bytes); index++ {
			n |= int64(bytes[index]&127) << shift
			shift += 7
			if bytes[index] < 128 {
				break
			}
		}
	}

	if negative {
		n = -n
	}

	if index != len(bytes)-1 {
		return 0, fmt.Errorf("%v too many bytes", len(bytes)-index-1)
	}

	return n, nil
}

func encodeVarInt(n64 int64) []byte {
	bytes := make([]byte, 1)

	// Store as double
	if n64 != int64(int32(n64)) {
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
