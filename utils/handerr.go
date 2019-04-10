package utils

import "fmt"

func Handerr(where string, err error) {
	if err==nil{
		return
	}
	errStr := fmt.Sprintf("%s", err)
	fmt.Println(where, ":", errStr)
}
