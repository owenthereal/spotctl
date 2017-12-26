# Spotctl

`spotctl` is command-line interface to control Spotify from your favorite terminal.

## Demo

[![asciicast](https://asciinema.org/a/154262.png)](https://asciinema.org/a/154262)

## Installation

## Homebrew

If you're on a Mac, you can install with [Homebrew](https://brew.sh/) like:

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

## License

[MIT](https://github.com/jingweno/spotctl/blob/master/LICENSE)
