package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"

	"github.com/steipete/gogcli/internal/outfmt"
	"github.com/steipete/gogcli/internal/ui"
)

type YouTubeSearchCmd struct {
	Query     []string `arg:"" name:"query" help:"Search query"`
	Max       int64    `name:"max" aliases:"limit" help:"Max results" default:"10"`
	FailEmpty bool     `name:"fail-empty" aliases:"non-empty,require-results" help:"Exit with code 3 if no results"`
}

func (c *YouTubeSearchCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
	if err != nil {
		return err
	}

	query := strings.TrimSpace(strings.Join(c.Query, " "))
	if query == "" {
		return usage("empty search query")
	}

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	call := svc.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(c.Max).
		Type("video")

	response, err := call.Context(ctx).Do()
	if err != nil {
		return err
	}

	var items []*youtube.SearchResult
	items = response.Items
	if outfmt.IsJSON(ctx) {
		if err := outfmt.WriteJSON(ctx, os.Stdout, map[string]any{
			"videos": items,
		}); err != nil {
			return err
		}
		if len(items) == 0 {
			return failEmptyExit(c.FailEmpty)
		}
		return nil
	}

	if len(items) == 0 {
		u.Err().Println("No videos found")
		return failEmptyExit(c.FailEmpty)
	}

	w, flush := tableWriter(ctx)
	defer flush()
	fmt.Fprintln(w, "VIDEO_ID\tTITLE\tCHANNEL")
	for _, item := range items {
		videoID := item.Id.VideoId
		title := item.Snippet.Title
		channel := item.Snippet.ChannelTitle
		fmt.Fprintf(w, "%s\t%s\t%s\n", videoID, title, channel)
	}
	return nil
}
