// Demo code for the TreeView primitive.
package main

import (
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
}

func makeHotkeyLine() t.Primitive {
	flex := NewFlexRow()

	for _, hk := range hotkeys {
		flex.AddItem(makeHotkeyButton(hk.key, hk.label), 0, 1, false)
	}

	return flex
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

//
// Panes
//

type paneInfo struct {
	name string
	view t.Primitive
}

var panes []paneInfo

func addPane(hotkey tc.Key, makePane func() (t.Primitive, string)) {
	p, name := makePane()
	panes = append(panes, paneInfo{name: name, view: p})

	// Add to hotkeys
	hotkeys = append(hotkeys, hotkeyInfo{
		key:   hotkey,
		label: name,
		handler: func() {
			showPane(p)
		},
	})
}

func showPane(pane t.Primitive) {
	contentView.Clear()
	contentView.AddItem(pane, 0, 1, true)
	application.SetFocus(pane)
}

func main() {
	application = t.NewApplication()

	// Build panes
	// TODO: Consider to make this lazy
	addPane(tc.KeyF1, makeHashPane)
	addPane(tc.KeyF2, makeHmacPane)
	addPane(tc.KeyF3, makeEncodePane)

	contentView = NewFlexRow()
	rootView = NewFlexColumn().
		AddItem(t.NewTextView().SetText("retool v0.0.1 beta =)"), 1, 0, false).
		AddItem(contentView, 0, 1, false).
		AddItem(makeHotkeyLine(), 1, 0, true)

	rootView.SetInputCapture(func(event *tc.EventKey) *tc.EventKey {
		for _, hk := range hotkeys {
			if event.Key() == hk.key {
				if hk.handler != nil {
					hk.handler()
				}
				return nil
			}
		}

		return event
	})

	application.
		SetRoot(rootView, true).
		EnableMouse(true)

	showPane(panes[2].view)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
