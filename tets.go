package main

import "fmt"

func main() {

	yinzhengjie := make(map[string]string)
	yinzhengjie["1"] = "运维"
	yinzhengjie["2"] = "商务"

	fmt.Println(yinzhengjie)

	fmt.Println(yinzhengjie["1"])

}
