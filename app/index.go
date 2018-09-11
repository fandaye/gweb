package app

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
)


func (G *Global) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/index/index.html", "templates/common/head.html", "templates/common/tail.html") ; err == nil{
			t.Execute(w, G)
		} else {
			log.Println( "函数 Index: " , err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
		}
	}
}
