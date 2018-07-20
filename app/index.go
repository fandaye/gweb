package app

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
)

func (G *Global) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/index/index.html"); err == nil {
			t.Execute(w, G)
		} else {
			log.Println(err)
			fmt.Fprintln(w, err)
		}
	}
}
