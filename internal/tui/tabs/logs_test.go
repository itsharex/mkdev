package tabs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/venkatkrishna07/mkdev/internal/tui/styles"
	"github.com/venkatkrishna07/mkdev/internal/tui/tabs"
)

func TestLogsViewShowsTailingHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tui.log")
	require.NoError(t, os.WriteFile(path, []byte("line one\nline two\n"), 0o600))
	l := tabs.NewLogs(styles.NewTheme(), path)
	require.Contains(t, l.View(), path)
}
