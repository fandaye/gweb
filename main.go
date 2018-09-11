package main

import (
	"net/http"
	"log"
	"github.com/fandaye/gweb/app"
)

func main() {
	g := app.Global{}
	g.ConfigFile = "config/config.ini"

	g.GlobalConfig = g.Config.GlobalConfig()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/login", g.AuthHandler(g.Login))
	http.HandleFunc("/logout", g.Logout)
	http.HandleFunc("/", g.AuthHandler(g.Index))

	http.HandleFunc("/ldap", g.AuthHandler(g.Ldap))

	ADDR := g.GlobalConfig["bind"] + ":" + g.GlobalConfig["port"]
	log.Println("开始启动监听:" + ADDR)
	start_log := http.ListenAndServe(ADDR, nil)
	log.Println(start_log)
}
