package main

import (
	"./handler"
	"./utils"
	"fmt"
	"net/http"
)

func main() {
	//路由设置
	http.HandleFunc("/file/upload/", handler.UploadHandler)
	http.HandleFunc("/file/uploadSuces",handler.SuceMsage)

	//初始化开启服务
	err := http.ListenAndServe(":8080", nil)
	if err!=nil{
		utils.Handerr("http.ListenAndServe：", err)
		return
	}
	fmt.Println("***Server is star =>")

}
