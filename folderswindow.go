package main

import (
	"github.com/jroimartin/gocui"
)

type FoldersWindow struct {
	*gocui.View
}

func NewFoldersWindow(v *gocui.View) *FoldersWindow {
	return &FoldersWindow{v}
}

func (self *FoldersWindow) nextRightView(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView(SubjectsViewKey)
}

func (self *FoldersWindow) cursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := self.Cursor()
	if err := self.SetCursor(cx, cy+1); err != nil {
		ox, oy := self.Origin()
		if err := self.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func (self *FoldersWindow) cursorDown(g *gocui.Gui, v *gocui.View) error {
	ox, oy := self.Origin()
	cx, cy := self.Cursor()
	if err := self.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := self.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}
