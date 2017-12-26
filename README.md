# Spotctl

`spotctl` is command-line interface to control Spotify from your favorite terminal.

## Demo

One of the highlights is that `spotctl player` shows a real-time Spotify player that allows you to control it, right in your terminal!

[![asciicast](https://asciinema.org/a/154262.png)](https://asciinema.org/a/154262)

## Installation

## Homebrew

If you're on a Mac, you can install with [Homebrew](https://brew.sh/):

```
brew install https://raw.githubusercontent.com/jingweno/spotctl/master/Formula/spotctl.rb
```

## Download

You can download the latest release for your operating system [here](https://github.com/jingweno/spotctl/releases).

## Manual Instllation

`spotctl` needs to connect to Spotify's API in order to control it.
To manually build it, you first need to sign up (or into) Spotify's developer site and [create an Application](https://developer.spotify.com/my-applications/#!/applications/create).
Once you've done so, you can find its Client ID and Client Secret values and run the following command:

```
SPOTIFY_CLIENT_ID=XXX SPOTIFY_CLIENT_SECRET=YYY ./bin/build
```

## Running

**Please make sure the Spotify app is opened before running any `spotctl` commands**, since it talks to the Spotify API which in turns talks to the Spotify app in your local box.
Here is a list of available commands:

```
$ spotctl -h
A command-line interface to Spotify.

Usage:
  spotctl [command]

Available Commands:
  help        Help about any command
  login       Login with your Spotify credentials
  logout      Clear your local Spotify credentials
  next        Skip to the next track
  pause       Pause Spotify playback
  play        Resume playback or play a track, album, artist or playlist by name
  player      Show the live player panel
  prev        Return to the previous track
  repeat      Toggle repeat playback mode
  shuffle     Toggle shuffle playback mode
  status      Show the current player status
  version     Show version.
  vol         Set or return volume percentage

Flags:
  -h, --help   help for spotctl

Use "spotctl [command] --help" for more information about a command.
```

## License

[MIT](https://github.com/jingweno/spotctl/blob/master/LICENSE)
