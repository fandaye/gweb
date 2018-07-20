package app

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
)

func (G *Global) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/login/index.html"); err == nil {
			t.Execute(w, ' ')
		} else {
			log.Println(err)
			fmt.Fprintln(w, err)
		}
	} else if r.Method == "POST" {
		var data map[string]string
		email := r.PostFormValue("email")
		passwd := r.PostFormValue("passwd")
		G.Username = email
		if email == "admin" {
			if passwd == "123456" {
				resp, _ := getJson(1, "", data)
				fmt.Fprintf(w, string(resp))
			} else {
				resp, _ := getJson(-1, "密码错误", data)
				fmt.Fprintf(w, string(resp))
			}
		} else {
			resp, _ := getJson(-1, "账户错误", data)
			fmt.Fprintf(w, string(resp))
		}
	}
}
