package main

import (
	"context"
	"math/rand"
	"regexp"
	"strings"

	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/gorilla/websocket"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

/**
************************************************* VARIABLES ***************************************************
**/

type PetitBacSettings struct {
	CurrentLetter string
	letterlist    []string
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

type answers struct {
	Artiste    string
	Album      string
	Groupe     string
	Instrument string
	Featuring  string
}

type Cookie struct {
	Name  string
	Value string
}

var tpl *template.Template

var clients []websocket.Conn
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/**
************************************************* FIN  VARIABLES ***************************************************
**/

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

	servBac := PetitBacSettings{"u", []string{}}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/testing", testHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/Guessong", GuessongHandler)
	http.HandleFunc("/BlindTest", BlindTestHandler)
	http.HandleFunc("/PetitBac", PetitBacHandler)

	http.HandleFunc("/GuessongGame", func(w http.ResponseWriter, r *http.Request) {
		leaderboardHandler(w, r)
		Guessong(w, r, &music)
	})
	http.HandleFunc("/BlindTestGame", func(w http.ResponseWriter, r *http.Request) {
		leaderboardHandler(w, r)
		BlindTest(w, r, &music)
	})
	http.HandleFunc("/PetitBacGame", func(w http.ResponseWriter, r *http.Request) {
		leaderboardHandler(w, r)
		PetitBac(w, r, &servBac)
	})

	http.HandleFunc("/PetitBacValidation", func(w http.ResponseWriter, r *http.Request) {
		PetitBacValidation(w, r, &servBac)
	})
	http.HandleFunc("/temp", TempHandler)
	http.HandleFunc("/createBlind", createCodeBlindTestHandler)
	http.HandleFunc("/createGuess", createCodeGuessHandler)
	http.HandleFunc("/createPTB", createCodePTBHandler)
	//.....................//
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe("localhost:8800", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cacao")
	// Supprimer le cookie
	cookie := http.Cookie{
		Name:    "pseudo",        // Nom du cookie à supprimer
		Value:   "",              // Effacer la valeur du cookie
		Expires: time.Unix(0, 0), // Rendre le cookie expiré
		MaxAge:  -1,              // Fixer le temps de vie négatif pour rendre le cookie expiré
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****homeHandler running*****")
	tpl.ExecuteTemplate(w, "index.html", nil)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	if _, err := r.Cookie("pseudo"); err == nil {
		loggedIn = true
	}

	data := struct {
		LoggedIn bool
	}{
		LoggedIn: loggedIn,
	}
	fmt.Println("*****homeHandler running*****")
	tpl.ExecuteTemplate(w, "HomePage.html", data)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "Log.html", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "Sign-in.html", nil)
}
func GuessongHandler(w http.ResponseWriter, r *http.Request) {
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
				Name:  "pseudo",
				Value: pseudo,
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "Log.html", "check username and password")

}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")

	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confipassword := r.FormValue("confipassword")
	fmt.Println("pseudo =", username, ", email =", email, ", password =", password, ", confipassword =", confipassword)

	if password != confipassword {
		errorMessage := "Les mots de passe ne correspondent pas"
		tpl.ExecuteTemplate(w, "Sign-in.html", map[string]interface{}{"Error": errorMessage})
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
		errorMessage := ("Ce pseudo est déjà utilisé")
		tpl.ExecuteTemplate(w, "Sign-in.html", map[string]interface{}{"Error": errorMessage})
		return
	}

	// Check if email already exists
	err = db.QueryRow("SELECT COUNT(*) FROM USER WHERE email = ?", email).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		errorMessage := ("Cet email est déjà utilisé")
		tpl.ExecuteTemplate(w, "Sign-in.html", map[string]interface{}{"Error": errorMessage})
		return
	}

	// Validate password CNIL
	isValid, message := validatePassword(password)
	if !isValid {
		errorMessage := "Mot de passe non valide: " + message
		tpl.ExecuteTemplate(w, "Sign-in.html", map[string]interface{}{"Error": errorMessage})
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
	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	template, err := template.ParseFiles("page/BlindTestInGame.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("r.formvalue" + r.FormValue("letter"))
	if r.FormValue("letter") != "" {
		if compare(r.FormValue("letter"), track.Title) {
			fmt.Println("gg")
			stmt := `
        UPDATE ROOM_USERS
        SET score = score + 1
        WHERE id_room = (
            SELECT id
            FROM ROOMS
            WHERE name = ?
        )
        AND id_user = (
            SELECT id
            FROM USER
            WHERE pseudo = ?
        )
    `
			cookie, err := r.Cookie("CodeRoom")
			if err != nil {
				log.Fatal(err)
			}
			roomName := cookie.Value
			cookie2, err := r.Cookie("pseudo")
			if err != nil {
				log.Fatal(err)
			}
			username := cookie2.Value
			_, err = db.Exec(stmt, roomName, username)
			if err != nil {
				http.Error(w, "Failed to increment score in database", http.StatusInternalServerError)
				return
			}

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
	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	l := lyrics.New()
	template, err := template.ParseFiles("page/GuessongInGame.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("r.formvalue" + r.FormValue("letter"))
	if r.FormValue("letter") != "" {
		if compare(r.FormValue("letter"), track.Title) {
			fmt.Println("gg")
			stmt := `
        UPDATE ROOM_USERS
        SET score = score + 1
        WHERE id_room = (
            SELECT id
            FROM ROOMS
            WHERE name = ?
        )
        AND id_user = (
            SELECT id
            FROM USER
            WHERE pseudo = ?
        )
    `
			cookie, err := r.Cookie("CodeRoom")
			if err != nil {
				log.Fatal(err)
			}
			roomName := cookie.Value
			cookie2, err := r.Cookie("pseudo")
			if err != nil {
				log.Fatal(err)
			}
			username := cookie2.Value
			_, err = db.Exec(stmt, roomName, username)
			if err != nil {
				http.Error(w, "Failed to increment score in database", http.StatusInternalServerError)
				return
			}
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

func PetitBac(w http.ResponseWriter, r *http.Request, setting *PetitBacSettings) {
	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	nbCoorect := 0
	defer db.Close()
	done := 0
	if r.FormValue("Album") == "on" {
		nbCoorect += 1

	}
	if r.FormValue("Artiste") == "on" {
		nbCoorect += 1
	}
	if r.FormValue("Groupe") == "on" {
		nbCoorect += 1
	}
	if r.FormValue("Instrument") == "on" {
		nbCoorect += 1
	}
	if r.FormValue("Featuring") == "on" {
		nbCoorect += 1
	}
	fmt.Println("NB" + strconv.Itoa(nbCoorect))
	stmt := `
	UPDATE ROOM_USERS
	SET score = score + ?
	WHERE id_room = (
		SELECT id
		FROM ROOMS
		WHERE name = ?
	)
	AND id_user = (
		SELECT id
		FROM USER
		WHERE pseudo = ?
	)
`
	cookie, err := r.Cookie("CodeRoom")
	if err != nil {
		log.Fatal(err)
	}
	roomName := cookie.Value
	cookie2, err := r.Cookie("pseudo")
	if err != nil {
		log.Fatal(err)
	}
	username := cookie2.Value
	_, err = db.Exec(stmt, nbCoorect, roomName, username)
	if err != nil {
		http.Error(w, "Failed to increment score in database", http.StatusInternalServerError)
		return
	}

	answers := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	letter := "a"

	for done != 1 {

		letter = answers[rand.Intn(len(answers))]
		done = 1
		for i := 0; i < len(setting.letterlist); i++ {
			if letter == setting.letterlist[i] {
				done = 0
			}
		}
	}
	setting.CurrentLetter = letter
	setting.letterlist = append(setting.letterlist, letter)
	template, err := template.ParseFiles("page/PetitBacInGame1.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, setting)
}

func PetitBacValidation(w http.ResponseWriter, r *http.Request, setting *PetitBacSettings) {
	data := answers{r.FormValue("artiste"), r.FormValue("Album"), r.FormValue("groupe"), r.FormValue("instrum"), r.FormValue("chanson")}

	template, err := template.ParseFiles("page/PetitBacInGame2.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
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

func Temp(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/temp.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func createCodeBlindTestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")

	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cookie, err := r.Cookie("pseudo")
	if err != nil {
		log.Fatal(err)
	}

	// Vérifier si le nom de la ROOMS existe déjà
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM ROOMS WHERE name = ?", code).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		fmt.Println("Le nom de la ROOMS existe déjà")
		http.Error(w, "Le nom de la ROOMS existe déjà", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO ROOMS (created_by, max_player, name, id_game) VALUES (?, ?, ?, ?)", cookie.Value, 4, code, 1)
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT id FROM USER WHERE pseudo = @pseudo"
	var userID int
	err = db.QueryRow(query, sql.Named("pseudo", cookie.Value)).Scan(&userID)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return
	}
	fmt.Println("FINTEST")

	// Récupérer l'ID de la dernière ligne insérée
	idroom, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(idroom)

	result2, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", idroom, userID, 0)
	if err != nil {
		log.Fatal(err)
	}
	idRoomUser, err := result2.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(idRoomUser)
	cookieCode := http.Cookie{
		Name:  "CodeRoom",
		Value: code,
	}
	http.SetCookie(w, &cookieCode)

	http.Redirect(w, r, "/BlindTestGame", http.StatusSeeOther)

}

func createCodeGuessHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")

	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cookie, err := r.Cookie("pseudo")
	if err != nil {
		log.Fatal(err)
	}

	// Vérifier si le nom de la ROOMS existe déjà
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM ROOMS WHERE name = ?", code).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		fmt.Println("Le nom de la ROOMS existe déjà")
		http.Error(w, "Le nom de la ROOMS existe déjà", http.StatusBadRequest)
		return
	}

	// Insérer les données dans la base de données
	result, err := db.Exec("INSERT INTO ROOMS (created_by, max_player, name, id_game) VALUES (?, ?, ?, ?)", cookie.Value, 4, code, 2)
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT id FROM USER WHERE pseudo = @pseudo"
	var userID int
	err = db.QueryRow(query, sql.Named("pseudo", cookie.Value)).Scan(&userID)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return
	}
	fmt.Println("FINTEST")

	// Récupérer l'ID de la dernière ligne insérée
	idroom, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(idroom)

	result2, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", idroom, userID, 0)
	if err != nil {
		log.Fatal(err)
	}
	idRoomUser, err := result2.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(idRoomUser)
	cookieCode := http.Cookie{
		Name:  "CodeRoom",
		Value: code,
	}
	http.SetCookie(w, &cookieCode)
	http.Redirect(w, r, "/GuessongGame", http.StatusSeeOther)

}
func createCodePTBHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")

	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cookie, err := r.Cookie("pseudo")
	if err != nil {
		log.Fatal(err)
	}

	// Vérifier si le nom de la ROOMS existe déjà
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM ROOMS WHERE name = ?", code).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		fmt.Println("Le nom de la ROOMS existe déjà")
		// Vous pouvez choisir de renvoyer un message d'erreur au client ou effectuer une redirection
		http.Error(w, "Le nom de la ROOMS existe déjà", http.StatusBadRequest)
		return
	}

	// Insérer les données dans la base de données
	result, err := db.Exec("INSERT INTO ROOMS (created_by, max_player, name, id_game) VALUES (?, ?, ?, ?)", cookie.Value, 4, code, 3)
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT id FROM USER WHERE pseudo = @pseudo"
	var userID int
	err = db.QueryRow(query, sql.Named("pseudo", cookie.Value)).Scan(&userID)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return
	}
	fmt.Println("FINTEST")

	// Récupérer l'ID de la dernière ligne insérée
	idroom, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(idroom)

	result2, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", idroom, userID, 0)
	if err != nil {
		log.Fatal(err)
	}
	idRoomUser, err := result2.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(idRoomUser)
	cookieCode := http.Cookie{
		Name:  "CodeRoom",
		Value: code,
	}
	http.SetCookie(w, &cookieCode)
	http.Redirect(w, r, "/PetitBacGame", http.StatusSeeOther)

}

func compare(tocompare string, compareto string) bool {
	// comparetobigger := 0
	fmt.Println(tocompare + " " + compareto)
	tocompare = strings.ToLower(tocompare)
	compareto = strings.ToLower(compareto)
	// maxmistake := len(compareto) - (len(compareto) / 10)
	// if len(tocompare) < len(compareto) {
	// 	comparetobigger = 1
	// }
	// mistake := 0
	fmt.Println(len(compareto))
	fmt.Println(len(tocompare))
	// if comparetobigger == 1 {
	// 	for i := 0; i < len(tocompare); i++ {
	// 		if compareto[i] != tocompare[i] {
	// 			mistake++
	// 		}
	// 	}
	// } else {
	// 	for i := 0; i < len(compareto); i++ {
	// 		if compareto[i] != tocompare[i] {
	// 			mistake++
	// 		}
	// 	}
	// }
	if compareto != tocompare {
		return false
	}
	// return maxmistake > mistake
	return true
}

type LeaderboardRow struct {
	Pseudo string
	Score  int
}

type User struct {
	ID     int
	Pseudo string
	Email  string
}

type RoomUser struct {
	UserID int
	Score  int
}

type LeaderboardEntry struct {
	User  User
	Score int
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "bdd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cookie, err := r.Cookie("CodeRoom")
	if err != nil {
		log.Fatal(err)
	}
	roomName := cookie.Value
	fmt.Println(roomName)
	// Requête SQL pour récupérer les données du leaderboard
	rows, err := db.Query(`
	SELECT u.id, u.pseudo, ru.score
	FROM ROOM_USERS ru
	INNER JOIN USER u ON ru.id_user = u.id
	INNER JOIN ROOMS r ON ru.id_room = r.id
	WHERE r.name = ?
	ORDER BY ru.score DESC
    `, roomName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Stockage des données du leaderboard dans une slice de LeaderboardEntry
	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var user User
		var score int
		err := rows.Scan(&user.ID, &user.Pseudo, &score)
		if err != nil {
			log.Fatal(err)
		}
		leaderboard = append(leaderboard, LeaderboardEntry{User: user, Score: score})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Utilisation d'un moteur de template pour générer le HTML
	tmpl, err := template.New("leaderboard").Parse(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Leaderboard</title>
            <link rel="stylesheet" type="text/css" href="./static/BT.css" />
        </head>
        <body>
            <div id="nexto">
                <div class="container">
                    <h1>Leaderboard</h1>
                    <table>
                        <tr>
                            <th>User</th>
                            <th>Score</th>
                        </tr>
                        {{range .}}
                        <tr>
                            <td>{{.User.Pseudo}}</td>
                            <td>{{.Score}}</td>
                        </tr>
                        {{end}}
                    </table>
                </div>
                
            </div>
        </body>
        </html>
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Exécution du template avec les données du leaderboard
	err = tmpl.Execute(w, leaderboard)
	if err != nil {
		log.Fatal(err)
	}
}

func validatePassword(password string) (bool, string) {
	// Longueur minimale
	if len(password) < 8 {
		return false, "Le mot de passe doit contenir au moins 8 caractères."
	}

	// Vérification de la présence de chiffres
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	if !hasDigit {
		return false, "Le mot de passe doit contenir au moins un chiffre."
	}

	// Vérification de la présence de lettres majuscules
	hasUpperCase, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUpperCase {
		return false, "Le mot de passe doit contenir au moins une lettre majuscule."
	}

	// Vérification de la présence de lettres minuscules
	hasLowerCase, _ := regexp.MatchString(`[a-z]`, password)
	if !hasLowerCase {
		return false, "Le mot de passe doit contenir au moins une lettre minuscule."
	}

	// Vérification de la présence de caractères spéciaux
	hasSpecialChar, _ := regexp.MatchString(`[!@#$%^&*()_+{}|:"<>?~]`, password)
	if !hasSpecialChar {
		return false, "Le mot de passe doit contenir au moins un caractère spécial."
	}

	return true, "Le mot de passe satisfait toutes les recommandations de sécurité."
}
