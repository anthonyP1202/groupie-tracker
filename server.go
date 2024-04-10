package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HomePage(w, r)
	})
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./HomePage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
