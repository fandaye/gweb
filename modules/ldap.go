package modules

import (
	"gopkg.in/ldap.v2"
	"strconv"
	"unicode/utf16"
	"gitee.com/zabbix_669/md4"
	"strings"
	"encoding/hex"
	"encoding/binary"
)

type Ldaps struct {
	Config map[string]string
	Conn   ldap.Conn
}

//查询函数
func (L *Ldaps) Search(search_dn, objectClass string, attributes []string) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		search_dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		objectClass,
		attributes,
		nil,
	)
	SearchInfo, Err := L.Conn.Search(searchRequest)
	return SearchInfo.Entries, Err
}

//添加用户
func (L *Ldaps) LdapUserADD(givenName, sn, uid, gidNumber, mobile, mail, UnixTime, password string) (error) {
	DN := "uid=" + uid + "," + L.Config["ldap_user_dn"]
	uidNumber := L.GetRandomuidNumber()
	//user_passwd := L.GetRandomPasswd(10) //生成随机密码
	add := ldap.NewAddRequest(DN)
	add.Attribute("objectClass", []string{"inetOrgPerson", "posixAccount", "top", "shadowAccount", "ldapPublicKey", "sambaSamAccount", "hostObject"})
	add.Attribute("shadowLastChange", []string{UnixTime})
	add.Attribute("sambaPwdLastSet", []string{UnixTime})
	add.Attribute("sambaNTPassword", []string{L.GetSambantPassword(password)})
	add.Attribute("sambaAcctFlags", []string{"[U]"})
	add.Attribute("userPassword", []string{""})
	add.Attribute("sshPublicKey", []string{""})
	add.Attribute("displayName", []string{givenName + sn})
	add.Attribute("givenName", []string{givenName})
	add.Attribute("gidNumber", []string{gidNumber})
	add.Attribute("uidNumber", []string{uidNumber})
	add.Attribute("sambaSID", []string{L.Config["ldap_sambasid"] + uidNumber})
	add.Attribute("mobile", []string{mobile})
	add.Attribute("mail", []string{mail})
	add.Attribute("uid", []string{uid})
	add.Attribute("cn", []string{givenName + sn})
	add.Attribute("sn", []string{sn})
	add.Attribute("homeDirectory", []string{"/home/" + uid})
	add.Attribute("loginShell", []string{"/bin/bash"})

	Err := L.Conn.Add(add) //执行添加用户操作
	if Err != nil {
		return Err
	} else { //添加成功操作
		L.Conn.PasswordModify(ldap.NewPasswordModifyRequest(DN, "", password))
		return nil
	}
	return nil
}

//添加用户
func (L *Ldaps) LdapGroupADD(cn, description string) (error) {
	DN := "cn=" + cn + "," + L.Config["ldap_group_dn"]
	add := ldap.NewAddRequest(DN)
	//添加对象
	add.Attribute("objectClass", []string{"posixGroup", "top"})
	add.Attribute("gidNumber", []string{L.GetRandomgidNumber()})
	add.Attribute("description", []string{description})
	add.Attribute("cn", []string{cn})
	Err := L.Conn.Add(add) //执行添加用户操作
	return Err
}

//删除DN操作
func (L *Ldaps) Del(DN string) error {
	if Err := L.Conn.Del(ldap.NewDelRequest(DN, nil)); Err != nil {
		return Err
	} else {
		return nil
	}
}

///重置密码
func (L *Ldaps) ResePasswd(uid, passwd, UnixTime string) (error) {

	DN := "uid=" + uid + "," + L.Config["ldap_user_dn"]
	L.ModifAttrtype(DN, "sambaNTPassword", L.GetSambantPassword(passwd))   //修改sambaNTPassword
	L.ModifAttrtype(DN, "sambaPwdLastSet", UnixTime)                       //修改时间
	L.ModifAttrtype(DN, "shadowLastChange", UnixTime)                      //修改时间
	passwordModifyRequest := ldap.NewPasswordModifyRequest(DN, "", passwd) //修改userPassword
	_, Err := L.Conn.PasswordModify(passwordModifyRequest)
	return Err
}

//修改指定DN中某个Attrtype值
func (L *Ldaps) ModifAttrtype(DN, Attrtype, Value string) error {
	Modif := ldap.NewModifyRequest(DN)
	Modif.Replace(Attrtype, []string{Value})
	if Err := L.Conn.Modify(Modif); Err != nil {
		return Err
	}
	return nil
}

//添加指定DN中某个Attrtype值
func (L *Ldaps) AddAttrtype(DN, Attrtype, Value string) error {
	Modify_ADD := ldap.NewModifyRequest(DN)
	Modify_ADD.Add(Attrtype, []string{Value})
	if Err := L.Conn.Modify(Modify_ADD); Err != nil {
		return Err
	}
	return nil
}

//删除DN某个Attrtype值
func (L *Ldaps) DelAttrtype(DN, Attrtype, Value string) error {
	Modify_DEL := ldap.NewModifyRequest(DN)
	Modify_DEL.Delete(Attrtype, []string{Value})
	if Err := L.Conn.Modify(Modify_DEL); Err != nil {
		return Err
	}
	return nil
}

//生成用户uidnumber
func (L *Ldaps) GetRandomuidNumber() string {
	SearchDn := L.Config["ldap_user_dn"]
	ObjectClass := "(&(objectClass=shadowAccount))"
	Attributes := []string{"uidNumber"}
	uidNumber := 10000
	if LdapUserAll, Err := L.Search(SearchDn, ObjectClass, Attributes); Err == nil {
		for _, entry := range LdapUserAll {
			b, _ := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
			if b > uidNumber {
				uidNumber = b
			}
		}
	}
	return strconv.Itoa(uidNumber + 1)
}

//生成分组gidNumber
func (L *Ldaps) GetRandomgidNumber() string {
	SearchDn := L.Config["ldap_group_dn"]
	ObjectClass := "(&(objectClass=posixGroup)(objectClass=top))"
	Attributes := []string{"gidNumber"}
	uidNumber := 10000
	if LdapGroupAll, Err := L.Search(SearchDn, ObjectClass, Attributes); Err == nil {
		for _, entry := range LdapGroupAll {
			b, _ := strconv.Atoi(entry.GetAttributeValue("gidNumber"))
			if b > uidNumber {
				uidNumber = b
			}
		}
	}
	return strconv.Itoa(uidNumber + 1)
}

//生成samb密码
func (L *Ldaps) GetSambantPassword(passwd string) string {
	passwd_utf16 := utf16.Encode([]rune(passwd))
	b := make([]byte, 2*len(passwd_utf16))
	for index, value := range passwd_utf16 {
		binary.LittleEndian.PutUint16(b[index*2:], value)
	}
	h_m4 := md4.New()
	h_m4.Write(b)
	return strings.ToUpper(hex.EncodeToString(h_m4.Sum(nil)))
}
