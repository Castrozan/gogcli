package cmd

import (
	"context"
	"fmt"

	"github.com/steipete/gogcli/internal/ui"
)

type YouTubePlaylistRemoveCmd struct {
	ItemIDs []string `arg:"" name:"itemIds" help:"Playlist item IDs to remove"`
}

func (c *YouTubePlaylistRemoveCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
	if err != nil {
		return err
	}

	if len(c.ItemIDs) == 0 {
		return usage("no item IDs provided")
	}

	if confirmErr := confirmDestructive(ctx, flags, fmt.Sprintf("remove %d item(s) from playlist", len(c.ItemIDs))); confirmErr != nil {
		return confirmErr
	}

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	for _, itemID := range c.ItemIDs {
		if err := svc.PlaylistItems.Delete(itemID).Context(ctx).Do(); err != nil {
			return err
		}
	}

	return writeResult(ctx, u,
		kv("deleted", true),
		kv("count", len(c.ItemIDs)),
	)
}
