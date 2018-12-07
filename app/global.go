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
)

type Global struct {
	Code                string
	GlobalConfig        map[string]string
	CookieExpiration    int // cookie 过期时间
	MaxLoginErrorNumber int // 最大登录失败次数
	LoginErrorLockTime  int // 登录失败锁定时间

	modules.Config
	modules.DB
	modules.Redis
	modules.UserInfo
	modules.Ldaps
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

func getJson(code int, msg string, data []map[string]string) ([]byte) {
	type Message struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data []map[string]string `json:"data"`
	}
	if mess, err := json.MarshalIndent(Message{code, msg, data}, "", ""); err == nil {
		return mess
	} else {
		log.Println("getJson函数：", err)
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

// 全局修饰器
func (G *Global) AuthHandler(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if G.DB.Conn.Driver() == nil {
			if err := G.DB.Connect(); err != nil { // 初始化MYSQL 连接
				log.Println("初始化mysql连接失败: ", err)
				return
			}
		}
		if (r.URL.Path == "/login" || r.URL.Path == "/favicon.ico" || r.Method == "GET") { // 对登录页面不验证cookie
			f(w, r)
			return
		}

		// 判断redis中菜单key是否存在，不存在则创建
		if _, err := G.Redis.Get("menus_"); err != nil {
			if err := G.GetMenuWRedis(); err != nil {
				log.Println("初始化系统菜单错误: ", err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
				return
			}
		}

		CookieAuthToken, CookieAuthTokenErr := r.Cookie("AuthToken")
		CookieUsername, CookieUsernameErr := r.Cookie("Username")
		if CookieAuthTokenErr != nil || CookieUsernameErr != nil { // 获取cookie 失败重定向到登录页面
			fmt.Fprintf(w, string(getJson(-3, "获取Cookie失败,请重新登录", nil)))
		} else {
			if UserInfo, err := G.Redis.Get(CookieUsername.Value); err == nil { // 获取redis 信息
				var JsonRes UserInfoUnmarshal
				json.Unmarshal([]byte(UserInfo), &JsonRes)
				if CookieAuthToken.Value == JsonRes.AuthToken { /// Cookie 正确
					if G.UserInfo.Expiration == 0 { /// 过期时间为0 表示服务器重启过
						if r.Method == "POST" {
							fmt.Fprintf(w, string(getJson(-2, "服务完成重启,请重新登录!", nil)))
							return
						}
					} else { /// 刷新redis过期时间
						G.Redis.ExPire(CookieUsername.Value, G.UserInfo.Expiration)
						f(w, r)
						return
					}
				} else { /// redis中的AuthToken 与 Cookie中的AuthToken 的值不相等 说明账户异地登录
					if r.Method == "POST" {
						fmt.Fprintf(w, string(getJson(-3, "账户异地登录，被迫下线!", nil)))
						return
					}
				}
			} else { /// 未获取到redis 信息 说明cookie 已经过期， 即登录超时
				if r.Method == "POST" {
					fmt.Fprintf(w, string(getJson(-4, "登录超时，请重新登录!", nil)))
					return
				}
			}
		}
	}
}
