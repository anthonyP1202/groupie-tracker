// package main

// import (
// 	"log"
// 	"net/http"
// 	"text/template"
// )

// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		HomePage(w, r)
// 	})
// 	fs := http.FileServer(http.Dir("static/"))
// 	http.Handle("/static/", http.StripPrefix("/static", fs))
// 	http.ListenAndServe(":8080", nil)
// }

// func HomePage(w http.ResponseWriter, r *http.Request) {
// 	template, err := template.ParseFiles("page/HomePage.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	template.Execute(w, nil)
// }

package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"fmt"
	"html/template"
	"net/http"
	// "unicode"
)

var tpl *template.Template
var db *sql.DB

func maino() {
	tpl, _ = template.ParseGlob("page/*.html")
	var err error
	db, err = sql.Open("sqlite3", "bdd.sqlite3")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", logintest)
	http.HandleFunc("/register", registerHandler)
	// http.HandleFunc("/registerauth", registerAuthHandler)
	http.ListenAndServe("localhost:5500", nil)
}

// loginHandler serves form for users to login with
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func logintest(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	fmt.Println(username)

	basedd := "bdd.sqlite3"
	theDB, err := sql.Open("sqlite3", basedd)
	defer theDB.Close()

	if err != nil {
		log.Fatal(err)
	}
	query := "SELECT * FROM USER"
	rows, _ := theDB.Query(query)
	var id int
	var pseudo string
	var email string
	var password string

	for rows.Next() {
		rows.Scan(&id, &pseudo, &email, &password)
		fmt.Println(strconv.Itoa(id) + " : " + pseudo + "  " + email + "    " + password)
	}
	tpl.ExecuteTemplate(w, "login.html", "check username and password")

}

// loginAuthHandler authenticates user login
// func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("*****loginAuthHandler running*****")
// 	r.ParseForm()
// 	username := r.FormValue("username")
// 	password := r.FormValue("password")
// 	fmt.Println("username:", username, "password:", password)

// 	row, _ := db.Query("SELECT pseudo FROM USER;")
// 	fmt.Println("")
// 	fmt.Println("")
// 	fmt.Println(row)
// 	fmt.Println("")
// 	fmt.Println("")
// 	fmt.Println("")

// 	err := row.Scan()
// 	if err != nil {
// 		tpl.ExecuteTemplate(w, "login.html", "check username and password")
// 		return
// 	}
// 	// func CompareHashAndPassword(hashedPassword, password []byte) error
// 	verif := password == "s"
// 	// returns nill on succcess
// 	if verif {
// 		fmt.Fprint(w, "You have successfully logged in :)")
// 		return
// 	}
// 	fmt.Println("incorrect password")
// 	tpl.ExecuteTemplate(w, "login.html", "check username and password")
// }

// registerHandler serves form for registring new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// // registerAuthHandler creates new user in database
// func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
// 	/*
// 		1. check username criteria
// 		2. check password criteria
// 		3. check if username is already exists in database
// 		4. create bcrypt hash from password
// 		5. insert username and password hash in database
// 		(email validation will be in another video)
// 	*/
// 	fmt.Println("*****registerAuthHandler running*****")
// 	r.ParseForm()
// 	username := r.FormValue("username")
// 	// check username for only alphaNumeric characters
// 	var nameAlphaNumeric = true
// 	for _, char := range username {
// 		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
// 		// if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
// 		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
// 			nameAlphaNumeric = false
// 		}
// 	}
// 	// check username pswdLength
// 	var nameLength bool
// 	if 5 <= len(username) && len(username) <= 50 {
// 		nameLength = true
// 	}
// 	// check password criteria
// 	password := r.FormValue("password")
// 	fmt.Println("password:", password, "\npswdLength:", len(password))
// 	// variables that must pass for password creation criteria
// 	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
// 	pswdNoSpaces = true
// 	for _, char := range password {
// 		switch {
// 		// func IsLower(r rune) bool
// 		case unicode.IsLower(char):
// 			pswdLowercase = true
// 		// func IsUpper(r rune) bool
// 		case unicode.IsUpper(char):
// 			pswdUppercase = true
// 		// func IsNumber(r rune) bool
// 		case unicode.IsNumber(char):
// 			pswdNumber = true
// 		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
// 		case unicode.IsPunct(char) || unicode.IsSymbol(char):
// 			pswdSpecial = true
// 		// func IsSpace(r rune) bool, type rune = int32
// 		case unicode.IsSpace(int32(char)):
// 			pswdNoSpaces = false
// 		}
// 	}
// 	if 11 < len(password) && len(password) < 60 {
// 		pswdLength = true
// 	}
// 	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
// 	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
// 		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
// 		return
// 	}
// 	// check if username already exists for availability
// 	stmt := "SELECT UserID FROM bcrypt WHERE username = ?"
// 	row := db.QueryRow(stmt, username)
// 	var uID string
// 	err := row.Scan(&uID)
// 	if err != sql.ErrNoRows {
// 		fmt.Println("username already exists, err:", err)
// 		tpl.ExecuteTemplate(w, "register.html", "username already taken")
// 		return
// 	}
// 	// create hash from password
// 	var hash []byte
// 	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
// 	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Println("bcrypt err:", err)
// 		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
// 		return
// 	}
// 	fmt.Println("hash:", hash)
// 	fmt.Println("string(hash):", string(hash))
// 	// func (db *DB) Prepare(query string) (*Stmt, error)
// 	var insertStmt *sql.Stmt
// 	insertStmt, err = db.Prepare("INSERT INTO bcrypt (Username, Hash) VALUES (?, ?);")
// 	if err != nil {
// 		fmt.Println("error preparing statement:", err)
// 		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
// 		return
// 	}
// 	defer insertStmt.Close()
// 	var result sql.Result
// 	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
// 	result, err = insertStmt.Exec(username, hash)
// 	rowsAff, _ := result.RowsAffected()
// 	lastIns, _ := result.LastInsertId()
// 	fmt.Println("rowsAff:", rowsAff)
// 	fmt.Println("lastIns:", lastIns)
// 	fmt.Println("err:", err)
// 	if err != nil {
// 		fmt.Println("error inserting new user")
// 		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
// 		return
// 	}
// 	fmt.Fprint(w, "congrats, your account has been successfully created")
// }
