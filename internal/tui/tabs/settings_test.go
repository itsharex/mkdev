package tabs_test

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
	"github.com/venkatkrishna07/mkdev/internal/config"
	"github.com/venkatkrishna07/mkdev/internal/tui/styles"
	"github.com/venkatkrishna07/mkdev/internal/tui/tabs"
)

func TestSettingsRendersFieldsFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Default()
	cfg.ProxyPort = 8443
	require.NoError(t, config.Save(filepath.Join(dir, "config.toml"), cfg))
	s := tabs.NewSettings(styles.NewTheme(), dir)
	out := s.View()
	require.Contains(t, out, "tld")
	require.Contains(t, out, "proxy_port")
	require.Contains(t, out, "8443")
}

func TestSettingsSavePersistsConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Default()
	require.NoError(t, config.Save(filepath.Join(dir, "config.toml"), cfg))
	s := tabs.NewSettings(styles.NewTheme(), dir)
	_, _ = s.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	saved, err := config.Load(filepath.Join(dir, "config.toml"))
	require.NoError(t, err)
	require.Equal(t, cfg.TLD, saved.TLD)
	require.Equal(t, cfg.ProxyPort, saved.ProxyPort)
}
