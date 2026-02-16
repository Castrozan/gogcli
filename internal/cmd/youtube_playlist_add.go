package cmd

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/youtube/v3"

	"github.com/steipete/gogcli/internal/outfmt"
	"github.com/steipete/gogcli/internal/ui"
)

type YouTubePlaylistAddCmd struct {
	PlaylistIDOrURL string   `arg:"" name:"playlistIdOrUrl" help:"Playlist ID or URL"`
	VideoIDsOrURLs  []string `arg:"" name:"videoIdsOrUrls" help:"Video IDs or URLs to add"`
}

func (c *YouTubePlaylistAddCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
	if err != nil {
		return err
	}

	playlistID, err := extractPlaylistID(c.PlaylistIDOrURL)
	if err != nil {
		return err
	}

	if len(c.VideoIDsOrURLs) == 0 {
		return usage("no video IDs or URLs provided")
	}

	videoIDs := make([]string, 0, len(c.VideoIDsOrURLs))
	for _, videoIDOrURL := range c.VideoIDsOrURLs {
		videoID, err := extractVideoID(videoIDOrURL)
		if err != nil {
			return err
		}
		videoIDs = append(videoIDs, videoID)
	}

	if dryRunErr := dryRunExit(ctx, flags, "youtube.playlist-add", map[string]any{
		"playlist_id": playlistID,
		"video_ids":   videoIDs,
	}); dryRunErr != nil {
		return dryRunErr
	}

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	addedItems := make([]*youtube.PlaylistItem, 0, len(videoIDs))
	for _, videoID := range videoIDs {
		item := &youtube.PlaylistItem{
			Snippet: &youtube.PlaylistItemSnippet{
				PlaylistId: playlistID,
				ResourceId: &youtube.ResourceId{
					Kind:    "youtube#video",
					VideoId: videoID,
				},
			},
		}

		addedItem, err := svc.PlaylistItems.Insert([]string{"snippet"}, item).Context(ctx).Do()
		if err != nil {
			return err
		}
		addedItems = append(addedItems, addedItem)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(ctx, os.Stdout, map[string]any{
			"items": addedItems,
			"count": len(addedItems),
		})
	}

	if len(addedItems) == 1 {
		item := addedItems[0]
		u.Out().Printf("id\t%s", item.Id)
		u.Out().Printf("video_id\t%s", item.ContentDetails.VideoId)
		u.Out().Printf("title\t%s", item.Snippet.Title)
		return nil
	}

	w, flush := tableWriter(ctx)
	defer flush()
	fmt.Fprintln(w, "ITEM_ID\tVIDEO_ID\tTITLE")
	for _, item := range addedItems {
		fmt.Fprintf(w, "%s\t%s\t%s\n", item.Id, item.ContentDetails.VideoId, item.Snippet.Title)
	}
	return nil
}
