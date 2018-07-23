package app

import (
	"encoding/json"
	"github.com/fandaye/gweb/modules"
	"net/http"
	"fmt"
	"log"
	"crypto/md5"
	"encoding/hex"
	"time"
	"math/rand"
	"strconv"
)

type Global struct {
	GlobalConfig     map[string]string
	CookieExpiration int
	modules.Config
	modules.DB
	modules.Redis
	modules.UserInfo

}


type UserInfoUnmarshal struct {
	UserName  string `json:"UserName"`
	Passwd    string `json:"Passwd"`
	AuthToken string `json:"AuthToken"`
	Role      string    `json:"Role"`
	Prems []struct {
		MenuID string `json:"menu_id"`
	} `json:"Prems"`
}



func getJson(code int, msg string, data map[string]string) ([]byte) {
	type Message struct {
		Code int
		Msg  string
		Data map[string]string
	}
	if mess,err:= json.MarshalIndent(Message{code, msg, data}, "", "") ; err == nil{
		return mess
	}else {
		log.Println("getJson函数：",err)
		return nil
	}
}

func (G *Global) RandStr(num int) string {
	codes := "123456789abcdefghjkmnpqrstuvwxyzABCDEFGHKMNPQRSTUVWXYZ"
	codeLen := len(codes)
	data := make([]byte, num)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	return string(data)
}

func (G *Global) PasswdMD5(p string) (string) {
	h_m5 := md5.New()
	h_m5.Write([]byte(p))
	return hex.EncodeToString(h_m5.Sum(nil)) //密码MD5加密
}


func (G *Global) AuthHandler(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/favicon.ico"){ // 排除网址图标
			f(w,r)
			return
		}


		if CookieExpiration, err := strconv.Atoi(G.GlobalConfig["CookieExpiration"]); err == nil {
			G.CookieExpiration = CookieExpiration
		} else {
			log.Println("函数 AuthHandler CookieExpiration:",err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
			return
		}


		G.DB.Config = G.Config.MysqlConfig() // 初始化mysql配置
		if err := G.DB.Connect(); err != nil { // 初始化MYSQL 连接
			log.Println("函数 AuthHandler 初始化mysql连接失败:",err)
			fmt.Fprintf(w, string(getJson(-1, "系统错误", nil)))
			return
		}
		defer G.DB.Conn.Close() //关闭Mysql 连接


		G.Redis.Config = G.Config.RedisConfig() // 初始化redis 配置
		G.Redis.Connect() // 初始化redis 连接
		defer G.Redis.Conn.Close() // 关闭redis 连接

		if (r.URL.Path == "/login") { // 对登录页面不验证cookie
			f(w,r)
			return
		}

		CookieAuthToken, CookieAuthTokenErr := r.Cookie("AuthToken") // 获取cookie 失败重定向到登录页面
		CookieUsername, CookieUsernameErr := r.Cookie("Username")
		if CookieAuthTokenErr != nil || CookieUsernameErr != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}else {
			if UserInfo, err := G.Redis.Get(CookieUsername.Value); err == nil { // 获取redis 信息
				var JsonRes UserInfoUnmarshal
				json.Unmarshal([]byte(UserInfo), &JsonRes)
				if CookieAuthToken.Value == JsonRes.AuthToken {
					G.Redis.ExPire(CookieUsername.Value, G.CookieExpiration)
					f(w,r)
				}else {
					fmt.Fprintf(w, string(getJson(-1, "账户异地登录，被迫下线!", nil)))
				}
			}else {
				http.SetCookie(w, &http.Cookie{Name: "AuthToken", Path: "/", MaxAge: -1})
				http.Redirect(w, r, "/login", http.StatusFound)
			}
		}
	}
}
