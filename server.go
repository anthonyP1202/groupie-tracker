package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"text/template"

	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     "2243f558d2644e81a6b121bd763acd00",
		ClientSecret: "bb0c6c22f67b4440b1f506673f0d6a32",
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	// fmt.Println(token)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	// var Realclient *spotify.Client
	client := spotify.New(httpClient)
	client.Play(ctx)

	// fmt.Println(Realclient.GetPlaylist(ctx, "2024"))

	// fmt.Println(client.PlayerCurrentlyPlaying(ctx))

	// playlist, err := client.CurrentUsersPlaylists(ctx)
	playlist, err := client.Search(ctx, "2024", spotify.SearchTypePlaylist)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("megatest", playlist.Playlists.Playlists[0])
	if playlist.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range playlist.Playlists.Playlists[0].Tracks.Endpoint {
			{
				fmt.Println("   ", item)
			}
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			HomePage(w, r)
		})
		fs := http.FileServer(http.Dir("static/"))
		http.Handle("/static/", http.StripPrefix("/static", fs))
		http.ListenAndServe(":8080", nil)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/HomePage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
