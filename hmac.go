package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"strings"

	tc "github.com/gdamore/tcell/v2"
	t "github.com/rivo/tview"
)

type hmacInfo struct {
	name string
	calc func(m, k []byte) []byte
	view *t.TextView
}

var hmacPane = struct {
	view        *t.Flex
	messageView *t.InputField
	keyView     *t.InputField
	tab         []t.Primitive
	message     string
	key         string
	hmacs       []*hmacInfo
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
	hmacPane.messageView = t.NewInputField().
		SetLabel("Message: ").
		SetChangedFunc(updateMessage)
	hmacPane.keyView = t.NewInputField().
		SetLabel("    Key: ").
		SetChangedFunc(updateKey)

	hmacPane.tab = append(hmacPane.tab, hmacPane.messageView, hmacPane.keyView)

	hmacPane.view = NewFlexColumn().
		AddItem(hmacPane.messageView, 1, 0, true).
		AddItem(hmacPane.keyView, 1, 0, true).
		AddItem(t.NewBox(), 1, 0, false)

	// TODO: Move the Tab/Shift+Tab handling into some shared code
	hmacPane.view.SetInputCapture(func(event *tc.EventKey) *tc.EventKey {
		dir := 0
		if event.Key() == tc.KeyTab {
			dir = 1
		} else if event.Key() == tc.KeyBacktab {
			dir = -1
		}

		if dir != 0 {
			index := 0

			for i, t := range hmacPane.tab {
				if t.HasFocus() {
					index = (i + dir + len(hmacPane.tab)) % len(hmacPane.tab)
					break
				}
			}

			application.SetFocus(hmacPane.tab[index])
			return nil
		}

		return event
	})

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
