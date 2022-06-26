// Demo code for the TreeView primitive.
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	tc "github.com/gdamore/tcell/v2"
	t "github.com/rivo/tview"
)

var application *t.Application
var rootView *t.Flex
var contentView *t.Flex

type hotkeyInfo struct {
	key     tc.Key
	label   string
	handler func()
}

var hotkeys = []hotkeyInfo{
	{
		key:     tc.KeyEsc,
		label:   "Quit",
		handler: quit,
	},
	{
		key:     tc.KeyF1,
		label:   "Hash",
		handler: showHashScreen,
	},
	{
		key:   tc.KeyF2,
		label: "Escape",
	},
	{
		key:     tc.KeyF12,
		label:   "Key log",
		handler: showKeyLogScreen,
	},
}

func makeHotkeyLine() t.Primitive {
	flex := NewFlexRow()

	for _, hk := range hotkeys {
		flex.AddItem(makeHotkeyButton(hk.key, hk.label), 0, 1, false)
	}

	return flex
}

func makeHotkeyLine2() t.Primitive {
	grid := t.NewGrid()
	grid.SetRows(1)
	grid.SetColumns(-1, -1, -1, -1)

	for i, hk := range hotkeys {
		grid.AddItem(makeHotkeyButton(hk.key, hk.label), 0, i, 1, 1, 0, 0, false)
	}

	return grid
}

func makeHotkeyLine3() t.Primitive {
	table := t.NewTable()

	for i, hk := range hotkeys {
		table.SetCellSimple(0, i, formatHotkeyText(hk.key, hk.label))
	}

	return table
}

func makeHotkeyLine4() t.Primitive {
	text := t.NewTextView()
	text.SetDynamicColors(true)

	w := text.BatchWriter()
	defer w.Close()

	for i, hk := range hotkeys {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, formatHotkeyText(hk.key, hk.label))
	}

	return text
}

func makeHotkeyButton(key tc.Key, title string) t.Primitive {
	label := formatHotkeyText(key, title)
	return t.NewBox().
		SetDrawFunc(func(screen tc.Screen, x, y, width, height int) (int, int, int, int) {
			t.Print(screen, label, x, y+height/2, width, t.AlignLeft, tc.ColorWhite)
			return x, y, width, height
		})
}

func formatHotkeyText(key tc.Key, title string) string {
	return fmt.Sprintf("[:#7f0000:b] %s [:-:-] %s", tc.KeyNames[key], title)
}

func quit() {
	application.Stop()
}

func showScreen(screen t.Primitive) {
	contentView.Clear()
	contentView.AddItem(screen, 0, 1, true)
	application.SetFocus(screen)
}

//
// Hash screen
//

var hashView *t.Flex

type hashInfo struct {
	name string
	calc func([]byte) []byte
	view *t.TextView
}

var hashes = []*hashInfo{
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
}

func makeHashScreen() {
	hashView = NewFlexColumn().
		AddItem(t.NewInputField().SetLabel("Input: ").SetChangedFunc(updateHashes), 1, 0, true)

	for _, h := range hashes {
		h.view = t.NewTextView()
		hashView.
			AddItem(t.NewBox(), 1, 0, false).
			AddItem(wrapWithLabel(h.view, h.name+": "), 1, 0, false)
	}

	hashView.AddItem(t.NewBox(), 0, 1, false)

	updateHashes("")
}

func updateHashes(input string) {
	for _, h := range hashes {
		h.view.SetText(hex.EncodeToString(h.calc([]byte(input))))
	}
}

func showHashScreen() {
	showScreen(hashView)
}

//
// Key log screen
//

var keyboardLogScreen *t.TextView

func makeKeyLogScreen() {
	keyboardLogScreen = t.NewTextView()
}

func showKeyLogScreen() {
	showScreen(keyboardLogScreen)
}

func main() {
	application = t.NewApplication()

	keyboardLog := ""

	// Build screens
	// TODO: Consider to make this lazy
	makeHashScreen()
	makeKeyLogScreen()

	contentView = NewFlexRow()
	rootView = NewFlexColumn().
		AddItem(contentView, 0, 1, false).
		AddItem(makeHotkeyLine(), 1, 0, true)
		// AddItem(makeHotkeyLine2(), 1, 0, false).
		// AddItem(makeHotkeyLine3(), 1, 0, false).
		// AddItem(makeHotkeyLine4(), 1, 0, false)

	rootView.SetInputCapture(func(event *tc.EventKey) *tc.EventKey {
		for _, hk := range hotkeys {
			if event.Key() == hk.key {
				if hk.handler != nil {
					hk.handler()
				}
				return nil
			}
		}

		// TODO: Trim the old entries
		keyboardLog += fmt.Sprintf("%v\n", event.Name())
		keyboardLogScreen.SetText(keyboardLog).ScrollToEnd()

		return event
	})

	application.
		SetRoot(rootView, true).
		EnableMouse(true)

	showHashScreen()

	if err := application.Run(); err != nil {
		panic(err)
	}
}
