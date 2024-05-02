package main

import (
	"database/sql"
	"grouptrack/webSocket"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"

	"fmt"
	"html/template"
	"net/http"
	// "unicode"
)

/**
************************************************* VARIABLES ***************************************************
**/

var tpl *template.Template
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/**
************************************************* MAIN CODE ***************************************************
**/

func main() {

	tpl, _ = template.ParseGlob("page/*.html")

	http.HandleFunc("/ws", webSocket.WebsocketHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)

	http.ListenAndServe("localhost:8800", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****homeHandler running*****")
	tpl.ExecuteTemplate(w, "HomePage.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// loginAuthHandler authenticates user login
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")

	db, err := sql.Open("sqlite3", "bdd.db")

	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()
	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")
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
			fmt.Fprint(w, "You have successfully logged in :)")
			return
		}
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
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
		tpl.ExecuteTemplate(w, "register.html", "Les mots de passe ne correspondent pas")
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
		tpl.ExecuteTemplate(w, "register.html", "Ce pseudo est déjà utilisé")
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
		tpl.ExecuteTemplate(w, "register.html", "Cet email est déjà utilisé")
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
