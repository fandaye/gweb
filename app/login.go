package app

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
	"encoding/json"
	"time"
	"strconv"
)

func (G *Global) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/login/index.html"); err == nil {
			t.Execute(w, ' ')
		} else {
			log.Println("函数 Login GET :", err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
		}
	} else if r.Method == "POST" {

		G.Email = r.PostFormValue("email")
		G.AuthToken = G.RandStr(64)
		G.Passwd = G.PasswdMD5(r.PostFormValue("passwd"))
		G.LoginErrorNumber = 0

		if Value, err := strconv.Atoi(G.GlobalConfig["max_login_error_number"]); err == nil {
			G.MaxLoginErrorNumber = Value
		} else {
			log.Println("函数 Login :", err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
			return
		}

		if Value, err := strconv.Atoi(G.GlobalConfig["login_error_lock_time"]); err == nil {
			G.LoginErrorLockTime = Value
		} else {
			log.Println("函数 Login :", err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
			return
		}

		if ValueS, err := G.Redis.Get("login_error_count_" + G.Email); err == nil {
			if ValueI, err := strconv.Atoi(ValueS); err == nil {
				G.LoginErrorNumber = ValueI
				if ValueI >= G.MaxLoginErrorNumber {
					fmt.Fprintf(w, string(getJson(-1, "账户密码输入错误已达上限，请稍后再试!", nil)))
					return
				}
			} else {
				log.Println("函数 Login :", err)
			}
		}

		if UserInfo, Err := G.SelectOne("SELECT * FROM users WHERE email=? OR username=?", G.Email, G.Email); Err != nil { //查询用户信息
			log.Println("函数 Login 查询用户信息 :", Err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
		} else {
			G.Username = UserInfo["username"]
			if UserInfo["passwd"] == G.Passwd { // 判断密码是否正确
				if UserInfo["status"] == "0" { // 判断账户是否被禁用
					user_json, _ := json.Marshal(G.UserInfo)
					if err := G.Redis.Set(UserInfo["email"], string(user_json), G.CookieExpiration); err == nil {
						http.SetCookie(w, &http.Cookie{Name: "AuthToken", Value: G.UserInfo.AuthToken, Path: "/", MaxAge: 0})
						http.SetCookie(w, &http.Cookie{Name: "Username", Value: UserInfo["email"], Path: "/", MaxAge: 0})
						if _, Err := G.DB.Update("UPDATE users SET login_time=? WHERE email=? OR username=?", time.Now().Format("2006-01-02 15:04:05"), G.Email, G.Email); err == nil { //更新用户最后登录时间
							fmt.Fprintf(w, string(getJson(1, "", nil)))
						} else {
							log.Println("函数 Login 更新用户登录时间 :", Err)
							fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
						}
					} else {
						log.Println("函数 Login 写入Redis :", err.Error())
						fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
					}
				} else {
					fmt.Fprintf(w, string(getJson(-1, "账号被禁用!", nil)))
				}
			} else {
				if err := G.Redis.Set("login_error_count_"+G.Email, strconv.Itoa(G.LoginErrorNumber+1), G.LoginErrorLockTime); err == nil {
					fmt.Fprintf(w, string(getJson(-1, "账号或者密码错误!", nil)))
				} else {
					log.Println("函数 Login 写入Redis :", err.Error())
				}
			}
		}
	}
}

//用户注销登录
func (G *Global) Logout(w http.ResponseWriter, r *http.Request) {
	//G.Redis.Del(G.UserInfo.Email)                                          //删除redis cookie
	http.SetCookie(w, &http.Cookie{Name: "AuthToken", Path: "/", MaxAge: -1}) //设置浏览器cookie
	http.SetCookie(w, &http.Cookie{Name: "Username", Path: "/", MaxAge: -1})  //设置浏览器cookie
	http.Redirect(w, r, "/login", http.StatusFound)                           //返回登录页面
}
