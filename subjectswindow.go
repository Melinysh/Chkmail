package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type SubjectsWindow struct {
	*gocui.View
	subjects []string
}

func NewSubjectsWindow(v *gocui.View) *SubjectsWindow {
	return &SubjectsWindow{v, []string{}}
}

func (self *SubjectsWindow) nextRightView(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView(MainViewKey)
}

func (self *SubjectsWindow) nextLeftView(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView(FoldersViewKey)
}

func (self *SubjectsWindow) cursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := self.Cursor()
	if err := self.SetCursor(cx, cy+1); err != nil {
		ox, oy := self.Origin()
		if err := self.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func (self *SubjectsWindow) cursorUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := self.Origin()
	cx, cy := self.Cursor()
	if err := self.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := self.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func (self *SubjectsWindow) SetSubjects(msgs []EmailMessage) {
	var subjects []string
	for _, m := range msgs {
		subjects = append(subjects, m.Subject)
	}
	self.subjects = subjects
	self.Clear()
	for _, s := range subjects {
		fmt.Fprintln(self, s)
	}
}
