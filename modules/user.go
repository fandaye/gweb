package modules

type UserInfo struct {
	Username string
	Email  string
	Passwd    string
	AuthToken string
	Role      string
	Prems     []map[string]string
}
