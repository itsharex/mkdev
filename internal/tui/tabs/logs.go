package tabs

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/venkatkrishna07/mkdev/internal/tui/styles"
)

// LogsTickMsg is sent on each refresh tick.
type LogsTickMsg time.Time

// Logs is the Logs tab. It tails a daemon log file.
type Logs struct {
	th       styles.Theme
	logPath  string
	viewport viewport.Model
	width    int
	height   int
	paused   bool
}

// NewLogs constructs a Logs tab tailing path.
func NewLogs(th styles.Theme, logPath string) Logs {
	vp := viewport.New(100, 10)
	vp.SetContent("(no log entries yet)")
	return Logs{th: th, logPath: logPath, viewport: vp}
}

// Title implements tabs.Tab.
func (l Logs) Title() string { return "Logs" }

// Init starts the tail tick.
func (l Logs) Init() tea.Cmd { return logsTickCmd() }

// Update handles ticks, viewport scrolling, and tab-local keys.
func (l Logs) Update(msg tea.Msg) (Logs, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = m.Width
		l.height = m.Height
		l.viewport.Width = m.Width - 2
		l.viewport.Height = max(m.Height-10, 5)
	case LogsTickMsg:
		if !l.paused {
			l.refresh()
		}
		return l, logsTickCmd()
	case tea.KeyMsg:
		switch m.String() {
		case " ":
			l.paused = !l.paused
			return l, nil
		case "c":
			l.viewport.SetContent("")
			return l, nil
		}
	}
	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

// refresh reads the tail of the log file into the viewport. The last
// 2*viewport.Height lines are retained so a scroll-up gesture still has
// recent context without keeping the entire file in memory.
func (l *Logs) refresh() {
	f, err := os.Open(l.logPath)
	if err != nil {
		l.viewport.SetContent(l.th.Dim.Render("log file not yet present at " + l.logPath))
		return
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	keep := l.viewport.Height
	if keep < 1 {
		keep = 20
	}
	if len(lines) > keep*2 {
		lines = lines[len(lines)-keep*2:]
	}
	l.viewport.SetContent(strings.Join(lines, "\n"))
	if !l.paused {
		l.viewport.GotoBottom()
	}
}

// View renders the header (path + paused pill) above the viewport body.
func (l Logs) View() string {
	hdr := l.th.Dim.Render("tailing " + l.logPath)
	if l.paused {
		hdr += "  " + l.th.PillDown.Render("PAUSED")
	}
	return hdr + "\n" + l.viewport.View()
}

func logsTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return LogsTickMsg(t) })
}
