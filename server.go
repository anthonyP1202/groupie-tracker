package main

import (
	"database/sql"

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

type Cookie struct {
	Name  string
	Value string
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	// set cookie for storing token
	cookie := http.Cookie{}
	cookie.Name = "accessToken"
	cookie.Value = "ro8BS6Hiivgzy8Xuu09JDjlNLnSLldY5"
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "This is cookies!\n")
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

/**
************************************************* MAIN CODE ***************************************************
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "page/index.html")
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)

	println("Your server run 8080")
	http.ListenAndServe(":8080", nil)
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
			cookie := http.Cookie{
				Name:  "accessToken",
				Value: "ro8BS6Hiivgzy8Xuu09JDjlNLnSLldY5",
			}
			http.SetCookie(w, &cookie)
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

// func getUser(id int) string {
// 	db, err := sql.Open("sqlite3", "bdd.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var username string
// 	defer db.Close() // Ensure the database connection is closed when function returns
// 	username, err = db.Exec("SELECT pseudo FROM USER WHERE id =?", 5)
// 	return

// }
