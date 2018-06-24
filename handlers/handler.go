package handlers

import (
	"github.com/evil-router/isfired/models"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"github.com/oschwald/geoip2-golang"
	"net"

)

type person struct {
	Name   string
	Reason string
}

type response struct {
	Host   string
	Source string
	City   string
}

func getRequest(r *http.Request) response {
	var req response
	if r.Header.Get("X-Forwarded-Server") != "" {
		req.Host = r.Header.Get("X-Forwarded-Host")
		req.Source = r.Header.Get("X-Forwarded-For")
	} else {
		req.Host = r.Host
		req.Source = r.RemoteAddr
	}
	db, err := geoip2.Open("GeoLite2/GeoLite2-City.mmdb")
	if err != nil {
		log.Print(err)
	}
	log.Printf(" city %v", req.Source)
	ip := net.ParseIP(req.Source)
	res, err := db.City(ip)
	if err != nil {
		log.Print(err)
	}
	log.Printf(" city %v", res.City.Names)
	req.City = res.City.Names["en"]
	return req
}

func Default(w http.ResponseWriter, r *http.Request) {
	req := getRequest(r)
	t, _ := template.ParseFiles("./tmpl/welcome.html")
	s, err := models.GetComment(req.Host, 10, 0)
	if err != nil {
		log.Print(err)
	}
	err = t.Execute(w, s) //step 2
	if err != nil {
		log.Print(err)
	}
	log.Println("Req: %v", req)
	log.Println("Getting status: %v", s)

}

func Seter(w http.ResponseWriter, r *http.Request) {
	req := getRequest(r)

	param := r.URL.Query()
	comment := bluemonday.UGCPolicy().Sanitize(param.Get("comment"))
	status,err := strconv.ParseBool(param.Get("status"))
	if err != nil{
		status = false
	}
	//key := param.Get("key")

	models.SetComment(req.Host,comment,req.City,status)

	Default(w,r)
}
func History(w http.ResponseWriter, r *http.Request) {
	req := getRequest(r)
	t, _ := template.ParseFiles("./tmpl/history.html")
	s, err := models.GetComment(req.Host, 10, 0)
	if err != nil {
		log.Print(err)
	}
	err = t.Execute(w, s) //step 2
	if err != nil {
		log.Print(err)
	}
	log.Println("Creating a new connection: %v", s)

}
