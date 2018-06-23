package handlers

import (
	"net/http"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"github.com/evil-router/isfired/models"
	"log"
)


type person struct {
	Name   string
	Reason string
}

func Default(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./tmpl/welcome.html")
	s,err := models.GetComment(r.Host,10,0)
	if err != nil {
		log.Print(err)
	}
	err = t.Execute(w, s)   //step 2
	if err != nil {
		log.Print(err)
	}
	log.Println("Creating a new connection: %v", s)

}

func Seter(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./welcome.html")
	param := r.URL.Query()

	person := new(person)
	person.Name = bluemonday.UGCPolicy().Sanitize(param.Get("foo"))
	person.Reason = param.Get("key")
	t.Execute(w, person) //step 2
}
	func History(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./tmpl/history.html")
	s,err := models.GetComment(r.Host,10,0)
	if err != nil {
	log.Print(err)
	}
	err = t.Execute(w, s)   //step 2
	if err != nil {
	log.Print(err)
	}
	log.Println("Creating a new connection: %v", s)

	}
