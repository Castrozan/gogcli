# Add YouTube Data API v3 Support

## Overview

This PR adds comprehensive YouTube Data API v3 support to gogcli, following the exact patterns established by the Tasks API implementation.

## Changes

### Core Infrastructure

- **Authentication** (`internal/googleauth/service.go`)
  - Added `ServiceYouTube` constant
  - Configured scope: `https://www.googleapis.com/auth/youtube`
  - Readonly scope: `https://www.googleapis.com/auth/youtube.readonly`
  - Registered in service order and service info map

- **API Service Factory** (`internal/googleapi/youtube.go`)
  - Created `NewYouTube` function following the existing pattern
  - Returns configured `*youtube.Service` with proper authentication

- **Command Registration** (`internal/cmd/root.go`)
  - Added `YouTubeCmd` to CLI with alias `yt`
  - Updated account help text to include youtube

### Implemented Commands

All commands support both `--json` and `--plain` output modes following existing patterns.

#### 1. Search Videos
```bash
gog youtube search "query" --max 10
```
- Searches YouTube videos
- Supports `--max` for result limits (default: 10)
- Supports `--fail-empty` for exit code 3 when no results
- Accepts both YouTube URLs and raw IDs for all video/playlist arguments

#### 2. List Playlists
```bash
gog youtube playlists
```
- Lists user's playlists
- Supports pagination with `--page`, `--all`, `--max`
- Shows playlist ID, title, and item count

#### 3. List Playlist Videos
```bash
gog youtube playlist <id-or-url>
```
- Lists videos in a specified playlist
- Accepts both playlist IDs and YouTube URLs
- Supports pagination
- Shows item ID, video ID, title, and channel

#### 4. Add Videos to Playlist
```bash
gog youtube playlist-add <playlist-id-or-url> <video-id-or-url>...
```
- Adds one or more videos to a playlist
- Supports `--dry-run` flag
- Accepts both playlist and video IDs or URLs
- Returns added item details

#### 5. Remove Items from Playlist
```bash
gog youtube playlist-remove <item-id>...
```
- Removes items from playlist by item ID
- Includes confirmation prompt (unless `--force`)
- Supports multiple item IDs in one command

#### 6. Create Playlist
```bash
gog youtube playlist-create "title" --description "desc" --privacy private
```
- Creates a new playlist
- Privacy options: `private`, `public`, `unlisted` (default: private)
- Optional description
- Supports `--dry-run` flag

#### 7. Get Video Details
```bash
gog youtube info <video-id-or-url>...
```
- Fetches detailed information for one or more videos
- Shows title, channel, duration, views, likes, comments, description
- Accepts both video IDs and YouTube URLs
- Multiple videos shown in table format

### Helper Functions

- **URL Parsing** (`internal/cmd/youtube_helpers.go`)
  - `extractVideoID`: Extracts video ID from various YouTube URL formats
    - Supports: `youtube.com/watch?v=ID`, `youtu.be/ID`, `/embed/ID`, `/v/ID`
  - `extractPlaylistID`: Extracts playlist ID from YouTube URLs
    - Supports: `youtube.com/playlist?list=ID` and query parameters

### Code Quality

- ✅ Follows exact patterns from Tasks API implementation
- ✅ Consistent error handling with existing codebase
- ✅ All commands support `--json` and `--plain` output modes
- ✅ Proper use of `--dry-run` for destructive operations
- ✅ Confirmation prompts for destructive commands
- ✅ Pagination support where applicable
- ✅ Uses Kong CLI framework consistently
- ✅ Builds successfully: `go build ./...`
- ✅ Passes static analysis: `go vet ./...`
- ✅ No comments (self-documenting code with descriptive names)

## Files Added

- `internal/googleapi/youtube.go` - YouTube service factory
- `internal/cmd/youtube.go` - Main command structure
- `internal/cmd/youtube_helpers.go` - URL parsing utilities
- `internal/cmd/youtube_search.go` - Search command
- `internal/cmd/youtube_playlists.go` - List playlists command
- `internal/cmd/youtube_playlist.go` - List playlist videos command
- `internal/cmd/youtube_playlist_add.go` - Add to playlist command
- `internal/cmd/youtube_playlist_remove.go` - Remove from playlist command
- `internal/cmd/youtube_playlist_create.go` - Create playlist command
- `internal/cmd/youtube_info.go` - Video info command

## Files Modified

- `internal/googleauth/service.go` - Added YouTube service
- `internal/cmd/root.go` - Registered YouTube command

## Additional Fix

Fixed pre-existing compilation errors in `internal/cmd/contacts_crud.go`:
- Corrected function parameter types in `contactsApplyPersonName`
- Corrected function parameter types in `contactsApplyPersonOrganization`

## Testing

The implementation has been verified to:
- Build successfully without errors
- Pass `go vet` checks
- Generate proper help text for all commands
- Follow exact patterns from existing Tasks API implementation

## Dependencies

No new dependencies required. The YouTube Data API v3 is already included in `google.golang.org/api v0.260.0`.

## Usage Example

```bash
# Authenticate with YouTube scope
gog login --services youtube

# Search for videos
gog youtube search "golang tutorial" --max 5

# List your playlists
gog youtube playlists

# Get video details
gog youtube info "dQw4w9WgXcQ"

# Create a playlist
gog youtube playlist-create "My Favorites" --privacy private

# Add videos to playlist
gog youtube playlist-add PLxxxx videoID1 videoID2
```

## Notes

- All commands accept both raw IDs and full YouTube URLs for convenience
- JSON output mode enabled via `--json` flag for scripting
- Follows existing patterns for pagination, error handling, and dry-run mode
- Help text auto-generated by Kong CLI framework
