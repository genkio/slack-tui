package mcp

import (
	"context"

	"github.com/genkio/slack-tui/internal/config"
	"github.com/genkio/slack-tui/internal/slack"
)

// Unreads lists conversations that have unread messages (summary only, no
// message bodies) for a fast list view.
func (c *Client) Unreads(ctx context.Context, u config.UnreadsConfig) ([]slack.Conversation, error) {
	args := map[string]any{
		"include_messages":         false,
		"channel_types":            u.ChannelTypes,
		"max_channels":             u.MaxChannels,
		"max_messages_per_channel": u.MaxMessagesPerChannel,
		"mentions_only":            u.MentionsOnly,
	}
	text, err := c.CallTool(ctx, ToolUnreads, args)
	if err != nil {
		return nil, err
	}
	return DecodeConversations(text)
}

// History fetches recent messages for a conversation. limit follows the
// server's dual-mode string: a duration like "1d"/"7d", or a message count
// like "50". An empty limit uses the server default.
func (c *Client) History(ctx context.Context, channelID, limit string) ([]slack.Message, error) {
	args := map[string]any{"channel_id": channelID}
	if limit != "" {
		args["limit"] = limit
	}
	text, err := c.CallTool(ctx, ToolHistory, args)
	if err != nil {
		return nil, err
	}
	msgs, _, err := DecodeMessages(text)
	return msgs, err
}

// Replies fetches every message in a thread, including the root message.
func (c *Client) Replies(ctx context.Context, channelID, threadTS string) ([]slack.Message, error) {
	args := map[string]any{"channel_id": channelID, "thread_ts": threadTS}
	text, err := c.CallTool(ctx, ToolReplies, args)
	if err != nil {
		return nil, err
	}
	msgs, _, err := DecodeMessages(text)
	return msgs, err
}

// MarkRead marks a conversation read up to ts. An empty ts marks the whole
// conversation read.
func (c *Client) MarkRead(ctx context.Context, channelID, ts string) error {
	args := map[string]any{"channel_id": channelID}
	if ts != "" {
		args["ts"] = ts
	}
	_, err := c.CallTool(ctx, ToolMark, args)
	return err
}
