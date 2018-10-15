package app

import (
	"net/http"
	"html/template"
	"log"
	"fmt"
	"encoding/json"
)

type SysMenuJson struct {
	Menu_url    string `json:"menu_url"`
	Menu_name   string `json:"menu_name"`
	Instruction string `json:"instruction"`
}

func (G *Global) Sys(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/sys/index.html", "templates/common/head.html", "templates/common/tail.html"); err == nil {
			t.Execute(w, G)
		} else {
			log.Println("函数 Sys: ", err)
			fmt.Fprintf(w, err.Error())
		}
	} else if r.Method == "POST" {
		action := r.PostFormValue("action") //获取方法
		Data := make([]map[string]string, 0)
		if action == "SysUserInfo" {
			if data, err := G.SelectAll("SELECT * FROM users"); err == nil {
				for _, va := range *data {
					vmap := make(map[string]string)
					vmap["email"] = va["email"]
					vmap["username"] = va["username"]
					vmap["create_time"] = va["create_time"]
					vmap["login_time"] = va["login_time"]
					vmap["status"] = va["status"]
					vmap["role_id"] = va["role_id"]
					Data = append(Data, vmap)
				}
				fmt.Fprintf(w, string(getJson(0, "数据请求成功!", Data)))
			} else {
				log.Println(err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", Data)))
			}

		} else if action == "SysMenuInfo" {
			if data, err := G.SelectAll("SELECT * FROM menu"); err == nil {
				for _, va := range *data {
					vmap := make(map[string]string)
					vmap["menu_lcon"] = va["menu_lcon"]
					vmap["menu_url"] = va["menu_url"]
					vmap["menu_name"] = va["menu_name"]
					vmap["create_time"] = va["create_time"]
					vmap["instruction"] = va["instruction"]
					Data = append(Data, vmap)
				}
				fmt.Fprintf(w, string(getJson(0, "数据请求成功!", Data)))
			} else {
				log.Println(err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", Data)))
			}
		} else if action == "SysRoleInfo" {
			if data, err := G.SelectAll("SELECT * FROM roles"); err == nil {
				for _, va := range *data {
					vmap := make(map[string]string)
					vmap["role_name"] = va["role_name"]
					vmap["create_time"] = va["create_time"]
					vmap["instruction"] = va["instruction"]
					Data = append(Data, vmap)
				}
				fmt.Fprintf(w, string(getJson(0, "数据请求成功!", Data)))
			} else {
				log.Println(err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", Data)))
			}
		} else if action == "AddMenu" {
			var m SysMenuJson
			if err := json.Unmarshal([]byte(r.PostFormValue("data")), &m); err == nil {
				fmt.Println(m.Instruction, m.Menu_name, m.Menu_url)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))

			} else {
				log.Println("函数Sys AddMenu 解析Json失败 :", err)
				fmt.Fprintf(w, string(getJson(1, "添加成功!", nil)))
			}

		}
	}
}
