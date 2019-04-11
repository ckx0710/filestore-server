package handler

import (
	"../meta"
	"../utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
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
		//接收文件流存储到本地目录,FormFile方法参数key对应html页面的name
		surFile, header, e := r.FormFile("File001")
		if e != nil {
			utils.Handerr("r.FormFile", e)
			return
		}
		defer surFile.Close()

		//创建文件元信息
		fileMeta := meta.FileMeta{
			FileName:   header.Filename,
			Path:       "./temp/" + header.Filename,
			UpdateTime: time.Now().Format("2019-4-10 18:16:25"),
		}

		//创建本地文件用来接收文件流
		newFile, e := os.Create(fileMeta.Path)
		if e != nil {
			utils.Handerr("os.Create", e)
			return
		}
		defer newFile.Close()

		//用io.copy方法完成文件流对本地文件的下载
		fileMeta.FileSize, e = io.Copy(newFile, surFile)
		if e != nil {
			utils.Handerr("io.Copy", e)
			return
		}

		//这里还不明白
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		meta.UpdateFileMeta(&fileMeta)
		fmt.Println("上传成功！", fileMeta.FileName, ":", fileMeta.FileSha1)

		//重定向
		http.Redirect(w, r, "/file/uploadSuces", http.StatusFound)

	}

}

//完成文件上传，反馈消息
func SuceMsage(wrt http.ResponseWriter, rst *http.Request) {
	io.WriteString(wrt, "upload file success!")
}

func QueryMetaBySha1(w http.ResponseWriter, r *http.Request) {
	//解析url请求参数，并将解析结果更新到r.Form字段。
	err := r.ParseForm()
	if err != nil {
		utils.Handerr("r.ParseForm()", err)
		return
	}

	//获取请求表单力的fileHash参数值
	fileHash := r.Form["sha1"][0]
	fileMeta := *meta.GetFileMeta(fileHash)

	//页面输出打印查询到的元文件json信息
	bytes, err := json.Marshal(fileMeta)
	if err != nil {
		utils.Handerr("json.Marshal", err)
		return
	}
	w.Write(bytes)
}

//查询以时间为参照最新上传的N个文件
func QueryMetasBySha1(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limit, _ := strconv.Atoi(r.Form["limit"][0])
	fileMetaSplic := meta.GetFileMetaSplic()

	//把获取到的fileMetaSplic按时间重新排序
	sort.Sort(meta.ByUpdateTime(fileMetaSplic))
	len := len(fileMetaSplic)
	//防止用户输入参数大于切片自身长度，返回最后上传文件
	if limit <= len {
		len -= limit
	} else {
		len--
	}
	fileMetas := fileMetaSplic[len:]

	//页面输出打印查询到的元文件json信息
	bytes, err := json.Marshal(fileMetas)
	if err != nil {
		utils.Handerr("json.Marshal", err)
		return
	}
	w.Write(bytes)
}

//文件下载
func DownloadFileMeta(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sha1 := r.Form.Get("sha1")
	fileMeta := meta.GetFileMeta(sha1)

	file, e := os.Open(fileMeta.Path)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	bytes, e := ioutil.ReadAll(file)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//★★★★★启动浏览器下载
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fileMeta.FileName+"\"")

	w.Write(bytes)
}

//更新元信息
func UpdateFileMeta(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opTepy := r.Form.Get("op")
	sha1 := r.Form.Get("sha1")
	newFileName := r.Form.Get("name")
	fmt.Println("r.Method：",r.Method)
	if opTepy != "0" {
		w.WriteHeader(http.StatusForbidden)//403 :参数错误
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)//405 请求方式错误
		return
	}

	fileMeta := meta.GetFileMeta(sha1)
	fileMeta.FileName=newFileName
	//更新元文件信息到集合Map
	meta.UpdateFileMeta(fileMeta)

	bytes, e := json.Marshal(fileMeta)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}
	w.Write(bytes)
	w.WriteHeader(http.StatusOK)//返回 200
}

//文件删除：物理删除+逻辑删除
func DeleteFileMeta(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	sha1 := r.Form.Get("sha1")
	fileMeta := meta.GetFileMeta(sha1)
	//物理删除
	os.Remove(fileMeta.Path)

	//逻辑删除
	meta.DeleteFileMetaMap(fileMeta)
	w.WriteHeader(http.StatusOK)

}
