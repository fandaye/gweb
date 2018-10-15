package app

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
	"gopkg.in/ldap.v2"
	"strconv"
	"encoding/json"
)

type LdapUserInfo struct {
	Code int
	Msg  ldap.Conn
}

type LdapGroupJson struct {
	Cn          string `json:"cn"`
	Description string `json:"description"`
}

func (G *Global) Ldap(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/ldap/index.html", "templates/common/head.html", "templates/common/tail.html"); err == nil {
			t.Execute(w, G)
		} else {
			log.Println("函数 Ldap: ", err)
			G.Code = "-1"
		}
	} else if r.Method == "POST" {
		action := r.PostFormValue("action")    //获取方法
		G.Ldaps.Config = G.Config.LdapConfig() // 初始化ldap配置
		if ldap_port, err := strconv.Atoi(G.Ldaps.Config["ldap_port"]); err == nil { //端口号 字符串 转 数字
			ldap_concatenon, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", G.Ldaps.Config["ldap_host"], ldap_port)) //初始化ldap连接
			if err != nil {
				log.Println("函数 Ldap 初始化连接错误 : ", err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
				return
			} else {
				if Err := ldap_concatenon.Bind(G.Ldaps.Config["ldap_user"], G.Ldaps.Config["ldap_passwd"]); Err != nil {
					log.Println("函数 Ldap 验证用户密码错误 : ", err)
					return
				} else {
					G.Ldaps.Conn = *ldap_concatenon
					defer ldap_concatenon.Close()
				}
			}
		} else {
			log.Println("函数 Ldap 转换ldap端口号错误 : ", err)
			return
		}
		if action == "LdapUserInfo" {
			LdapGroupAll := make(map[string]string)
			GroupObjectClass := "(&(objectClass=top))"
			GroupAttributes := []string{"description", "gidNumber"}
			if Info, Err := G.Ldaps.Search("ou=Group,"+G.Ldaps.Config["ldap_base_dn"], GroupObjectClass, GroupAttributes); Err == nil {
				for _, entry := range Info {
					LdapGroupAll[entry.GetAttributeValue("gidNumber")] = entry.GetAttributeValue("description")
				}
			} else {
				log.Println("函数 Ldap 查询用户组信息 : ", Err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
			}
			LdapUserAll := make([]map[string]string, 0)
			PeopleObjectClass := "(&(objectClass=shadowAccount))"
			PeopleAttributes := []string{"displayName", "mobile", "uid", "mail", "gidNumber"}
			if Info, Err := G.Ldaps.Search("ou=People,"+G.Ldaps.Config["ldap_base_dn"], PeopleObjectClass, PeopleAttributes); Err == nil {
				for _, entry := range Info {
					vmap := make(map[string]string)
					vmap["uid"] = entry.GetAttributeValue("uid")
					vmap["displayName"] = entry.GetAttributeValue("displayName")
					vmap["mail"] = entry.GetAttributeValue("mail")
					vmap["mobile"] = entry.GetAttributeValue("mobile")
					vmap["group"] = LdapGroupAll[entry.GetAttributeValue("gidNumber")]
					LdapUserAll = append(LdapUserAll, vmap)
				}
				fmt.Fprintf(w, string(getJson(0, "数据请求成功!", LdapUserAll)))
			} else {
				log.Println("函数 Ldap 查询用户信息 : ", Err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
			}
		} else if action == "LdapGroupInfo" {
			LdapGroupAll := make([]map[string]string, 0)
			GroupObjectClass := "(&(objectClass=top))"
			GroupAttributes := []string{"description", "gidNumber", "cn"}
			if Info, Err := G.Ldaps.Search("ou=Group,"+G.Ldaps.Config["ldap_base_dn"], GroupObjectClass, GroupAttributes); Err == nil {
				for _, entry := range Info {
					vmap := make(map[string]string)
					if entry.GetAttributeValue("cn") != "" {
						vmap["cn"] = entry.GetAttributeValue("cn")
						vmap["gidNumber"] = entry.GetAttributeValue("gidNumber")
						vmap["description"] = entry.GetAttributeValue("description")
						LdapGroupAll = append(LdapGroupAll, vmap)
					}
				}
				fmt.Fprintf(w, string(getJson(0, "数据请求成功!", LdapGroupAll)))
			} else {
				log.Println("函数 Ldap 查询用户组信息 : ", Err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
			}
		}else if action == "AddGroup" {
			var m LdapGroupJson
			if err:= json.Unmarshal([]byte(r.PostFormValue("data")), &m) ; err == nil{
				add := ldap.NewAddRequest("cn=" + m.Cn + "," + "ou=Group,"+G.Ldaps.Config["ldap_base_dn"])
				//添加对象
				add.Attribute("objectClass", []string{"posixGroup", "top"})
				add.Attribute("gidNumber", []string{G.GetRandomgidNumber()})
				add.Attribute("description", []string{m.Description})
				add.Attribute("cn", []string{m.Cn})
				if Err := G.Ldaps.Conn.Add(add) ;Err ==nil{
					fmt.Fprintf(w, string(getJson(1, "添加成功!", nil)))
				}else {
					log.Println("函数Ldap AddGroup错误: " ,Err)
					fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
				}
			}else {
				log.Println("函数Ldap 解析Json失败 :" , err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
			}
		}else if action == "DleGroup" {
			if Err := G.Ldaps.Del("cn=" + r.PostFormValue("cn") + "," + "ou=Group,"+G.Ldaps.Config["ldap_base_dn"]); Err == nil {
				fmt.Fprintf(w, string(getJson(1, "删除成功!", nil)))
			} else {
				log.Println(Err.Error())
				fmt.Fprintf(w, string(getJson(-1, "系统错误!", nil)))
			}
		}
	}
}
