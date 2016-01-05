package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type MainWindow struct {
	*gocui.View
}

func NewMainWindow(v *gocui.View) *MainWindow {
	return &MainWindow{v}
}

func (self *MainWindow) nextLeftView(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView(SubjectsViewKey)
}

func (self *MainWindow) SetMessage(msg string) {
	self.Clear()
	fmt.Fprintln(self, msg)
}

func (self *MainWindow) cursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := self.Cursor()
	if err := self.SetCursor(cx, cy+1); err != nil {
		ox, oy := self.Origin()
		if err := self.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func (self *MainWindow) cursorUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := self.Origin()
	cx, cy := self.Cursor()
	if err := self.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := self.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}
