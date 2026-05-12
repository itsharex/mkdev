package tabs_test

import (
	"path/filepath"
	"testing"

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
