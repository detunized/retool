package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"strings"

	t "github.com/rivo/tview"
)

type hmacInfo struct {
	name string
	calc func(m, k []byte) []byte
	view *t.TextView
}

var hmacPane = struct {
	view    *t.Flex
	message string
	key     string
	hmacs   []*hmacInfo
}{
	message: "",
	key:     "",
	hmacs: []*hmacInfo{
		{
			name: "MD5",
			calc: func(m, k []byte) []byte {
				return hmacWith(m, k, md5.New)
			},
		},
		{
			name: "SHA1",
			calc: func(m, k []byte) []byte {
				return hmacWith(m, k, sha1.New)
			},
		},
		{
			name: "SHA256",
			calc: func(m, k []byte) []byte {
				return hmacWith(m, k, sha256.New)
			},
		},
	},
}

func hmacWith(m, k []byte, hash func() hash.Hash) []byte {
	h := hmac.New(hash, k)
	h.Write(m)
	return h.Sum(nil)
}

func makeHmacPane() t.Primitive {
	hmacPane.view = NewFlexColumn().
		AddItem(t.NewInputField().SetLabel("Message: ").SetChangedFunc(updateMessage), 1, 0, true).
		AddItem(t.NewInputField().SetLabel("    Key: ").SetChangedFunc(updateKey), 1, 0, true).
		AddItem(t.NewBox(), 1, 0, false)

	// TODO: Merge code with hash.go

	maxNameLength := 0
	for _, h := range hmacPane.hmacs {
		if len(h.name) > maxNameLength {
			maxNameLength = len(h.name)
		}
	}

	for _, h := range hmacPane.hmacs {
		h.view = t.NewTextView()
		name := h.name
		if len(name) < maxNameLength {
			name = strings.Repeat(" ", maxNameLength-len(name)) + name
		}
		hmacPane.view.
			AddItem(wrapWithLabel(h.view, name+": "), 1, 0, false)
	}

	hmacPane.view.AddItem(t.NewBox(), 0, 1, false)

	updateHmacs()

	return hmacPane.view
}

func updateMessage(message string) {
	hmacPane.message = message
	updateHmacs()
}

func updateKey(key string) {
	hmacPane.key = key
	updateHmacs()
}

func updateHmacs() {
	for _, hi := range hmacPane.hmacs {
		h := hi.calc([]byte(hmacPane.message), []byte(hmacPane.key))
		hi.view.SetText(hex.EncodeToString(h))
	}
}
