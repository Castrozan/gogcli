package cmd

import (
	"github.com/steipete/gogcli/internal/googleapi"
)

var newYouTubeService = googleapi.NewYouTube

type YouTubeCmd struct {
	Search         YouTubeSearchCmd         `cmd:"" name:"search" help:"Search videos"`
	Playlists      YouTubePlaylistsCmd      `cmd:"" name:"playlists" help:"List user's playlists"`
	Playlist       YouTubePlaylistCmd       `cmd:"" name:"playlist" help:"List videos in a playlist"`
	PlaylistAdd    YouTubePlaylistAddCmd    `cmd:"" name:"playlist-add" help:"Add videos to playlist"`
	PlaylistRemove YouTubePlaylistRemoveCmd `cmd:"" name:"playlist-remove" help:"Remove items from playlist"`
	PlaylistCreate YouTubePlaylistCreateCmd `cmd:"" name:"playlist-create" help:"Create a playlist"`
	Info           YouTubeInfoCmd           `cmd:"" name:"info" help:"Get video details"`
}
