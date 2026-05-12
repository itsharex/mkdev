package tabs

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/venkatkrishna07/mkdev/internal/config"
	"github.com/venkatkrishna07/mkdev/internal/tui/styles"
)

const numSettingsFields = 5

// SettingsSavedMsg is emitted after a successful save so the root model can
// refresh its in-memory config copy.
type SettingsSavedMsg struct{ Cfg config.Config }

// Settings is the Settings tab. It edits ~/.mkdev/config.toml in place.
type Settings struct {
	th     styles.Theme
	home   string
	fields [numSettingsFields]textinput.Model
	labels [numSettingsFields]string
	focus  int
	status string
}

// NewSettings constructs a Settings tab seeded from the on-disk config.
func NewSettings(th styles.Theme, home string) Settings {
	s := Settings{th: th, home: home}
	cfg, _ := config.Load(filepath.Join(home, "config.toml"))
	s.labels = [numSettingsFields]string{"tld", "proxy_port", "theme", "log_retention", "log_max_size"}
	s.fields[0] = mkField(cfg.TLD)
	s.fields[1] = mkField(strconv.Itoa(cfg.ProxyPort))
	s.fields[2] = mkField(cfg.Theme)
	s.fields[3] = mkField(cfg.LogRetention)
	s.fields[4] = mkField(cfg.LogMaxSize)
	s.fields[0].Focus()
	return s
}

func mkField(value string) textinput.Model {
	t := textinput.New()
	t.SetValue(value)
	return t
}

// Title implements tabs.Tab.
func (s Settings) Title() string { return "Settings" }

// Init starts the textinput cursor blink.
func (s Settings) Init() tea.Cmd { return textinput.Blink }

// Update advances focus (Tab/Shift+Tab), commits on Enter, reloads on 'r',
// and otherwise forwards the keystroke to the focused textinput.
func (s Settings) Update(msg tea.Msg) (Settings, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.Type {
		case tea.KeyTab:
			s.fields[s.focus].Blur()
			s.focus = (s.focus + 1) % numSettingsFields
			s.fields[s.focus].Focus()
			return s, textinput.Blink
		case tea.KeyShiftTab:
			s.fields[s.focus].Blur()
			s.focus = (s.focus - 1 + numSettingsFields) % numSettingsFields
			s.fields[s.focus].Focus()
			return s, textinput.Blink
		case tea.KeyEnter:
			return s.save()
		}
		if k.String() == "r" {
			return NewSettings(s.th, s.home), textinput.Blink
		}
	}
	var cmd tea.Cmd
	s.fields[s.focus], cmd = s.fields[s.focus].Update(msg)
	return s, cmd
}

func (s Settings) save() (Settings, tea.Cmd) {
	port, err := strconv.Atoi(strings.TrimSpace(s.fields[1].Value()))
	if err != nil || port < 1 || port > 65535 {
		s.status = "✗ proxy_port must be 1–65535"
		return s, nil
	}
	cfg := config.Config{
		TLD:          strings.TrimSpace(s.fields[0].Value()),
		ProxyPort:    port,
		Theme:        strings.TrimSpace(s.fields[2].Value()),
		LogRetention: strings.TrimSpace(s.fields[3].Value()),
		LogMaxSize:   strings.TrimSpace(s.fields[4].Value()),
	}
	if err := config.Save(filepath.Join(s.home, "config.toml"), cfg); err != nil {
		s.status = "✗ save failed: " + err.Error()
		return s, nil
	}
	s.status = "✓ saved · restart for some changes (e.g. proxy_port) to take effect"
	return s, func() tea.Msg { return SettingsSavedMsg{Cfg: cfg} }
}

// View renders the field stack with a focus arrow and a footer hint line.
func (s Settings) View() string {
	var out strings.Builder
	for i, f := range s.fields {
		label := s.th.Dim.Render(s.labels[i] + ":")
		if i == s.focus {
			label = s.th.Title.Render("▶ " + s.labels[i] + ":")
		}
		out.WriteString(label + " " + f.View() + "\n")
	}
	out.WriteString("\n")
	if s.status != "" {
		out.WriteString(s.th.Dim.Render(s.status) + "\n")
	}
	out.WriteString(s.th.Dim.Render("tab next · enter save · r reload"))
	return out.String()
}
