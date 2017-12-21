package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var playCmdFlagType string

var playCmd = &cobra.Command{
	Use:   "play [name]",
	Short: "Resume playback or play a song, album, artist or playlist by name.",
	Long:  `Resume playback or find a song, album, artist or playlist by name and play it. The search type is specified with --type.`,
	RunE:  play,
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause Spotify playback.",
	RunE:  pause,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current player status.",
	RunE:  status,
}

func play(cmd *cobra.Command, args []string) error {
	var (
		opt *spotify.PlayOptions
		err error
	)

	if len(args) > 0 {
		// if args start with a spotify ID, play it directly, otherwise search for songs
		if strings.Contains(args[0], "spotify:") {
			opt = playByID(args[0]) // only play the first id
		} else {
			opt, err = searchToPlay(strings.Join(args, " "), playCmdFlagType)
			if err != nil {
				return err
			}
		}
	}

	return client.PlayOpt(opt)
}

func playByID(id string) *spotify.PlayOptions {
	var (
		uris    []spotify.URI
		context *spotify.URI
	)

	if strings.Contains(id, "spotify:track") {
		uris = append(uris, spotify.URI(id))
	} else {
		uri := spotify.URI(id)
		context = &uri
	}

	return &spotify.PlayOptions{
		PlaybackContext: context,
		URIs:            uris,
	}
}

func searchToPlay(query, t string) (*spotify.PlayOptions, error) {
	var st spotify.SearchType
	switch t {
	case "track":
		st = spotify.SearchTypeTrack
	case "album":
		st = spotify.SearchTypeAlbum
	case "artist":
		st = spotify.SearchTypeArtist
	case "playlist":
		st = spotify.SearchTypePlaylist
	default:
		return nil, fmt.Errorf("unsupported search type %s", t)
	}

	result, err := client.Search(query, st)
	if err != nil {
		return nil, err
	}

	var opt *spotify.PlayOptions
	switch t {
	case "track":
		if result.Tracks != nil && len(result.Tracks.Tracks) > 0 {
			opt = &spotify.PlayOptions{
				URIs: []spotify.URI{result.Tracks.Tracks[0].URI},
			}
		}
	case "album":
		if result.Albums != nil && len(result.Albums.Albums) > 0 {
			opt = &spotify.PlayOptions{
				PlaybackContext: &result.Albums.Albums[0].URI,
			}
		}
	case "artist":
		if result.Artists != nil && len(result.Artists.Artists) > 0 {
			opt = &spotify.PlayOptions{
				PlaybackContext: &result.Artists.Artists[0].URI,
			}
		}
	case "playlist":
		if result.Playlists != nil && len(result.Playlists.Playlists) > 0 {
			opt = &spotify.PlayOptions{
				PlaybackContext: &result.Playlists.Playlists[0].URI,
			}
		}
	}

	return opt, nil
}

func pause(cmd *cobra.Command, args []string) error {
	return client.Pause()
}

func status(cmd *cobra.Command, args []string) error {
	state, err := client.PlayerState()
	if err != nil {
		return err
	}

	if state.Playing && state.Item != nil {
		var artists []string
		for _, a := range state.Item.Artists {
			artists = append(artists, a.Name)
		}

		fmt.Println("Spoitfy is currently playing.")
		fmt.Printf("Artist: %s\n", strings.Join(artists, ", "))
		fmt.Printf("Album: %s\n", state.Item.Album.Name)
		fmt.Printf("Track: %s\n", state.Item.Name)
		fmt.Printf("Position: %s / %s\n", formatDurationInMillisecond(state.Progress), formatDurationInMillisecond(state.Item.Duration))
	} else {
		fmt.Println("Spotify is currently paused.")
	}

	return nil
}

func formatDurationInMillisecond(d int) time.Duration {
	return (time.Duration(d) * time.Millisecond).Round(time.Second)
}
