package model

import (
	"github.com/webservice/api"
	"io/ioutil"
	"log"
	"net/http"
)

type Person struct {
	Name    string
	Age     int
	Address string
}

func NewPerson() *Person {
	person := &Person{
		Name:    "jared",
		Age:     19,
		Address: "BJ",
	}

	return person
}

func (p *Person) Showtime() {
	log.Println("Person: Welcome access your home!")
}

func (p *Person) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	body, err := ioutil.ReadAll(r.Body)
	api.CheckErr(err)
	defer r.Body.Close()

	log.Println("body:", string(body))
	w.Write(append([]byte("Person:"), body...))
	log.Println("Leave...")
}
