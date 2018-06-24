package handlers

import (
	"github.com/evil-router/isfired/models"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"log"
	"net/http"
	"github.com/oschwald/geoip2-golang"
	"net"

	"strings"
	"regexp"
)

type response struct {
	Host   string
	Source string
	City   string
}

func getRequest(r *http.Request) response {
	var req response
	if r.Header.Get("X-Forwarded-Server") != "" {
		req.Host = r.Header.Get("X-Forwarded-Host")
		list := r.Header.Get("X-Forwarded-For")
		req.Source = strings.Split(list,",")[0]
	} else {
		req.Host = r.Host
		req.Source,_,_ = net.SplitHostPort(r.RemoteAddr)
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
	s, err := models.GetComment(req.Host, 1, 0)
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
	t, _ := template.ParseFiles("./tmpl/set.html")
	fired := false
	req := getRequest(r)
	r.ParseForm()
	comment := bluemonday.UGCPolicy().Sanitize(r.Form.Get("comment") )
	status := r.Form.Get("status")
	if status == "on"{
		fired = true
	}
	log.Printf("form data %v ",r.Form.Encode() )
	//key := param.Get("key")
	if len(comment)  >0 {
		models.SetComment(req.Host, comment, req.City, fired)
		Default(w,r)
		return
	}
	s, err := models.GetComment(req.Host, 10, 0)
	err = t.Execute(w, s) //step 2
	if err != nil {
		log.Print(err)
	}
	log.Println("Req: %v", req)
	log.Println("Getting status: %v", s)
}
func History(w http.ResponseWriter, r *http.Request) {
	req := getRequest(r)
	t, _ := template.ParseFiles("./tmpl/history.html")
	s, err := models.GetComment(req.Host, 1000, 0)
	if err != nil {
		log.Print(err)
	}
	err = t.Execute(w, s) //step 2
	if err != nil {
		log.Print(err)
	}

}

func AddSite(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./tmpl/add.html")
	r.ParseForm()
	name := r.PostForm.Get("name")
	if len(name) >0 {
		log.Printf("name: %v", name)
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Print(err)
		}
		site := reg.ReplaceAllString(name,"") + ".isfired.com"
		models.AddSite(site,name)
		http.Redirect(w,r,"http://" + site , 302)
	}

	s,_ := models.GetActiveSites()
	log.Printf("sites %v ",s)
	err := t.Execute(w,s) //step 2
	if err != nil {
		log.Print(err)
	}
}
