package model

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

var lastMouseEvent time.Time

// MouseEventFilter rate-limits mouse wheel / motion events so trackpad
// flooding doesn't overwhelm the update loop (mirrors crush's pattern).
func MouseEventFilter(_ tea.Model, msg tea.Msg) tea.Msg {
	switch msg.(type) {
	case tea.MouseWheelMsg, tea.MouseMotionMsg:
		now := time.Now()
		if now.Sub(lastMouseEvent) < 15*time.Millisecond {
			return nil
		}
		lastMouseEvent = now
	}
	return msg
}
