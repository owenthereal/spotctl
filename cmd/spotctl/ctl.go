package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	playCmdFlagType string
	deviceNameFlag  string
)

var playCmd = &cobra.Command{
	Use:   "play [name]",
	Short: "Resume playback or play a track, album, artist or playlist by name",
	Long:  `Resume playback or find a track, album, artist or playlist by name and play it. The search type can be specified with --type.`,
	RunE:  play,
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause Spotify playback",
	RunE:  pause,
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Skip to the next track",
	RunE:  next,
}

var prevCmd = &cobra.Command{
	Use:   "prev",
	Short: "Return to the previous track",
	RunE:  prev,
}

var volCmd = &cobra.Command{
	Use:   "vol [up|down|amount]",
	Short: "Set or return volume percentage",
	Long:  `Set volume percentage to an amount between 0 and 100. If arg is up, volume is increased by 10%. If arg is down, volume is decreased by 10%. If no arg is provided, current volume percentage is returned.`,
	RunE:  vol,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current player status",
	RunE:  status,
}

var shuffleCmd = &cobra.Command{
	Use:   "shuffle",
	Short: "Toggle shuffle playback mode",
	RunE:  shuffle,
}

var repeatCmd = &cobra.Command{
	Use:   "repeat",
	Short: "Toggle repeat playback mode",
	RunE:  repeat,
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Show list of available devices",
	RunE:  devices,
}

var setDeviceCmd = &cobra.Command{
	Use:   "setdevice [name]",
	Short: "Set playback device",
	RunE:  setDevice,
}

var getAlbumsCmd = &cobra.Command{
    Use: "albums",
    Short: "Show list of saved albums.",
    RunE: getAlbums,
}

var getSongsCmd = &cobra.Command{
    Use: "tracks",
    Short: "Show list of saved tracks.",
    RunE: getSongs,
}

var getArtistsCmd = &cobra.Command{
    Use: "artists",
    Short: "Show list of saved artists.",
    RunE: getArtists,
}

func devices(cmd *cobra.Command, args []string) error {
	devices, err := client.PlayerDevices()
	if err != nil {
		return err
	}

	for _, device := range devices {
		active := ""
		if device.Active {
			active = "* "
		}
		fmt.Printf("%s%s - %s (volume %d%%)\n", active, device.Name, device.Type, device.Volume)
	}

	return nil
}

func setDevice(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
        return nil
    }
	devices, err := client.PlayerDevices()
	if err != nil {
		return err
	}
    var wanted string
    wanted = strings.Join(args, " ")
	for _, device := range devices {
		if wanted == device.Name {
            fmt.Printf ("Transfering playback to \"%s\"\n", device.Name);
            err = client.TransferPlayback(device.ID, true)
            if err != nil {
                return err
            } else {
                return nil
            }
		}
	}

	return nil
}

func getAlbums(cmd *cobra.Command, args []string) error {
    limit := int(50)
    start := int(0)
    var opt *spotify.Options
    opt = &spotify.Options{
        Limit: &limit,
        Offset: &start,
    }
    for {
        albums, err := client.CurrentUsersAlbumsOpt(opt)
        if err != nil {
            return err
        }
        for _, album := range albums.Albums {
            fmt.Printf("%s\n", album.Name)
        }
        if len(albums.Albums) < limit {
            return nil
        }
        start += limit
    }
    return nil
}

func getSongs(cmd *cobra.Command, args []string) error {
    limit := int(50)
    start := int(0)
    var opt *spotify.Options
    opt = &spotify.Options{
        Limit: &limit,
        Offset: &start,
    }
    for {
        songs, err := client.CurrentUsersTracksOpt(opt)
        if err != nil {
            return err
        }
        for _, song := range songs.Tracks {
            fmt.Printf("%s\n", song.Name)
        }
        if len(songs.Tracks) < limit {
            return nil
        }
        start += limit
    }
    return nil
}

func getArtists(cmd *cobra.Command, args []string) error {
    start := ""
    for {
        artists, err := client.CurrentUsersFollowedArtistsOpt(50, start)
        if err != nil {
            return err
        }
        for _, artist := range artists.Artists {
            fmt.Printf("%s\n", artist.Name)
        }
        if len(artists.Artists) == 50 {
            start = string(artists.Artists[49].URI)
            start = strings.Split(start,":")[2]
            fmt.Printf("%s\n", start)
        } else {
            return nil
        }
    }
    return nil
}

func shuffle(cmd *cobra.Command, args []string) error {
	state, err := client.PlayerState()
	if err != nil {
		return err
	}

	return client.Shuffle(!state.ShuffleState)
}

func repeat(cmd *cobra.Command, args []string) error {
	state, err := client.PlayerState()
	if err != nil {
		return err
	}

	var repeat string
	if state.RepeatState == "off" {
		repeat = "track"
	} else if state.RepeatState == "track" {
		repeat = "context"
	} else if state.RepeatState == "context" {
		repeat = "off"
	} else {
		return fmt.Errorf("unsupported repeat state %s", state.RepeatState)
	}

	opt := &spotify.PlayOptions{
		DeviceID: findDeviceByName(deviceNameFlag),
	}
	return client.RepeatOpt(repeat, opt)
}

func play(cmd *cobra.Command, args []string) error {
	var (
		opt = &spotify.PlayOptions{}
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

	opt.DeviceID = findDeviceByName(deviceNameFlag)

	return client.PlayOpt(opt)
}

func vol(cmd *cobra.Command, args []string) error {
	state, err := client.PlayerState()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		fmt.Printf("Current volume is %d%%.\n", state.Device.Volume)
		return nil
	}

	var currVolume int
	switch vol := args[0]; vol {
	case "up":
		currVolume = state.Device.Volume + 10
	case "down":
		currVolume = state.Device.Volume - 10
	default:
		currVolume, err = strconv.Atoi(vol)
		if err != nil {
			return err
		}
	}

	if currVolume < 0 {
		currVolume = 0
	}

	if currVolume > 100 {
		currVolume = 100
	}

	opt := &spotify.PlayOptions{
		DeviceID: findDeviceByName(deviceNameFlag),
	}
	return client.VolumeOpt(currVolume, opt)
}

func pause(cmd *cobra.Command, args []string) error {
	opt := &spotify.PlayOptions{
		DeviceID: findDeviceByName(deviceNameFlag),
	}
	return client.PauseOpt(opt)
}

func next(cmd *cobra.Command, args []string) error {
	opt := &spotify.PlayOptions{
		DeviceID: findDeviceByName(deviceNameFlag),
	}
	return client.NextOpt(opt)
}

func prev(cmd *cobra.Command, args []string) error {
	opt := &spotify.PlayOptions{
		DeviceID: findDeviceByName(deviceNameFlag),
	}
	return client.PreviousOpt(opt)
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

		fmt.Printf("Spotify is currently playing on %s.\n", state.Device.Name)
		fmt.Printf("Artist: %s\n", strings.Join(artists, ", "))
		fmt.Printf("Album: %s\n", state.Item.Album.Name)
		fmt.Printf("Track: %s\n", state.Item.Name)
		fmt.Printf("Position: %s / %s\n", durationToStr(state.Progress), durationToStr(state.Item.Duration))
	} else {
		fmt.Println("Spotify is currently paused.")
	}

	return nil
}

// findDeviceByName finds the device by name.
// If name is empty, the first Computer device ID is returned if it's available;
// otherwise it returns the first device ID.
func findDeviceByName(name string) *spotify.ID {
	devices, err := client.PlayerDevices()
	if err != nil {
		return nil
	}

	for _, device := range devices {
		if name != "" && device.Name == name {
			return &device.ID
		} else if device.Type == "Computer" {
			return &device.ID
		}
	}

	if len(devices) > 0 {
		return &devices[0].ID
	}

	return nil
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
