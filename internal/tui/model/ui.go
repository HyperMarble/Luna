package model

import (
	"context"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/HyperMarble/Luna/internal/agent"
	"github.com/HyperMarble/Luna/internal/config"
	"github.com/HyperMarble/Luna/internal/tui/events"
	tuilayout "github.com/HyperMarble/Luna/internal/tui/layout"
	"github.com/HyperMarble/Luna/internal/tui/slash"
	"github.com/HyperMarble/Luna/internal/tui/types"
	"github.com/HyperMarble/Luna/internal/tui/view"
)

type pickerState int

const (
	pickerStateProviders pickerState = iota
	pickerStateModels
	pickerStateAPIKey
	pickerStateCustomModel
)

// UI is the main Bubble Tea model and state owner.
type UI struct {
	width  int
	height int
	input  textinput.Model
	layout tuilayout.UI

	spinner   spinner.Model
	messages  []types.Message
	thinking  bool
	verbIdx   int
	pickerIdx int

	// Model picker tree.
	modelPickerOpen    bool
	modelPickerState   pickerState
	modelPickerProvIdx int
	modelPickerModIdx  int
	expandedProv       int // -1 means no provider is expanded.
	apiKeyInput        textinput.Model
	apiKeyProvider     agent.ProviderInfo
	customModelInput   textinput.Model
	activeModel        string

	agent agent.Service
}

// New returns the initial UI model.
func New() UI {
	ti := textinput.New()
	ti.Placeholder = "Ask Luna..."
	ti.Focus()
	ti.CharLimit = 2000
	ti.Prompt = ""
	ti.Width = 76 // Updated on WindowSizeMsg.

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	ki := textinput.New()
	ki.EchoMode = textinput.EchoPassword
	ki.Placeholder = "Paste API key..."
	ki.CharLimit = 200

	ci := textinput.New()
	ci.Placeholder = "e.g. openai/gpt-4o"
	ci.CharLimit = 200

	_ = config.Load()

	return UI{
		input:            ti,
		spinner:          sp,
		layout:           tuilayout.Compute(80),
		agent:            agent.New(nil),
		expandedProv:     -1,
		apiKeyInput:      ki,
		customModelInput: ci,
	}
}

// Init starts the cursor blink when the program launches.
func (m UI) Init() tea.Cmd { return textinput.Blink }

// Input exposes the text input (used in tests).
func (m UI) Input() textinput.Model { return m.input }

// Messages exposes the conversation history (used in tests).
func (m UI) Messages() []types.Message { return m.messages }

// Update routes all Bubble Tea messages and mutates model state.
func (m UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout = tuilayout.Compute(msg.Width)
		m.input.Width = max(10, m.layout.ComposerWidth-4)
		cmds = append(cmds, tea.ClearScreen)

	case tea.KeyMsg:
		cmd, done := m.onKey(msg)
		if done {
			return m, cmd
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case events.UserSubmitMsg:
		m.messages = append(m.messages, types.Message{Role: "user", Content: msg.Text})
		m.thinking = true

	case events.AgentResponseMsg:
		m.thinking = false
		m.messages = append(m.messages, types.Message{Role: "luna", Content: msg.Text})

	case events.SaveAPIKeyMsg:
		if msg.Err != nil {
			m.closePicker()
			m.messages = append(m.messages, types.Message{
				Role:    "luna",
				Content: "Failed to save API key: " + msg.Err.Error(),
			})
			return m, nil
		}
		// Expand the newly-unlocked provider and move to model selection.
		providers := agent.Providers()
		for i, p := range providers {
			if p.EnvKey == msg.EnvKey {
				m.expandedProv = i
				m.modelPickerProvIdx = i
				break
			}
		}
		m.modelPickerState = pickerStateModels
		m.modelPickerModIdx = 0
		m.apiKeyInput.SetValue("")

	case spinner.TickMsg:
		if m.thinking {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			m.verbIdx++
			cmds = append(cmds, spinCmd)
		}
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	// Route non-key messages (e.g. cursor blink) to the active dialog input.
	if m.modelPickerOpen && m.modelPickerState == pickerStateAPIKey {
		var keyCmd tea.Cmd
		m.apiKeyInput, keyCmd = m.apiKeyInput.Update(msg)
		cmds = append(cmds, keyCmd)
	}
	if m.modelPickerOpen && m.modelPickerState == pickerStateCustomModel {
		var keyCmd tea.Cmd
		m.customModelInput, keyCmd = m.customModelInput.Update(msg)
		cmds = append(cmds, keyCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the full UI.
func (m UI) View() string {
	return view.Render(view.State{
		Width:              m.width,
		Height:             m.height,
		Layout:             m.layout,
		Input:              m.input,
		Messages:           m.messages,
		Thinking:           m.thinking,
		VerbIdx:            m.verbIdx,
		PickerIndex:        m.pickerIdx,
		ModelPickerOpen:    m.modelPickerOpen,
		ModelPickerState:   int(m.modelPickerState),
		ModelPickerProvIdx: m.modelPickerProvIdx,
		ModelPickerModIdx:  m.modelPickerModIdx,
		ExpandedProv:       m.expandedProv,
		APIKeyInput:        m.apiKeyInput,
		APIKeyProvider:     m.apiKeyProvider,
		CustomModelInput:   m.customModelInput,
		ActiveModel:        m.activeModel,
	})
}

func (m *UI) onKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if m.modelPickerOpen {
		switch m.modelPickerState {
		case pickerStateProviders:
			return m.onProvidersKey(msg)
		case pickerStateModels:
			return m.onModelsKey(msg)
		case pickerStateAPIKey:
			return m.onAPIKeyInput(msg)
		case pickerStateCustomModel:
			return m.onCustomModelInput(msg)
		}
	}

	switch msg.String() {
	case "ctrl+c":
		return tea.Quit, true
	case "up", "down":
		m.movePicker(msg.String())
	case "tab":
		m.completePicker()
	case "esc":
		m.dismissPicker()
	case "enter":
		return m.onEnter()
	}

	if msg.String() != "up" && msg.String() != "down" {
		m.pickerIdx = 0
	}
	return nil, false
}

func (m *UI) onProvidersKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	providers := agent.Providers()
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit, true
	case "esc":
		m.closePicker()
	case "up":
		if m.modelPickerProvIdx > 0 {
			m.modelPickerProvIdx--
		}
	case "down":
		if m.modelPickerProvIdx < len(providers)-1 {
			m.modelPickerProvIdx++
		}
	case "enter":
		if m.modelPickerProvIdx < len(providers) {
			prov := providers[m.modelPickerProvIdx]
			if prov.Free || m.providerHasKey(prov) {
				m.expandedProv = m.modelPickerProvIdx
				m.modelPickerState = pickerStateModels
				m.modelPickerModIdx = 0
			} else {
				m.apiKeyProvider = prov
				m.apiKeyInput.SetValue("")
				m.apiKeyInput.Focus()
				m.modelPickerState = pickerStateAPIKey
			}
		}
	}
	return nil, true
}

func (m *UI) onModelsKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	providers := agent.Providers()
	var models []agent.ModelEntry
	if m.expandedProv >= 0 && m.expandedProv < len(providers) {
		models = providers[m.expandedProv].Models
	}
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit, true
	case "esc":
		m.modelPickerState = pickerStateProviders
		m.expandedProv = -1
	case "up":
		if m.modelPickerModIdx > 0 {
			m.modelPickerModIdx--
		}
	case "down":
		if m.modelPickerModIdx < len(models)-1 {
			m.modelPickerModIdx++
		}
	case "enter":
		if m.modelPickerModIdx < len(models) {
			e := models[m.modelPickerModIdx]
			if e.ModelID == "__custom__" {
				m.customModelInput.SetValue("")
				m.customModelInput.Focus()
				m.modelPickerState = pickerStateCustomModel
			} else {
				m.selectModel(e)
			}
		}
	}
	return nil, true
}

func (m *UI) onCustomModelInput(msg tea.KeyMsg) (tea.Cmd, bool) {
	var keyCmd tea.Cmd
	m.customModelInput, keyCmd = m.customModelInput.Update(msg)

	switch msg.String() {
	case "ctrl+c":
		return tea.Quit, true
	case "esc":
		m.modelPickerState = pickerStateModels
		m.customModelInput.SetValue("")
	case "enter":
		val := strings.TrimSpace(m.customModelInput.Value())
		if val != "" {
			providers := agent.Providers()
			if m.expandedProv >= 0 && m.expandedProv < len(providers) {
				prov := providers[m.expandedProv]
				e := agent.ModelEntry{
					Provider:    prov.Name,
					DisplayName: prov.DisplayName,
					ModelID:     val,
					ModelLabel:  val,
					Free:        prov.Free,
					EnvKey:      prov.EnvKey,
				}
				m.selectModel(e)
			}
		}
	}
	return keyCmd, true
}

func (m *UI) onAPIKeyInput(msg tea.KeyMsg) (tea.Cmd, bool) {
	// Update the textinput here because key messages return done=true and
	// bypass the generic update at the bottom of Update().
	var keyCmd tea.Cmd
	m.apiKeyInput, keyCmd = m.apiKeyInput.Update(msg)

	switch msg.String() {
	case "ctrl+c":
		return tea.Quit, true
	case "esc":
		m.modelPickerState = pickerStateProviders
		m.apiKeyInput.SetValue("")
	case "enter":
		val := strings.TrimSpace(m.apiKeyInput.Value())
		if val != "" {
			return tea.Batch(saveAPIKeyCmd(m.apiKeyProvider.EnvKey, val), keyCmd), true
		}
	}
	return keyCmd, true
}

func (m *UI) selectModel(e agent.ModelEntry) {
	m.activeModel = e.ModelID
	m.closePicker()
	m.agent = agent.NewWithModel(string(e.Provider), e.ModelID)
	m.messages = append(m.messages, types.Message{
		Role:    "luna",
		Content: "Switched to **" + e.ModelLabel + "** (" + e.DisplayName + ")",
	})
}

func (m *UI) closePicker() {
	m.modelPickerOpen = false
	m.modelPickerState = pickerStateProviders
	m.modelPickerProvIdx = 0
	m.modelPickerModIdx = 0
	m.expandedProv = -1
	m.apiKeyInput.SetValue("")
}

func (m *UI) providerHasKey(prov agent.ProviderInfo) bool {
	return os.Getenv(prov.EnvKey) != "" || config.KeyForProvider(prov.EnvKey) != ""
}

func (m *UI) movePicker(dir string) {
	if !strings.HasPrefix(m.input.Value(), "/") {
		return
	}
	filtered := slash.Filtered(m.input.Value())
	if len(filtered) == 0 {
		return
	}
	if dir == "up" {
		m.pickerIdx = max(0, m.pickerIdx-1)
	} else {
		m.pickerIdx = min(len(filtered)-1, m.pickerIdx+1)
	}
}

func (m *UI) completePicker() {
	if !strings.HasPrefix(m.input.Value(), "/") {
		return
	}
	filtered := slash.Filtered(m.input.Value())
	if len(filtered) > 0 {
		m.input.SetValue(filtered[m.pickerIdx].Name)
		m.input.CursorEnd()
	}
}

func (m *UI) dismissPicker() {
	if strings.HasPrefix(m.input.Value(), "/") {
		m.input.SetValue("")
		m.pickerIdx = 0
	}
}

func (m *UI) onEnter() (tea.Cmd, bool) {
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		return nil, false
	}
	if strings.HasPrefix(text, "/") {
		return m.executeSlash(text)
	}
	return m.submitText(text), false
}

func (m *UI) submitText(text string) tea.Cmd {
	m.input.SetValue("")
	m.messages = append(m.messages, types.Message{Role: "user", Content: text})
	m.thinking = true
	return tea.Batch(agentResponseCmd(m.agent, text), m.spinner.Tick)
}

func (m *UI) executeSlash(text string) (tea.Cmd, bool) {
	filtered := slash.Filtered(text)
	if len(filtered) > 0 && text != filtered[m.pickerIdx].Name {
		text = filtered[m.pickerIdx].Name
	}
	m.input.SetValue("")
	m.pickerIdx = 0
	return m.handleSlash(text), true
}

func (m *UI) handleSlash(cmd string) tea.Cmd {
	switch cmd {
	case "/exit":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/exit"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "Goodbye! Thanks for using Luna."})
		return tea.Quit
	case "/clear":
		m.thinking = false
		m.messages = []types.Message{{Role: "luna", Content: "Conversation cleared."}}
		return nil
	case "/help":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/help"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: helpText()})
		return nil
	case "/model":
		m.input.SetValue("")
		m.modelPickerOpen = true
		m.modelPickerState = pickerStateProviders
		m.modelPickerProvIdx = 0
		return nil
	case "/plugins":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/plugins"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "No plugins installed yet.\n\nPlugin system coming soon."})
		return nil
	default:
		m.messages = append(m.messages, types.Message{Role: "user", Content: cmd})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "Unknown command: `" + cmd + "`\n\nType `/help` to see available commands."})
		return nil
	}
}

func helpText() string {
	return `**Available commands**

| Command | Description |
|---------|-------------|
| ` + "`/help`" + ` | Show this help |
| ` + "`/clear`" + ` | Clear the conversation |
| ` + "`/model`" + ` | Change the active model |
| ` + "`/plugins`" + ` | Manage plugins |
| ` + "`/exit`" + ` | Exit Luna |`
}

func agentResponseCmd(svc agent.Service, text string) tea.Cmd {
	return func() tea.Msg {
		resp, err := svc.Run(context.Background(), agent.Request{Prompt: text})
		if err != nil {
			return events.AgentResponseMsg{Text: "Agent error: " + err.Error()}
		}
		return events.AgentResponseMsg{Text: resp.Text}
	}
}

func saveAPIKeyCmd(envKey, value string) tea.Cmd {
	return func() tea.Msg {
		err := config.SetKey(envKey, value)
		return events.SaveAPIKeyMsg{EnvKey: envKey, Value: value, Err: err}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
