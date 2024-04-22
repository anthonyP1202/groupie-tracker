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
	http.HandleFunc("/BlindTest", func(w http.ResponseWriter, r *http.Request) {
		BlindTest(w, r)

	})

	http.HandleFunc("/Guessong", func(w http.ResponseWriter, r *http.Request) {
		Guessong(w, r)

	})

	http.HandleFunc("/PetitBac", func(w http.ResponseWriter, r *http.Request) {
		PetitBac(w, r)

	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/HomePage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func BlindTest(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/BlindTest.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Guessong(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/Guessong.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func PetitBac(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("page/PetitBac.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
