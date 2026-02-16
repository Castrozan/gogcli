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

type YouTubePlaylistCmd struct {
	PlaylistIDOrURL string `arg:"" name:"playlistIdOrUrl" help:"Playlist ID or URL"`
	Max             int64  `name:"max" aliases:"limit" help:"Max results" default:"50"`
	Page            string `name:"page" aliases:"cursor" help:"Page token"`
	All             bool   `name:"all" aliases:"all-pages,allpages" help:"Fetch all pages"`
	FailEmpty       bool   `name:"fail-empty" aliases:"non-empty,require-results" help:"Exit with code 3 if no results"`
}

func (c *YouTubePlaylistCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
	if err != nil {
		return err
	}

	playlistID, err := extractPlaylistID(c.PlaylistIDOrURL)
	if err != nil {
		return err
	}

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	fetch := func(pageToken string) ([]*youtube.PlaylistItem, string, error) {
		call := svc.PlaylistItems.List([]string{"snippet", "contentDetails"}).
			PlaylistId(playlistID).
			MaxResults(c.Max)
		if strings.TrimSpace(pageToken) != "" {
			call = call.PageToken(strings.TrimSpace(pageToken))
		}

		resp, err := call.Context(ctx).Do()
		if err != nil {
			return nil, "", err
		}
		return resp.Items, resp.NextPageToken, nil
	}

	var items []*youtube.PlaylistItem
	nextPageToken := ""
	if c.All {
		all, err := collectAllPages(c.Page, fetch)
		if err != nil {
			return err
		}
		items = all
	} else {
		var err error
		items, nextPageToken, err = fetch(c.Page)
		if err != nil {
			return err
		}
	}

	if outfmt.IsJSON(ctx) {
		if err := outfmt.WriteJSON(ctx, os.Stdout, map[string]any{
			"items":         items,
			"nextPageToken": nextPageToken,
		}); err != nil {
			return err
		}
		if len(items) == 0 {
			return failEmptyExit(c.FailEmpty)
		}
		return nil
	}

	if len(items) == 0 {
		u.Err().Println("No items in playlist")
		return failEmptyExit(c.FailEmpty)
	}

	w, flush := tableWriter(ctx)
	defer flush()
	fmt.Fprintln(w, "ITEM_ID\tVIDEO_ID\tTITLE\tCHANNEL")
	for _, item := range items {
		videoID := item.ContentDetails.VideoId
		title := item.Snippet.Title
		channel := item.Snippet.ChannelTitle
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.Id, videoID, title, channel)
	}
	printNextPageHint(u, nextPageToken)
	return nil
}
