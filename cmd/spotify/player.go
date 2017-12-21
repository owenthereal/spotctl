package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Resumes playback where Spotify last left off.",
	RunE:  play,
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pauses Spotify playback.",
	RunE:  pause,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the current player status.",
	RunE:  status,
}

func play(cmd *cobra.Command, args []string) error {
	var opt *spotify.PlayOptions

	if len(args) > 0 {
		// if args start with a spotify ID, play it directly, otherwise search for tracks
		if strings.Contains(args[0], "spotify:") {
			arg := args[0] // only play for the first arg

			var (
				uris    []spotify.URI
				context *spotify.URI
			)

			if strings.Contains(arg, "spotify:track") {
				uris = append(uris, spotify.URI(arg))
			} else {
				uri := spotify.URI(arg)
				context = &uri
			}

			opt = &spotify.PlayOptions{
				PlaybackContext: context,
				URIs:            uris,
			}
		} else {
			result, err := client.Search(strings.Join(args, " "), spotify.SearchTypeTrack)
			if err != nil {
				return err
			}

			if result.Tracks != nil && len(result.Tracks.Tracks) > 0 {
				opt = &spotify.PlayOptions{
					URIs: []spotify.URI{spotify.URI(result.Tracks.Tracks[0].URI)},
				}
			}
		}
	}

	return client.PlayOpt(opt)
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
