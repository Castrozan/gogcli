package cmd

import (
	"fmt"
	"net/url"
	"strings"
)

func extractVideoID(videoIDOrURL string) (string, error) {
	videoIDOrURL = strings.TrimSpace(videoIDOrURL)
	if videoIDOrURL == "" {
		return "", fmt.Errorf("empty video ID or URL")
	}

	if !strings.Contains(videoIDOrURL, "://") && !strings.Contains(videoIDOrURL, "youtube.com") && !strings.Contains(videoIDOrURL, "youtu.be") {
		return videoIDOrURL, nil
	}

	if !strings.Contains(videoIDOrURL, "://") {
		videoIDOrURL = "https://" + videoIDOrURL
	}

	parsedURL, err := url.Parse(videoIDOrURL)
	if err != nil {
		return videoIDOrURL, nil
	}

	if strings.Contains(parsedURL.Host, "youtu.be") {
		videoID := strings.TrimPrefix(parsedURL.Path, "/")
		if videoID != "" {
			return videoID, nil
		}
	}

	if strings.Contains(parsedURL.Host, "youtube.com") {
		if videoID := parsedURL.Query().Get("v"); videoID != "" {
			return videoID, nil
		}

		if strings.HasPrefix(parsedURL.Path, "/embed/") {
			videoID := strings.TrimPrefix(parsedURL.Path, "/embed/")
			if videoID != "" {
				return videoID, nil
			}
		}

		if strings.HasPrefix(parsedURL.Path, "/v/") {
			videoID := strings.TrimPrefix(parsedURL.Path, "/v/")
			if videoID != "" {
				return videoID, nil
			}
		}

		if strings.HasPrefix(parsedURL.Path, "/watch/") {
			videoID := strings.TrimPrefix(parsedURL.Path, "/watch/")
			if videoID != "" {
				return videoID, nil
			}
		}
	}

	return videoIDOrURL, nil
}

func extractPlaylistID(playlistIDOrURL string) (string, error) {
	playlistIDOrURL = strings.TrimSpace(playlistIDOrURL)
	if playlistIDOrURL == "" {
		return "", fmt.Errorf("empty playlist ID or URL")
	}

	if !strings.Contains(playlistIDOrURL, "://") && !strings.Contains(playlistIDOrURL, "youtube.com") {
		return playlistIDOrURL, nil
	}

	if !strings.Contains(playlistIDOrURL, "://") {
		playlistIDOrURL = "https://" + playlistIDOrURL
	}

	parsedURL, err := url.Parse(playlistIDOrURL)
	if err != nil {
		return playlistIDOrURL, nil
	}

	if strings.Contains(parsedURL.Host, "youtube.com") {
		if playlistID := parsedURL.Query().Get("list"); playlistID != "" {
			return playlistID, nil
		}
	}

	return playlistIDOrURL, nil
}
