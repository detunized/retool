package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	t "github.com/rivo/tview"
)

var hashScreen = struct {
	view   *t.Flex
	hashes []*hashInfo
}{
	hashes: []*hashInfo{
		{
			name: "MD5",
			calc: func(b []byte) []byte {
				h := md5.Sum(b)
				return h[:]
			},
		},
		{
			name: "SHA1",
			calc: func(b []byte) []byte {
				h := sha1.Sum(b)
				return h[:]
			},
		},
		{
			name: "SHA256",
			calc: func(b []byte) []byte {
				h := sha256.New()
				return h.Sum(b)
			},
		},
	},
}

type hashInfo struct {
	name string
	calc func([]byte) []byte
	view *t.TextView
}

func makeHashScreen() t.Primitive {
	hashScreen.view = NewFlexColumn().
		AddItem(t.NewInputField().SetLabel("Input: ").SetChangedFunc(updateHashes), 1, 0, true)

	for _, h := range hashScreen.hashes {
		h.view = t.NewTextView()
		hashScreen.view.
			AddItem(t.NewBox(), 1, 0, false).
			AddItem(wrapWithLabel(h.view, h.name+": "), 1, 0, false)
	}

	hashScreen.view.AddItem(t.NewBox(), 0, 1, false)

	updateHashes("")

	return hashScreen.view
}

func updateHashes(input string) {
	for _, h := range hashScreen.hashes {
		h.view.SetText(hex.EncodeToString(h.calc([]byte(input))))
	}
}
