package handler

import (
	"../utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//获取上传主页
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			utils.Handerr("ioutil.ReadFile", err)
			return
		}

		io.WriteString(w, string(data))

	case "POST":
		//接收文件流存储到本地目录
		surFile, header, e := r.FormFile("surFile")
		if e != nil {
			utils.Handerr("r.FormFile",e)
			return
		}
		defer surFile.Close()

		//创建本地文件用来接收文件流
		newFile, e := os.Create("/temp/" + header.Filename)
		if e != nil {
			utils.Handerr("os.Create",e)
			return
		}
		defer  newFile.Close()

		//用io.copy方法完成文件流对本地文件的下载
		_, e = io.Copy(newFile, surFile)
		if e != nil {
			utils.Handerr("io.Copy",e)
			return 
		}

		//重定向
		http.Redirect(w,r,"/file/uploadSuces",http.StatusFound)

	}

}

func SuceMsage(wrt http.ResponseWriter,rst *http.Request)  {
	io.WriteString(wrt,"upload file success!")
}
