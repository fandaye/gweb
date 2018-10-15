package app

import (
	"net/http"
	"fmt"
	"log"
)

func (G *Global) Public(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.PostFormValue("action") //获取方法
		Menus := make([]map[string]string, 0)
		if action == "menu" {
			if data, err := G.SelectAll("SELECT * FROM menu"); err == nil {
				for _, va := range *data {
					vmap := make(map[string]string)
					vmap["menu_url"] = va["menu_url"]
					vmap["menu_name"] = va["menu_name"]
					if va["menu_url"] == r.PostFormValue("url") {
						vmap["active"] = "true"
					} else {
						vmap["active"] = "false"
					}
					Menus = append(Menus, vmap)
				}
				fmt.Fprintf(w, string(getJson(0, "数据请求成功!", Menus)))
			} else {
				log.Println(err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", Menus)))
			}
		}
	}
}
