package app

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
	"gopkg.in/ldap.v2"
	"strconv"
)

type LdapUserInfo struct {
	Code int
	Msg   ldap.Conn
}

func (G *Global) Ldap(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if t, err := template.ParseFiles("templates/ldap/index.html", "templates/common/head.html", "templates/common/tail.html") ; err == nil{
			t.Execute(w, G)
		} else {
			log.Println( "函数 Ldap: " , err)
			G.Code= "-1"
		}
	}else if r.Method == "POST"{
		action := r.PostFormValue("action") //获取方法
		G.Ldaps.Config = G.Config.LdapConfig() // 初始化ldap配置
		if ldap_port, err := strconv.Atoi(G.Ldaps.Config["ldap_port"]); err == nil { //端口号 字符串 转 数字
			ldap_concatenon, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", G.Ldaps.Config["ldap_host"], ldap_port)) //初始化ldap连接
			if err != nil {
				log.Println( "函数 Ldap 初始化连接错误 : " , err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!" ,nil)))
				return
			} else {
				if Err := ldap_concatenon.Bind(G.Ldaps.Config["ldap_user"], G.Ldaps.Config["ldap_passwd"]); Err != nil {
					log.Println( "函数 Ldap 验证用户密码错误 : " , err)
					return
				} else {
					G.Ldaps.Conn = *ldap_concatenon
					defer ldap_concatenon.Close()
				}
			}
		} else {
			log.Println( "函数 Ldap 转换ldap端口号错误 : " , err)
			return
		}

		if action == "list" {
			LdapGroupAll := make(map[string]string)
			GroupObjectClass := "(&(objectClass=top))"
			GroupAttributes := []string{"description", "gidNumber"}
			if Info, Err := G.Ldaps.Search("ou=Group,"+G.Ldaps.Config["ldap_base_dn"], GroupObjectClass, GroupAttributes); Err == nil {
				for _, entry := range Info {
					LdapGroupAll[entry.GetAttributeValue("gidNumber")] = entry.GetAttributeValue("description")
				}
			}else {
				log.Println( "函数 Ldap 查询用户组信息 : " , Err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!" ,nil)))
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
			}else {
				log.Println( "函数 Ldap 查询用户信息 : " , Err)
				fmt.Fprintf(w, string(getJson(-1, "系统错误!" ,nil)))
			}
		}
	}
}
