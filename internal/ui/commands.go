package ui

import (
	"context"
	"os/exec"
	"runtime"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/genkio/slack-tui/internal/config"
	"github.com/genkio/slack-tui/internal/mcp"
	"github.com/genkio/slack-tui/internal/slack"
)

// Messages flowing back into the update loop from background tool calls.
type (
	unreadsMsg struct{ convs []slack.Conversation }
	historyMsg struct {
		convID string
		msgs   []slack.Message
	}
	repliesMsg struct {
		threadTS string
		msgs     []slack.Message
	}
	markedMsg      struct{ label string }
	openedMsg      struct{}
	autoRefreshMsg struct{}
	errMsg         struct{ err error }
)

func fetchUnreads(ctx context.Context, c *mcp.Client, u config.UnreadsConfig) tea.Cmd {
	return func() tea.Msg {
		convs, err := c.Unreads(ctx, u)
		if err != nil {
			return errMsg{err}
		}
		return unreadsMsg{convs}
	}
}

// scheduleRefresh emits an autoRefreshMsg after d; the update loop reschedules
// it to form a recurring timer.
func scheduleRefresh(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return autoRefreshMsg{} })
}

func fetchHistory(ctx context.Context, c *mcp.Client, convID, limit string) tea.Cmd {
	return func() tea.Msg {
		msgs, err := c.History(ctx, convID, limit)
		if err != nil {
			return errMsg{err}
		}
		return historyMsg{convID: convID, msgs: msgs}
	}
}

func fetchReplies(ctx context.Context, c *mcp.Client, convID, threadTS string) tea.Cmd {
	return func() tea.Msg {
		msgs, err := c.Replies(ctx, convID, threadTS)
		if err != nil {
			return errMsg{err}
		}
		return repliesMsg{threadTS: threadTS, msgs: msgs}
	}
}

func markRead(ctx context.Context, c *mcp.Client, convID, ts, label string) tea.Cmd {
	return func() tea.Msg {
		if err := c.MarkRead(ctx, convID, ts); err != nil {
			return errMsg{err}
		}
		return markedMsg{label: label}
	}
}

func openURL(url string) tea.Cmd {
	return func() tea.Msg {
		if err := openInBrowser(url); err != nil {
			return errMsg{err}
		}
		return openedMsg{}
	}
}

func openInBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Run()
	case "windows":
		// The empty "" is start's window-title argument; without it a quoted URL
		// is mistaken for the title and nothing opens.
		return exec.Command("cmd", "/c", "start", "", url).Run()
	default:
		return exec.Command("xdg-open", url).Run()
	}
}
