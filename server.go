package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"text/template"

	lyrics "github.com/rhnvrm/lyric-api-go"
	"golang.org/x/oauth2/clientcredentials"
)

type test struct {
	Id         int
	Lyrics     string
	PreviewURL string
	OtherMusic string
	Music      *spotify.PlaylistTrackPage
	Artiste    []spotify.SimpleArtist
	Title      string
}

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     "2243f558d2644e81a6b121bd763acd00",
		ClientSecret: "bb0c6c22f67b4440b1f506673f0d6a32",
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	playlists, err := client.GetPlaylist(ctx, "6FPetUGNfFzaltVi4omGH0")
	fmt.Println(playlists.Name + "*")
	fmt.Println(playlists.ID + "*")
	if err != nil {
		log.Fatal(err)
	}
	playlistTrack, err := client.GetPlaylistTracks(ctx, playlists.ID)
	if err != nil {
		log.Fatalln(err)
	}

	rnd := 0
	if playlistTrack.Total < playlistTrack.Limit {
		rnd = rand.Intn(playlistTrack.Total)
	} else {
		rnd = rand.Intn(playlistTrack.Limit)
	}
	fmt.Println(playlistTrack.Tracks[rnd].Track.PreviewURL)
	l := lyrics.New()
	lyric, err := l.Search("John Lennon", "Imagine")

	if err != nil {
		fmt.Printf("Lyrics for John Lennon - imagine were not found")
	}

	music := test{rnd, lyric, playlistTrack.Tracks[rnd].Track.PreviewURL, string(playlistTrack.Tracks[rnd].Track.ID), playlistTrack, playlistTrack.Tracks[rnd].Track.Artists, playlistTrack.Tracks[rnd].Track.Name}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HomePage(w, r, &music)
	})

	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		BlindTest(w, r, &music)
	})

	http.HandleFunc("/GuessSong", func(w http.ResponseWriter, r *http.Request) {
		guessSong(w, r, &music)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)

}

func HomePage(w http.ResponseWriter, r *http.Request, track *test) {

	template, err := template.ParseFiles("page/HomePage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, track)
}

func BlindTest(w http.ResponseWriter, r *http.Request, track *test) {
	template, err := template.ParseFiles("page/BlindTest.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("r.formvalue" + r.FormValue("letter"))
	if r.FormValue("letter") != "" {
		if compare(r.FormValue("letter"), track.Title) {
			fmt.Println("gg")
		} else {
			fmt.Println("you're a failure like me")
		}
	}

	// rndList := []int{}
	track.PreviewURL = ""
	for track.PreviewURL == "" {

		rnd := 0
		// contained := 0
		if track.Music.Total < track.Music.Limit {
			rnd = rand.Intn(track.Music.Total)
		} else {
			rnd = rand.Intn(track.Music.Limit)
		}
		track.PreviewURL = track.Music.Tracks[rnd].Track.PreviewURL

		track.OtherMusic = string(track.Music.Tracks[rnd].Track.ID)
		track.Artiste = (track.Music.Tracks[rnd].Track.Artists)
		track.Title = (track.Music.Tracks[rnd].Track.Name)
		fmt.Println(track.Title)
	}
	template.Execute(w, track)

	// for {
	// 	for contained == 0 {
	// 		if track.Music.Total < track.Music.Limit {
	// 			rnd = rand.Intn(track.Music.Total)
	// 		} else {
	// 			rnd = rand.Intn(track.Music.Limit)
	// 		}
	// 		contained = 1
	// 		fmt.Println(rnd)
	// 		for i := 0; i < len(rndList); i++ {
	// 			if rnd == rndList[i] {
	// 				contained = 0

	// 			}

	// 		}
	// 	}
	// 	contained = 0
	// }
}

func guessSong(w http.ResponseWriter, r *http.Request, track *test) {
	l := lyrics.New()
	template, err := template.ParseFiles("page/Guessong.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("r.formvalue" + r.FormValue("letter"))
	if r.FormValue("letter") != "" {
		if compare(r.FormValue("letter"), track.Title) {
			fmt.Println("gg")
		} else {
			fmt.Println("you're a failure like me")
		}
	}
	track.Lyrics = ""
	for track.Lyrics == "" {

		rnd := 0
		// contained := 0
		if track.Music.Total < track.Music.Limit {
			rnd = rand.Intn(track.Music.Total)
		} else {
			rnd = rand.Intn(track.Music.Limit)
		}
		track.Lyrics, err = l.Search(track.Music.Tracks[rnd].Track.Artists[0].Name, track.Music.Tracks[rnd].Track.Name)

		track.OtherMusic = string(track.Music.Tracks[rnd].Track.ID)
		track.Artiste = (track.Music.Tracks[rnd].Track.Artists)
		track.Title = (track.Music.Tracks[rnd].Track.Name)
		fmt.Println(track.Title)
	}
	template.Execute(w, track)
}

func compare(tocompare string, compareto string) bool {
	comparetobigger := 0
	fmt.Println(tocompare + " " + compareto)
	tocompare = strings.ToLower(tocompare)
	compareto = strings.ToLower(compareto)
	maxmistake := len(compareto) - (len(compareto) / 10)
	if len(tocompare) < len(compareto) {
		comparetobigger = 1
	}
	mistake := 0
	fmt.Println(len(compareto))
	fmt.Println(len(tocompare))
	if comparetobigger == 1 {
		for i := 0; i < len(tocompare); i++ {
			if compareto[i] != tocompare[i] {
				mistake++
			}
		}
	} else {
		for i := 0; i < len(compareto); i++ {
			if compareto[i] != tocompare[i] {
				mistake++
			}
		}
	}

	return maxmistake > mistake
}
