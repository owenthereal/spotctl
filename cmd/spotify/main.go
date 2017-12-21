package main

import (
	"encoding/json"
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
	redirectURI  = "http://localhost:8080/callback"
	clientID     = "047627e6515f464a932fb6c4c6a1a446"
	clientSecret = "ac19b6260f7741dc800d8aac867871d4"
)

var (
	auth      spotify.Authenticator
	token     *oauth2.Token
	client    spotify.Client
	tokenPath string
)

var rootCmd = &cobra.Command{
	Use:               "spotify",
	Short:             "A command-line interface to Spotify.",
	PersistentPreRun:  preRootCmd,
	PersistentPostRun: postRootCmd,
}

func main() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(statusCmd)

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

	tokenPath = filepath.Join(usr.HomeDir, ".spotify")
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
	)
	auth.SetAuthInfo(clientID, clientSecret)

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
