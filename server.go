package main

import (
	"context"
	"math/rand"
	"strings"

	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/gorilla/websocket"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// "unicode"

/**
************************************************* VARIABLES ***************************************************
**/

type Cookie struct {
	Name  string
	Value string
}

type track struct {
	Id         int
	Lyrics     string
	PreviewURL string
	OtherMusic string
	Music      *spotify.PlaylistTrackPage
	Artiste    []spotify.SimpleArtist
	Title      string
}

// type Users struct {
// 	Id       int
// 	Pseudo   string
// 	Email    string
// 	Password string
// }

var tpl *template.Template
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var clients []websocket.Conn

func main() {

	tpl, _ = template.ParseGlob("page/*.html")

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		clients = append(clients, *conn)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			fmt.Printf("%s send : %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range clients {
				if err = client.WriteMessage(msgType, msg); err != nil {
					return
				}
			}
		}
	})

	//set token and create user
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

	//create a placeholder track structure
	playlistTrack, err := client.GetPlaylistTracks(ctx, playlists.ID)
	if err != nil {
		log.Fatalln(err)
	}

	// rnd := 0
	// if playlistTrack.Total < playlistTrack.Limit {
	// 	rnd = rand.Intn(int(playlistTrack.Total))
	// } else {
	// 	rnd = rand.Intn(int(playlistTrack.Total))
	// }
	// fmt.Println(playlistTrack.Tracks[rnd].Track.PreviewURL)
	l := lyrics.New()
	lyric, err := l.Search("John Lennon", "Imagine")

	if err != nil {
		fmt.Printf("Lyrics for John Lennon - imagine were not found")
	}

	music := track{0, lyric, playlistTrack.Tracks[0].Track.PreviewURL, string(playlistTrack.Tracks[0].Track.ID), playlistTrack, playlistTrack.Tracks[0].Track.Artists, playlistTrack.Tracks[0].Track.Name}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/testing", testHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.HandleFunc("/Guessong", func(w http.ResponseWriter, r *http.Request) {
		Guessong(w, r, &music)
	})
	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		BlindTest(w, r, &music)
	})
	http.HandleFunc("/PetitBac", PetitBacHandler)
	http.HandleFunc("/temp", TempHandler)
	//.....................//
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe("localhost:8800", nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****homeHandler running*****")
	tpl.ExecuteTemplate(w, "index.html", nil)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****homeHandler running*****")
	tpl.ExecuteTemplate(w, "HomePage.html", nil)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "Log.html", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "Sign-in.html", nil)
}
func guessSongHandeler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "Guessong.html", nil)
}
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "PetitBac.html", nil)
}
func BlindTestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "BlindTest.html", nil)
}

func TempHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "temp.html", nil)
}

// loginAuthHandler authenticates user login
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")

	db, err := sql.Open("sqlite3", "bdd.db")

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()
	pseudo := r.Form.Get("Username")
	password := r.Form.Get("password")
	fmt.Println("pseudo:", pseudo, "password:", password)

	rows, _ := db.Query("SELECT * FROM USER;")

	var idDB int
	var pseudoDB string
	var passwordDB string
	var emailDB string
	for rows.Next() {
		rows.Scan(&idDB, &pseudoDB, &emailDB, &passwordDB)
		fmt.Println(strconv.Itoa(idDB) + " " + pseudoDB + " " + emailDB + " " + passwordDB)
		verifPSW := password == passwordDB
		verifPSEUDO := pseudo == pseudoDB
		verifMAIL := pseudo == emailDB

		if (verifPSW && verifPSEUDO) || (verifPSW && verifMAIL) {
			cookie := http.Cookie{
				Name: pseudo,
			}
			http.SetCookie(w, &cookie)
			tpl.ExecuteTemplate(w, "HomePage.html", nil)

			return
		}
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "Log.html", "check username and password")
}

// registerAuthHandler creates new user in database
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")

	// Open database connection
	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Ensure the database connection is closed when function returns

	// Parse form data
	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confipassword := r.FormValue("confipassword")
	fmt.Println("pseudo =", username, ", email =", email, ", password =", password, ", confipassword =", confipassword)

	// Check if passwords match
	if password != confipassword {
		fmt.Println("Les mots de passe ne correspondent pas")
		tpl.ExecuteTemplate(w, "Sign-in.html", "Les mots de passe ne correspondent pas")
		return
	}

	// Check if username already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM USER WHERE pseudo = ?", username).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		fmt.Println("Ce pseudo est déjà utilisé")
		tpl.ExecuteTemplate(w, "Sign-in.html", "Ce pseudo est déjà utilisé")
		return
	}

	// Check if email already exists
	err = db.QueryRow("SELECT COUNT(*) FROM USER WHERE email = ?", email).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		fmt.Println("Cet email est déjà utilisé")
		tpl.ExecuteTemplate(w, "Sign-in.html", "Cet email est déjà utilisé")
		return
	}

	// Insert new user into database
	_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		http.Error(w, "Failed to insert user into database", http.StatusInternalServerError)
		return
	}

	fmt.Println("Utilisateur ajouté")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func BlindTest(w http.ResponseWriter, r *http.Request, track *track) {
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
			rnd = rand.Intn(int(track.Music.Total))
		} else {
			rnd = rand.Intn(int(track.Music.Limit))
		}
		track.PreviewURL = track.Music.Tracks[rnd].Track.PreviewURL

		track.OtherMusic = string(track.Music.Tracks[rnd].Track.ID)
		track.Artiste = (track.Music.Tracks[rnd].Track.Artists)
		track.Title = (track.Music.Tracks[rnd].Track.Name)
		fmt.Println(track.Title)
	}
	template.Execute(w, track)
}

func Guessong(w http.ResponseWriter, r *http.Request, track *track) {
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
			rnd = rand.Intn(int(track.Music.Total))
		} else {
			rnd = rand.Intn(int(track.Music.Limit))
		}
		track.Lyrics, err = l.Search(track.Music.Tracks[rnd].Track.Artists[0].Name, track.Music.Tracks[rnd].Track.Name)

		track.OtherMusic = string(track.Music.Tracks[rnd].Track.ID)
		track.Artiste = (track.Music.Tracks[rnd].Track.Artists)
		track.Title = (track.Music.Tracks[rnd].Track.Name)
		fmt.Println(track.Title)
	}
	template.Execute(w, track)
}

func PetitBac(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/PetitBac.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Sign(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/Sign-in.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/Log.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

// à supprimer à la fin
func Temp(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/temp.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
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
