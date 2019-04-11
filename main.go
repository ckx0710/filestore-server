package main

import (
	"./handler"
	"./utils"
	"fmt"
	"net/http"
)

func main() {
	//添加设置路由规则
	http.HandleFunc("/file/upload/", handler.UploadHandler)
	http.HandleFunc("/file/uploadSuces",handler.SuceMsage)
	http.HandleFunc("/file/queryMeta",handler.QueryMetaBySha1)
	http.HandleFunc("/file/queryMetas",handler.QueryMetasBySha1)
	http.HandleFunc("/file/down",handler.DownloadFileMeta)
	http.HandleFunc("/file/update",handler.UpdateFileMeta)
	http.HandleFunc("/file/delete",handler.DeleteFileMeta)

	//初始化开启服务
	fmt.Println("***Server is star =>")
	err := http.ListenAndServe(":8080", nil)
	if err!=nil{
		utils.Handerr("http.ListenAndServe：", err)
		return
	}

}
