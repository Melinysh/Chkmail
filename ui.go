package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	MainViewKey     = "main"
	FoldersViewKey  = "folders"
	SubjectsViewKey = "subjects"
)

type UI struct {
	UIPublisher
	Sub              EmailSubscriber
	gui              *gocui.Gui
	mainWindow       *MainWindow
	subjectsWindow   *SubjectsWindow
	foldersWindow    *FoldersWindow
	messages         []EmailMessage
	currPos          int
	currFolder       string
	folderToPos      map[string]int
	folderToMessages map[string][]EmailMessage
}

func NewUI() UI {
	return UI{
		NewUIPublisher(),
		NewEmailSubscriber(),
		gocui.NewGui(),
		nil,
		nil,
		nil,
		[]EmailMessage{},
		0,
		"Inbox",
		map[string]int{},
		map[string][]EmailMessage{},
	}
}

func NewUIWithSubscriber(sub UISubscriber) UI {
	ui := NewUI()
	ui.AddSubscriber(sub)
	return ui
}

func (self *UI) Close() {
	self.gui.Close()
}

func (self *UI) ListenForEmailChanges() {
	go func() {
		for {
			event := <-self.Sub.emailEvents
			switch event.Action {
			case Recieved:
				self.folderToMessages["Inbox"] = append(self.folderToMessages["Inbox"], event.Email)
				self.refreshUI(self.gui, nil)
				continue
			default:
				debugPrint("Unimplemented subscriber")
				// probably refresh gui here for the events
			}
		}
	}()
}

func (self *UI) Init() {
	if err := self.gui.Init(); err != nil {
		panic(err)
	}
	self.gui.SelBgColor = gocui.ColorGreen
	self.gui.SelFgColor = gocui.ColorBlack
	self.gui.ShowCursor = true
	self.gui.SetLayout(self.layout)
	if keyErr := self.keybindings(self.gui); keyErr != nil {
		panic(keyErr)
	}
}

func (self *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func (self *UI) layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(FoldersViewKey, -1, -1, maxX/8, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		v.Editable = false
		fmt.Fprintln(v, "Inbox")
		fmt.Fprintln(v, "Trash")
		fmt.Fprintln(v, "Drafts")
		fmt.Fprintln(v, "Sent")
		self.foldersWindow = NewFoldersWindow(v)
	}

	if v, err := g.SetView(SubjectsViewKey, maxX/8, -1, maxX/4+maxX/8, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		v.Editable = false
		self.subjectsWindow = NewSubjectsWindow(v)
		if err := g.SetCurrentView(SubjectsViewKey); err != nil {
			return err
		}
	}

	if v, err := g.SetView(MainViewKey, maxX/4+maxX/8, -1, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		self.mainWindow = NewMainWindow(v)
	}
	return nil
}

func (self *UI) Loop() {
	err := self.gui.MainLoop() // blocks until finished UI (gets an error)
	if err != nil && err != gocui.Quit {
		panic(err)
	}
}
func (self *UI) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, self.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, self.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(MainViewKey, gocui.KeyArrowLeft, gocui.ModNone, self.mainWindow.nextLeftView); err != nil {
		return err
	}
	if err := g.SetKeybinding(SubjectsViewKey, gocui.KeyArrowLeft, gocui.ModNone, self.subjectsWindow.nextLeftView); err != nil {
		return err
	}
	if err := g.SetKeybinding(SubjectsViewKey, gocui.KeyArrowRight, gocui.ModNone, self.subjectsWindow.nextRightView); err != nil {
		return err
	}
	if err := g.SetKeybinding(FoldersViewKey, gocui.KeyArrowRight, gocui.ModNone, self.foldersWindow.nextRightView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, self.quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, self.delMsg); err != nil {
		return err
	}
	return nil
}

func (self *UI) refreshUI(g *gocui.Gui, v *gocui.View) error {
	self.messages = self.folderToMessages[self.currFolder]
	self.subjectsWindow.Clear()
	self.mainWindow.Clear()
	if len(self.messages) == 0 {
		return self.gui.Flush()
	}
	if len(self.messages) < self.currPos || self.currPos < 0 {
		panic("currPos is out of order" + string(self.currPos) + " msgs: " + string(len(self.messages)))
	}

	self.mainWindow.SetMessage(self.messages[self.currPos].ToString())
	self.subjectsWindow.SetSubjects(self.messages)
	return nil
}

func (self *UI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		defer self.refreshUI(g, v)
		if v.Name() == SubjectsViewKey && self.currPos < len(self.messages)-1 {
			self.currPos++
			return self.subjectsWindow.cursorDown(g, v)
		} else if v.Name() == FoldersViewKey {
			err := self.foldersWindow.cursorDown(g, v)
			self.changeFolders()
			return err //can be nil, from call above
		} else if v.Name() == MainViewKey {
			return self.mainWindow.cursorDown(g, v)
		}
	}
	return nil
}

func (self *UI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		defer self.refreshUI(g, v)
		if v.Name() == SubjectsViewKey && self.currPos > 0 {
			self.currPos--
			return self.subjectsWindow.cursorUp(g, v)
		} else if v.Name() == FoldersViewKey {
			err := self.foldersWindow.cursorUp(g, v)
			self.changeFolders()
			return err
		} else if v.Name() == MainViewKey {
			return self.mainWindow.cursorUp(g, v)
		}

	}
	return nil
}

func (self *UI) debugMsg(msg string) {
	maxX, maxY := self.gui.Size()
	if v, err := self.gui.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		fmt.Fprintln(v, msg)
		self.gui.SetCurrentView("msg")
	}

}

func (self *UI) getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, l)
		if err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func (self *UI) delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if err := g.SetCurrentView(SubjectsViewKey); err != nil {
		return err
	}
	return nil
}

func (self *UI) changeFolders() error {
	v, err := self.gui.View(FoldersViewKey)
	if err != nil {
		debugPrint("Should be folders view", err)
		panic(err)
	}
	_, y := v.Cursor()
	line, lErr := v.Line(y)
	if lErr != nil {
		debugPrint("Unable to get line from folders view", err)
		self.debugMsg("Unable to get line from folders view" + err.Error())
		return nil
		//	panic(err)
	}

	// save current settings to map
	self.folderToPos[self.currFolder] = self.currPos
	self.folderToMessages[self.currFolder] = self.messages

	// fetch last used settings, or to defaults if not found
	if newPos, ok := self.folderToPos[line]; ok {
		self.currPos = newPos
	} else {
		self.currPos = 0
	}

	self.currFolder = line
	return self.refreshUI(self.gui, nil)
}
