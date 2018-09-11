package modules

type UserInfo struct {
	Username string
	Email  string
	Passwd    string
	AuthToken string
	Role      string
	Prems     []map[string]string
	Expiration int  // 过期时间
	LoginErrorNumber int // 登录失败次数
}
