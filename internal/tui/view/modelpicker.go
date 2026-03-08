package view

import (
	"os"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/HyperMarble/Luna/internal/agent"
	"github.com/HyperMarble/Luna/internal/config"
	"github.com/HyperMarble/Luna/internal/tui/style"
)

const (
	pickerStateProviders   = 0
	pickerStateModels      = 1
	pickerStateAPIKey      = 2
	pickerStateCustomModel = 3
)

// RenderModelPicker is the exported version used by ui.go's View().
func RenderModelPicker(s State) string { return renderModelPicker(s) }

// renderModelPicker dispatches to the correct sub-view based on picker state.
func renderModelPicker(s State) string {
	switch s.ModelPickerState {
	case pickerStateAPIKey:
		return renderAPIKeyDialog(s)
	case pickerStateCustomModel:
		return renderCustomModelDialog(s)
	default:
		return renderProviderTree(s)
	}
}

// renderProviderTree renders the two-zone provider list with optional
// inline model expansion under the selected provider.
func renderProviderTree(s State) string {
	providers := agent.Providers()
	width := s.Layout.ComposerWidth

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(style.WelcomeSub.Render("  Select a model") + "\n\n")

	for i, prov := range providers {
		selected := i == s.ModelPickerProvIdx
		badge := badgeForProvider(prov)
		b.WriteString(renderProviderRow(prov.DisplayName, badge, selected, width) + "\n")

		// Render inline models when this provider is expanded.
		if i == s.ExpandedProv && s.ModelPickerState == pickerStateModels {
			for j, m := range prov.Models {
				modelSelected := j == s.ModelPickerModIdx
				isLast := j == len(prov.Models)-1
				isActive := m.ModelID == s.ActiveModel

				prefix := "  ├── "
				if isLast {
					prefix = "  └── "
				}
				label := prefix + m.ModelLabel
				if isActive {
					label += " ✓"
				}

				if modelSelected {
					b.WriteString(style.PickerSelected.Render("  "+label) + "\n")
				} else {
					b.WriteString(style.PickerCmd.Render("  "+label) + "\n")
				}
			}
		}
	}

	b.WriteString("\n")
	hint := "  ↑↓ navigate · enter expand · esc cancel"
	if s.ModelPickerState == pickerStateModels {
		hint = "  ↑↓ navigate · enter select · esc back"
	}
	b.WriteString(style.WelcomePath.Render(hint) + "\n")
	return b.String()
}

// renderProviderRow renders a single provider line with badge aligned right.
func renderProviderRow(displayName, badge string, selected bool, width int) string {
	bullet := "    "
	if selected {
		bullet = "  › "
	}
	var label string
	if selected {
		label = style.PickerSelected.Render(bullet + displayName)
	} else {
		label = style.PickerCmd.Render(bullet + displayName)
	}
	gap := strings.Repeat(" ", max(0, width/2-lipgloss.Width(label)))
	return label + gap + badge
}

// badgeForProvider returns the styled badge string for a provider row.
func badgeForProvider(prov agent.ProviderInfo) string {
	if prov.Free {
		return style.BadgeFree.Render("[free]")
	}
	if prov.Name == agent.ProviderOllama {
		return style.BadgeLocked.Render("[local]")
	}
	hasKey := os.Getenv(prov.EnvKey) != "" || config.KeyForProvider(prov.EnvKey) != ""
	if hasKey {
		return style.BadgeUnlocked.Render("[unlocked]")
	}
	return style.BadgeLocked.Render("[API key required]")
}

// renderAPIKeyDialog renders the API key input screen for a paid provider.
func renderAPIKeyDialog(s State) string {
	prov := s.APIKeyProvider

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(style.WelcomeSub.Render("  "+prov.DisplayName+" API Key") + "\n\n")
	b.WriteString(style.PickerDesc.Render("  Paste your API key below:") + "\n")
	b.WriteString("  " + s.APIKeyInput.View() + "\n\n")
	b.WriteString(style.WelcomePath.Render("  enter confirm · esc cancel") + "\n")
	if prov.KeyURL != "" {
		b.WriteString(style.WelcomePath.Render("  Get key: "+prov.KeyURL) + "\n")
	}
	return b.String()
}

// renderCustomModelDialog renders the custom model ID input for OpenRouter.
func renderCustomModelDialog(s State) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(style.WelcomeSub.Render("  OpenRouter — Custom Model") + "\n\n")
	b.WriteString(style.PickerDesc.Render("  Paste the model ID (provider/model-name):") + "\n")
	b.WriteString("  " + s.CustomModelInput.View() + "\n\n")
	b.WriteString(style.WelcomePath.Render("  enter confirm · esc back") + "\n")
	b.WriteString(style.WelcomePath.Render("  Browse models: openrouter.ai/models") + "\n")
	return b.String()
}
