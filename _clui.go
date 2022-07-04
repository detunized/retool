// VladimirMarkelov/clui experiment

package main

func mainClui() {
	ui.InitLibrary()

	w := ui.AddWindow(0, 0, 10, 7, "Hello World!")
	//w.SetSizable(false)
	w.SetMaximized(true)
	w.SetPack(ui.Vertical)
	w.SetPaddings(1, 1)
	w.SetGaps(0, 0)
	//w.SetBorder(ui.BorderNone)

	pane := ui.CreateFrame(w, ui.AutoSize, ui.AutoSize, ui.BorderThick, 1)
	pane.SetPack(ui.Vertical)
	pane.SetTitle("| Calc |")
	pane.SetAlign(ui.AlignCenter)

	e1 := addEditWithLabel(pane, "Input", true)
	e2 := addEditWithLabel(pane, "Result", false)

	e1.OnChange(func(e ui.Event) {
		e2.SetTitle("r:" + e.Msg)
	})

	pane2 := ui.CreateFrame(w, ui.AutoSize, ui.AutoSize, ui.BorderThick, 1)
	pane2.SetPack(ui.Vertical)
	pane2.SetTitle("| Hash |")
	pane2.SetAlign(ui.AlignCenter)
	pane2.SetVisible(false)

	addEditWithLabel(pane2, "Input2", true)
	addEditWithLabel(pane2, "Result2", false)
	addEditWithLabel(pane2, "Result2", false)
	addEditWithLabel(pane2, "Result2", false)
	addEditWithLabel(pane2, "Result2", false)

	f := ui.CreateFrame(w, ui.AutoSize, 1, ui.BorderNone, ui.Fixed)
	f.SetPack(ui.Horizontal)
	f.SetGaps(3, 0)
	addHotkeyLabel(f, "Esc", "Quit")
	addHotkeyLabel(f, "F1", "Hash")
	addHotkeyLabel(f, "F2", "HMAC")
	addHotkeyLabel(f, "F3", "Escape")
	addHotkeyLabel(f, "F4", "Calc")

	w.OnKeyDown(func(e ui.Event, i interface{}) bool {
		if e.Key == termbox.KeyCtrlC {
			ui.Stop()
		} else if e.Key == termbox.KeyCtrl2 {
			pane.SetVisible(!pane.Visible())
			pane2.SetVisible(!pane2.Visible())
		}

		return false
	}, nil)

	ui.ActivateControl(w, e1)

	termbox.SetInputMode(termbox.InputEsc)

	//f := ui.CreateFrame(w, ui.AutoSize, ui.AutoSize, ui.BorderThin, ui.Fixed)

	//ui.CreateButton(f, ui.AutoSize, ui.AutoSize, "Quit", ui.Fixed).SetShadowType(ui.ShadowNone)

	//_ = f

	//w.SetPack(ui.Vertical)
	// ui.CreateButton(f, ui.AutoSize, 10, "Quit", 1).OnClick(func(e ui.Event) {
	// 	ui.Stop()
	// })
	// ui.CreateButton(f, 0, 0, "Quit", ui.Fixed).SetShadowType(ui.ShadowNone)
	// ui.CreateButton(f, 0, 0, "Quit", ui.Fixed).SetShadowType(ui.ShadowHalf)
	ui.MainLoop()
	defer ui.DeinitLibrary()
}

func addEditWithLabel(parent ui.Control, label string, enabled bool) *ui.EditField {
	f := ui.CreateFrame(parent, ui.AutoSize, 1, ui.BorderNone, ui.Fixed)
	f.SetPack(ui.Horizontal)
	ui.CreateLabel(f, ui.AutoSize, 1, label+": ", ui.Fixed)
	e := ui.CreateEditField(f, ui.AutoSize, "", 1)
	e.SetEnabled(enabled)
	return e
}

func addHotkeyLabel(parent ui.Control, hotkey, name string) {
	label := fmt.Sprintf("<b:red> %s <b:> %s", hotkey, name)
	width := len(ui.UnColorizeText(label))
	ui.CreateLabel(parent, width, 1, label, ui.Fixed)
}
