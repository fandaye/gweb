package app

import (
	"net/http"
	"html/template"
	"log"
	"fmt"
	"encoding/json"
	"time"
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
					vmap["id"] = va["id"]
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
				if _, err := G.DB.Insert("INSERT INTO menu (menu_url, menu_name, menu_lcon, create_time, instruction) VALUES (?, ?, 'lnr-tag', ?, ?)", m.Menu_url, m.Menu_name, time.Now().Format("2006-01-02 15:04:05"), m.Instruction); err == nil {
					fmt.Fprintf(w, string(getJson(0, "添加成功!", nil)))
				} else {
					log.Println("添加系统菜单失败:  ", err.Error())
					fmt.Fprintf(w, string(getJson(-1, "添加失败!", nil)))
				}
			} else {
				log.Println("添加系统菜单失败:  ", err.Error())
				fmt.Fprintf(w, string(getJson(-1, "添加失败!", nil)))
			}
		} else if action == "DelMenu" {
			if _, err := G.DB.Delete("DELETE FROM menu WHERE id = ?", r.PostFormValue("id")); err == nil {
				fmt.Fprintf(w, string(getJson(0, "删除成功!", nil)))
			} else {
				log.Println("删除系统菜单失败:  ", err.Error())
				fmt.Fprintf(w, string(getJson(-1, "删除失败!", nil)))
			}

		}else if action == "EditMenu" {
			var m SysMenuJson
			json.Unmarshal([]byte(r.PostFormValue("data")), &m)
			fmt.Println(m)

			fmt.Fprintf(w, string(getJson(0, "编辑成功!", nil)))
		}
	}
}
