package tui

import (
	"github.com/HyperMarble/Luna/internal/tui/events"
	"github.com/HyperMarble/Luna/internal/tui/model"
	"github.com/HyperMarble/Luna/internal/tui/types"
)

type Message = types.Message
type UserSubmitMsg = events.UserSubmitMsg
type AgentResponseMsg = events.AgentResponseMsg
type LunaStubMsg = events.LunaStubMsg

type Model = model.UI

func NewModel() Model {
	return model.New()
}
