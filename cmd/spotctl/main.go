package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	redirectURI = "http://localhost:10028/callback"
)

var (
	spotifyClientID     string
	spotifyClientSecret string
	version             string
)

var (
	auth      spotify.Authenticator
	token     *oauth2.Token
	client    spotify.Client
	tokenPath string
)

var rootCmd = &cobra.Command{
	Use:               "spotctl",
	Short:             "A command-line interface to Spotify.",
	PersistentPreRun:  preRootCmd,
	PersistentPostRun: postRootCmd,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version.",
	Run:   ver,
}

func ver(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}

func main() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(volCmd)
	rootCmd.AddCommand(shuffleCmd)
	rootCmd.AddCommand(repeatCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(playerCmd)
	rootCmd.AddCommand(versionCmd)

	playCmd.PersistentFlags().StringVarP(&playCmdFlagType, "type", "t", "track", "the type of [name] to play: track, album, artist or playlist.")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func preRootCmd(cmd *cobra.Command, args []string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	tokenPath = filepath.Join(usr.HomeDir, ".spotctl")
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
	)
	auth.SetAuthInfo(spotifyClientID, spotifyClientSecret)

	// skip reading token or login if this is a login command
	if cmd.Use == "login" {
		return
	}

	token, err = readToken()
	if err != nil {
		if os.IsNotExist(err) {
			if err := login(cmd, args); err != nil {
				log.Fatal(err)
			}

			// read token one more time
			token, err = readToken()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	client = auth.NewClient(token)
}

func postRootCmd(cmd *cobra.Command, args []string) {
	// skip reading token or login if this is a login command
	if cmd.Use == "login" {
		return
	}

	tokenInUse, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}

	if tokenInUse != token {
		if err := saveToken(tokenInUse); err != nil {
			log.Fatal(err)
		}
	}
}

func saveToken(tok *oauth2.Token) error {
	f, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(tok)
}

func readToken() (*oauth2.Token, error) {
	content, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	if err := json.Unmarshal(content, &tok); err != nil {
		return nil, err
	}

	return &tok, nil
}
