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

type YouTubePlaylistsCmd struct {
	Max       int64  `name:"max" aliases:"limit" help:"Max results" default:"25"`
	Page      string `name:"page" aliases:"cursor" help:"Page token"`
	All       bool   `name:"all" aliases:"all-pages,allpages" help:"Fetch all pages"`
	FailEmpty bool   `name:"fail-empty" aliases:"non-empty,require-results" help:"Exit with code 3 if no results"`
}

func (c *YouTubePlaylistsCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
	if err != nil {
		return err
	}

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	fetch := func(pageToken string) ([]*youtube.Playlist, string, error) {
		call := svc.Playlists.List([]string{"snippet", "contentDetails"}).
			Mine(true).
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

	var items []*youtube.Playlist
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
			"playlists":     items,
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
		u.Err().Println("No playlists")
		return failEmptyExit(c.FailEmpty)
	}

	w, flush := tableWriter(ctx)
	defer flush()
	fmt.Fprintln(w, "PLAYLIST_ID\tTITLE\tITEMS")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\t%d\n", item.Id, item.Snippet.Title, item.ContentDetails.ItemCount)
	}
	printNextPageHint(u, nextPageToken)
	return nil
}
