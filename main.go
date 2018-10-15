package main

import (
	"net/http"
	"log"
	"github.com/fandaye/gweb/app"
	"strconv"
)

func main() {
	G := app.Global{}
	G.ConfigFile = "config/config.ini"
	G.GlobalConfig = G.Config.GlobalConfig()
	G.DB.Config = G.Config.MysqlConfig() // 初始化mysql配置


	G.Redis.Config = G.Config.RedisConfig() // 初始化redis 配置
	G.Redis.Connect()                       // 初始化redis 连接

	if CookieExpiration, err := strconv.Atoi(G.GlobalConfig["CookieExpiration"]); err == nil {
		G.CookieExpiration = CookieExpiration
	} else {
		log.Println("CookieExpiration 参数配置错误")
		return
	}

	if Value, err := strconv.Atoi(G.GlobalConfig["max_login_error_number"]); err == nil {
		G.MaxLoginErrorNumber = Value
	} else {
		log.Println("max_login_error_number 参数配置错误")
		return
	}

	if Value, err := strconv.Atoi(G.GlobalConfig["login_error_lock_time"]); err == nil {
		G.LoginErrorLockTime = Value
	} else {
		log.Println("login_error_lock_time 参数配置错误")
		return
	}


	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/login", G.AuthHandler(G.Login))
	http.HandleFunc("/logout", G.Logout)
	http.HandleFunc("/", G.AuthHandler(G.Index))
	http.HandleFunc("/ldap", G.AuthHandler(G.Ldap))
	http.HandleFunc("/sys", G.AuthHandler(G.Sys))
	http.HandleFunc("/pubilc", G.AuthHandler(G.Public))

	ADDR := G.GlobalConfig["bind"] + ":" + G.GlobalConfig["port"]
	log.Println("开始启动监听:" + ADDR)
	start_log := http.ListenAndServe(ADDR, nil)
	log.Println(start_log)
}
