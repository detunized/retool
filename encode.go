package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"html"
	"net/url"
	"strconv"
	"strings"

	tc "github.com/gdamore/tcell/v2"
	t "github.com/rivo/tview"
)

type encoderInfo struct {
	name   string
	encode func([]byte) (string, error)
	decode func(string) ([]byte, error)
	view   *t.InputField
}

var encodePane = struct {
	name     string
	view     *t.Flex
	encoders []*encoderInfo
	updating bool
}{
	name: "Encode",
	encoders: []*encoderInfo{
		{
			name: "Text",
			encode: func(b []byte) (string, error) {
				return string(b), nil
			},
			decode: func(s string) ([]byte, error) {
				return []byte(s), nil
			},
		},
		{
			name: "Hex",
			encode: func(b []byte) (string, error) {
				return hex.EncodeToString(b), nil
			},
			decode: func(s string) ([]byte, error) {
				return hex.DecodeString(s)
			},
		},
		{
			name: "Base64/ANY",
			encode: func(b []byte) (string, error) {
				return base64.StdEncoding.EncodeToString(b), nil
			},
			decode: func(s string) ([]byte, error) {
				s = strings.ReplaceAll(s, "+", "-")
				s = strings.ReplaceAll(s, "/", "_")
				s = strings.TrimRight(s, "=")
				return base64.RawURLEncoding.DecodeString(s)
			},
		},
		{
			name: "Base64/STD",
			encode: func(b []byte) (string, error) {
				return base64.StdEncoding.EncodeToString(b), nil
			},
			decode: func(s string) ([]byte, error) {
				return base64.StdEncoding.DecodeString(s)
			},
		},
		{
			name: "Base64/URL",
			encode: func(b []byte) (string, error) {
				return base64.URLEncoding.EncodeToString(b), nil
			},
			decode: func(s string) ([]byte, error) {
				return base64.URLEncoding.DecodeString(s)
			},
		},
		{
			name: "Base32/STD",
			encode: func(b []byte) (string, error) {
				return base32.StdEncoding.EncodeToString(b), nil
			},
			decode: func(s string) ([]byte, error) {
				return base32.StdEncoding.DecodeString(s)
			},
		},
		{
			name: "Base32/HEX",
			encode: func(b []byte) (string, error) {
				return base32.HexEncoding.EncodeToString(b), nil
			},
			decode: func(s string) ([]byte, error) {
				return base32.HexEncoding.DecodeString(s)
			},
		},
		{
			name: "Escape/URL",
			encode: func(b []byte) (string, error) {
				return url.QueryEscape(string(b)), nil
			},
			decode: func(s string) ([]byte, error) {
				raw, err := url.QueryUnescape(s)
				return []byte(raw), err
			},
		},
		{
			name: "Escape/C",
			encode: func(b []byte) (string, error) {
				return removeQuotes(strconv.Quote(string(b))), nil
			},
			decode: func(s string) ([]byte, error) {
				raw, err := strconv.Unquote(`"` + s + `"`)
				return []byte(raw), err
			},
		},
		{
			name: "Escape/HTML",
			encode: func(b []byte) (string, error) {
				return html.EscapeString(string(b)), nil
			},
			decode: func(s string) ([]byte, error) {
				return []byte(html.UnescapeString(s)), nil
			},
		},
		{
			name: "Punycode",
			encode: func(b []byte) (string, error) {
				return idna.ToASCII(string(b))
			},
			decode: func(s string) ([]byte, error) {
				raw, err := idna.ToUnicode(s)
				return []byte(raw), err
			},
		},
	},
}

func makeEncodePane() (t.Primitive, string) {
	p := &encodePane
	p.view = NewFlexColumn()
	decoratePane(p.view.Box, p.name)

	for _, ei := range p.encoders {
		ei := ei
		ei.view = t.NewInputField().
			SetLabel(ei.name + ": ").
			SetChangedFunc(func(text string) {
				if p.updating {
					return
				}

				b, err := ei.decode(text)
				updateEncoded(ei.name, b, err)
			})

		p.view.AddItem(ei.view, 1, 0, true)
	}

	// Add clear box
	p.view.AddItem(t.NewBox(), 0, 1, false)

	return p.view, p.name
}

func updateEncoded(source string, raw []byte, err error) {
	p := &encodePane

	if p.updating {
		return
	}

	p.updating = true
	defer func() {
		p.updating = false
	}()

	for _, ei := range p.encoders {
		if ei.name == source {
			if err == nil {
				ei.view.SetFieldBackgroundColor(t.Styles.ContrastBackgroundColor)
			} else {
				ei.view.SetFieldBackgroundColor(tc.ColorCrimson)
			}
			continue
		} else {
			text, encodeErr := ei.encode(raw)
			if encodeErr == nil {
				ei.view.SetText(text)
				ei.view.SetFieldBackgroundColor(t.Styles.ContrastBackgroundColor)
			} else {
				ei.view.SetText(text)
				ei.view.SetFieldBackgroundColor(tc.ColorCrimson)
			}
		}
	}
}
