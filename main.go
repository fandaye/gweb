package main

import (
	"net/http"
	"log"
	"github.com/fandaye/gweb/app"
)

func main() {
	g := app.Global{}
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/login", g.Login)
	http.HandleFunc("/", g.Index)

	start_log := http.ListenAndServe("0.0.0.0:80", nil)
	log.Println(start_log)
}
