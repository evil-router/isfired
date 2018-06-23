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

type response struct {
	Host string
	Source string
}
func getRequest( r *http.Request) (response) {
   var req response
   if r.Header.Get( "X-Forwarded-Server" ) != "" {
   	req.Host = r.Header.Get( "X-Forwarded-Server" )
   	req.Source = r.Header.Get( "X-Forwarded-For" )
   } else {
	   req.Host = r.Host
	   req.Source = r.RemoteAddr
	}
	return req
}

func Default(w http.ResponseWriter, r *http.Request) {
	req := getRequest(r)
	t, _ := template.ParseFiles("./tmpl/welcome.html")
	s,err := models.GetComment(req.Host,10,0)
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
