package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strings"

	t "github.com/rivo/tview"
)

//
// UI
//

// TODO: Move to the ui package
type Pane interface {
	GetName() string
	GetView() t.Primitive
}

func DecoratePane(view *t.Box, name string) {
	view.SetBorder(true)
	view.SetTitle(" " + name + " ")
}

func SetFormField(form *t.Form, label string, text string) {
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

//
// Hex
//

func Float64ToHexLE(v float64, separator string) string {
	return Uint64ToHexLE(math.Float64bits(v), separator)
}

func Float64ToHexBE(v float64, separator string) string {
	return Uint64ToHexBE(math.Float64bits(v), separator)
}

func Uint64ToHexLE(v uint64, separator string) string {
	return uint64ToHex(v, separator, binary.LittleEndian)
}

func Uint64ToHexBE(v uint64, separator string) string {
	return uint64ToHex(v, separator, binary.BigEndian)
}

func uint64ToHex(v uint64, separator string, bo binary.ByteOrder) string {
	var buf [8]byte
	bo.PutUint64(buf[:], v)
	return BytesToHex(buf[:], separator)
}

func BytesToHex(v []byte, separator string) string {
	var s strings.Builder
	for i, x := range v {
		if i > 0 {
			s.WriteString(separator)
		}
		s.WriteString(fmt.Sprintf("%02x", x))
	}
	return s.String()
}

func HexToFloat64LE(hexStr string, strip string) (float64, error) {
	return hexToFloat64(hexStr, strip, binary.LittleEndian)
}

func HexToFloat64BE(hexStr string, strip string) (float64, error) {
	return hexToFloat64(hexStr, strip, binary.BigEndian)
}

func hexToFloat64(hexStr string, strip string, bo binary.ByteOrder) (float64, error) {
	u, err := hexToUint64(hexStr, strip, bo)
	if err != nil {
		return 0.0, nil
	}

	return math.Float64frombits(u), nil
}

func HexToUint64LE(hexStr string, strip string) (uint64, error) {
	return hexToUint64(hexStr, strip, binary.LittleEndian)
}

func HexToUint64BE(hexStr string, strip string) (uint64, error) {
	return hexToUint64(hexStr, strip, binary.BigEndian)
}

func hexToUint64(hexStr string, strip string, bo binary.ByteOrder) (uint64, error) {
	var s strings.Builder
	for _, r := range hexStr {
		if strings.ContainsRune(strip, r) {
			continue
		}

		s.WriteRune(r)
	}

	b, err := hex.DecodeString(s.String())
	if err != nil {
		return 0, err
	}

	if len(b) != 8 {
		return 0, fmt.Errorf("need 8 bytes, got '%d'", len(b))
	}

	return bo.Uint64(b), nil
}
