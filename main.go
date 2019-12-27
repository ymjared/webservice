package main

import (
	"fmt"
	"github.com/webservice/api"
	clog "github.com/webservice/example/log"
	"github.com/webservice/model"
	"net/http"
	"time"
)

var requestCount int
var Logger *clog.Logger

func init() {
	Logger = clog.CreateLogger()
}

func main() {
	Logger.Log(time.Now().String())
	person := model.NewPerson()
	person.Showtime()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		Logger.Log(fmt.Sprintln("Main: Accept requrst", r.URL.String(), requestCount))
		w.Header().Add("Cache-Control", "no-store")
		w.Write([]byte("Welcome you..."))
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	http.Handle("/person", person)

	Logger.Log(fmt.Sprintln("http://127.0.0.1:8889/"))
	api.CheckErr(http.ListenAndServe(":8889", nil))
}
