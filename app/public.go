package app

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

// 解析菜单json
type MenuToJson []struct {
	Menu_url  string `json:"menu_url"`
	Menu_name string `json:"menu_name"`
}

// 读取菜单信息 写入redis
func (G *Global) GetMenuWRedis() (error) {
	Menus := make([]map[string]string, 0)
	if data, err := G.SelectAll("SELECT * FROM menu"); err == nil {
		for _, va := range *data {
			vmap := make(map[string]string)
			vmap["menu_url"] = va["menu_url"]
			vmap["menu_name"] = va["menu_name"]
			Menus = append(Menus, vmap)
		}
	}
	MenusJson, _ := json.Marshal(Menus)
	if err := G.Redis.Set("menus_", string(MenusJson), 86400); err == nil {
		return nil
	} else {
		return err
	}
}

func (G *Global) Public(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.PostFormValue("action") //获取方法
		Menus := make([]map[string]string, 0)
		if action == "menu" { // 从redis中获取菜单信息
			if Data, err := G.Redis.Get("menus_"); err == nil {
				var MenuJson MenuToJson
				json.Unmarshal([]byte(Data), &MenuJson)
				for _, va := range MenuJson {
					vmap := make(map[string]string)
					vmap["menu_url"] = va.Menu_url
					vmap["menu_name"] = va.Menu_name
					if va.Menu_url == r.PostFormValue("url") {
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
