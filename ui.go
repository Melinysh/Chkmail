package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type UI struct {
	Sub      Subscriber
	gui      *gocui.Gui
	messages []EmailMessage
	currPos  int
}

func NewUI() UI {
	return UI{NewSubscriber(), gocui.NewGui(), []EmailMessage{}, 0}
}

func (self *UI) Close() {
	self.gui.Close()
}

func (self *UI) ListenForChanges() {
	go func() {
		for {
			event := <-self.Sub.emailEvents
			switch event.Action {
			case Recieved:
				self.messages = append(self.messages, event.Email)
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
	if v, err := g.SetView("side", -1, -1, maxX/4, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		v.Editable = false
		for _, e := range self.messages {
			fmt.Fprintln(v, e.Subject)
		}
		if err := g.SetCurrentView("side"); err != nil {
			return err
		}

	}
	if v, err := g.SetView("main", maxX/4, -1, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		if len(self.messages) == 0 {
			return nil
		}
		if len(self.messages) < self.currPos || self.currPos < 0 {
			panic("currPos is out of order")
		}

		e := self.messages[self.currPos]
		fmt.Fprintf(v, "%s", e.ToString())
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
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, self.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, self.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, self.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, self.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, self.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, self.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowRight, gocui.ModNone, self.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, self.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, self.getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, self.delMsg); err != nil {
		return err
	}

	return nil
}

func (self *UI) refreshUI(g *gocui.Gui, v *gocui.View) error {
	main, err := g.View("main")
	if err != nil {
		return err
	}
	if len(self.messages) == 0 {
		return nil
	}
	if len(self.messages) < self.currPos || self.currPos < 0 {
		panic("currPos is out of order")
	}
	e := self.messages[self.currPos]
	main.Clear()
	fmt.Fprintf(main, "%s", e.ToString())

	side, sideErr := g.View("side")
	if sideErr != nil {
		return sideErr
	}
	side.Clear()
	for _, e := range self.messages {
		fmt.Fprintln(side, e.Subject)
	}

	return nil
}

func (self *UI) nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("side")
}

func (self *UI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if self.currPos == len(self.messages)-1 {
		return nil
	}
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		if v.Name() == "side" {
			self.currPos++
			return self.refreshUI(g, v)
		}
	}
	return nil
}

func (self *UI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if self.currPos == 0 && v.Name() == "side" {
		return nil
	}
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		if v.Name() == "side" {
			self.currPos--
			return self.refreshUI(g, v)
		}
	}
	return nil
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
	if err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}
