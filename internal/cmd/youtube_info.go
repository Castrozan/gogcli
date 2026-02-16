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

type YouTubeInfoCmd struct {
	VideoIDsOrURLs []string `arg:"" name:"videoIdsOrUrls" help:"Video IDs or URLs"`
}

func (c *YouTubeInfoCmd) Run(ctx context.Context, flags *RootFlags) error {
	u := ui.FromContext(ctx)
	account, err := requireAccount(flags)
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

	svc, err := newYouTubeService(ctx, account)
	if err != nil {
		return err
	}

	call := svc.Videos.List([]string{"snippet", "contentDetails", "statistics"}).
		Id(strings.Join(videoIDs, ","))

	response, err := call.Context(ctx).Do()
	if err != nil {
		return err
	}

	var videos []*youtube.Video
	videos = response.Items
	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(ctx, os.Stdout, map[string]any{
			"videos": videos,
		})
	}

	if len(videos) == 0 {
		u.Err().Println("No videos found")
		return nil
	}

	if len(videos) == 1 {
		video := videos[0]
		u.Out().Printf("id\t%s", video.Id)
		u.Out().Printf("title\t%s", video.Snippet.Title)
		u.Out().Printf("channel\t%s", video.Snippet.ChannelTitle)
		u.Out().Printf("duration\t%s", video.ContentDetails.Duration)
		u.Out().Printf("views\t%d", video.Statistics.ViewCount)
		u.Out().Printf("likes\t%d", video.Statistics.LikeCount)
		u.Out().Printf("comments\t%d", video.Statistics.CommentCount)
		if strings.TrimSpace(video.Snippet.Description) != "" {
			u.Out().Printf("description\t%s", video.Snippet.Description)
		}
		return nil
	}

	w, flush := tableWriter(ctx)
	defer flush()
	fmt.Fprintln(w, "VIDEO_ID\tTITLE\tCHANNEL\tDURATION\tVIEWS")
	for _, video := range videos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			video.Id,
			video.Snippet.Title,
			video.Snippet.ChannelTitle,
			video.ContentDetails.Duration,
			video.Statistics.ViewCount)
	}
	return nil
}
