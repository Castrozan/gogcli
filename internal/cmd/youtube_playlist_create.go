package cmd

import (
	"context"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"

	"github.com/steipete/gogcli/internal/outfmt"
	"github.com/steipete/gogcli/internal/ui"
)

type YouTubePlaylistCreateCmd struct {
	Title       []string `arg:"" name:"title" help:"Playlist title"`
	Description string   `name:"description" help:"Playlist description"`
	Privacy     string   `name:"privacy" help:"Privacy status: private|public|unlisted" default:"private"`
}

func (c *YouTubePlaylistCreateCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
	if err != nil {
		return err
	}

	title := strings.TrimSpace(strings.Join(c.Title, " "))
	if title == "" {
		return usage("empty title")
	}

	privacy := strings.TrimSpace(strings.ToLower(c.Privacy))
	if privacy != "private" && privacy != "public" && privacy != "unlisted" {
		return usage("invalid privacy status (expected private, public, or unlisted)")
	}

	if dryRunErr := dryRunExit(ctx, flags, "youtube.playlist-create", map[string]any{
		"title":       title,
		"description": c.Description,
		"privacy":     privacy,
	}); dryRunErr != nil {
		return dryRunErr
	}

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	playlist := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:       title,
			Description: strings.TrimSpace(c.Description),
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: privacy,
		},
	}

	created, err := svc.Playlists.Insert([]string{"snippet", "status"}, playlist).Context(ctx).Do()
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(ctx, os.Stdout, map[string]any{"playlist": created})
	}

	u.Out().Printf("id\t%s", created.Id)
	u.Out().Printf("title\t%s", created.Snippet.Title)
	u.Out().Printf("privacy\t%s", created.Status.PrivacyStatus)
	if strings.TrimSpace(created.Snippet.Description) != "" {
		u.Out().Printf("description\t%s", created.Snippet.Description)
	}
	return nil
}
