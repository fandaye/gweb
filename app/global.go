package app

import (
	"encoding/json"
)

type Global struct {
	Username string
}

func getJson(code int, msg string, data map[string]string) ([]byte, error) {
	type Message struct {
		Code int
		Msg  string
		Data map[string]string
	}
	return json.MarshalIndent(Message{code, msg, data}, "", "")
}
