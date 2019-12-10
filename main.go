package main

import (
	"github.com/webservice/api"
	"github.com/webservice/model"
	"log"
	"net/http"
)

var requestCount int

func main() {
	person := model.NewPerson()
	person.Showtime()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		log.Println("Main: Accept requrst", r.URL.String(), requestCount)
		w.Header().Add("Cache-Control", "no-store")
		w.Write([]byte("Welcome you..."))
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	http.Handle("/person", person)

	log.Println("http://127.0.0.1:8889/")
	api.CheckErr(http.ListenAndServe(":8889", nil))
}
