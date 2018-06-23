package main

import (
	"net/http"
	"github.com/evil-router/isfired/handlers"
	"github.com/evil-router/isfired/config"
	"log"
	"flag"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	confptr := flag.String("conf","conf.json", "config file location")
	flag.Parse()
	err := config.GetConfig(*confptr)

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handlers.Default)
	http.HandleFunc("/set", handlers.Seter)
	http.HandleFunc("/history", handlers.History)

	server.ListenAndServe()
}
