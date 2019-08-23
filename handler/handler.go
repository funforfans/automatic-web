package handler

import (
	"encoding/json"
	"github.com/micro/go-micro/util/log"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Error 错误结构体
type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

const (
	UPLOAD_DIR = "./uploads"
)

func Trigger(w http.ResponseWriter, r *http.Request) {
	log.Log("Trigger")
	// 只接受POST请求
	if r.Method != "POST" {
		log.Logf("非法请求")
		http.Error(w, "非法请求", 400)
		return
	}
	// 解析参数
	if err:=r.ParseForm();err!=nil{
		log.Logf("参数解析异常")
		http.Error(w, "参数解析异常", 400)
		return
	}
	action := r.Form.Get("action")
	log.Logf("action is %s", action)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// 构造返回结果
	response := map[string]interface{}{
		"time": time.Now().UnixNano(),
	}
	response["action"]= action+" done!"
	// 返回JSON结构
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// 上传文件功能
//
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	fileId := r.FormValue("id")
	filePath := UPLOAD_DIR + "/" + fileId
	if exists := isExists(filePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "file")
	http.ServeFile(w, r, filePath)
}
// 是否存在该文件
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	var configDir interface{}
	value, ok := baseConfig["uploader_dirname"]
	if ok{
		configDir=value
	} else {
		log.Fatal("config uploader_dirname not existed")
	}
	uploadDir := configDir.(string)
	log.Log(uploadDir)
	// 先拉取页面
	if r.Method == "GET" {
		folderName := "./" + uploadDir + "/upload.html"
		t, err := template.ParseFiles(folderName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		defer f.Close()
		// 创建上传文件夹目录结构
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		// 跳转展示页面，触发ViewHandler
		http.Redirect(w, r, "/view?id="+filename,
			http.StatusFound)
	}
}

func WebhookHandler(w http.ResponseWriter, r *http.Request)  {

	defer r.Body.Close()

	log.Log("WebhookHandler")
	// 返回json格式的header
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// webhook 默认只接受POST请求
	if r.Method != "POST" {
		log.Logf("非POST的请求，sorry！")
		http.Error(w, "非POST的请求，sorry！", http.StatusBadRequest)
		return
	}
	// 解析参数
	if err:=r.ParseForm();err!=nil{
		log.Logf("参数解析异常")
		http.Error(w, "参数解析异常", http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var webData interface{}
	if err := json.Unmarshal(b, &webData); err!=nil{
		log.Logf("json解析异常")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 后面增加一些业务
	// TODO
	mapdata := webData.(map[string]interface{})
	for k, v := range mapdata{
		log.Logf("key:%v, value:%v\n", k,v)
	}
	// 构造返回结果
	response := map[string]interface{}{
		"time": time.Now().UnixNano(),
	}
	// 返回JSON结构
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}