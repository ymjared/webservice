package main

import (
	"github.com/webservice/api"
	"github.com/webservice/model"
	"log"
	"net/http"
)

func main() {
	person := model.NewPerson()
	person.Showtime()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Main: Accept requrst", r.URL.String())
		w.Write([]byte("Welcome you..."))
	})
	http.Handle("/person", person)

	log.Println("http://127.0.0.1:8889/")
	api.CheckErr(http.ListenAndServe(":8889", nil))
}
