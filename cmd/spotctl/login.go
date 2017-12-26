package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with your Spotify credentials",
	RunE:  login,
}

func login(cmd *cobra.Command, args []string) error {
	state, err := generateRandomString(32)
	if err != nil {
		return err
	}

	ch := make(chan *oauth2.Token)

	http.Handle("/callback", &authHandler{state: state, ch: ch, auth: auth})
	go http.ListenAndServe("localhost:8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	tok := <-ch

	if err := saveToken(tok); err != nil {
		return err
	}

	return nil
}

type authHandler struct {
	state string
	ch    chan *oauth2.Token
	auth  spotify.Authenticator
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tok, err := a.auth.Token(a.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != a.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, a.state)
	}

	fmt.Fprintf(w, "Login successfully. Please return to your terminal.")

	a.ch <- tok
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
