package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

func GetAppDir() string {
	appDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		applicationPath, _ := filepath.Abs(file)
		appDir, _ = filepath.Split(applicationPath)
	}
	return appDir
}

func GetCurrentPath() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return exPath
}

func HttpGet(url string, head bool) (body []byte, err error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	var rqt *http.Request
	rqt, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if head {
	}
	var response *http.Response
	response, err = client.Do(rqt)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	return
}

// 发送POST请求(带请求头)
// url:请求地址，data:POST请求提交的数据,
// contentType:
// (1)application/x-www-form-urlencoded  最常见的POST提交数据的方式，浏览器的原生form表单。后面可以跟charset=utf-8
// (2)multipart/form-data
// (3)application/json
// (4)text/xml    XML-RPC远程调用
// content:请求放回的内容
func HttpPostHeader(url string, data interface{}, contentType string, header map[string]interface{}) (content string, err error) {
	jsonStr, _ := json.Marshal(data)
	fmt.Println("jsonStr:", string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("content-type", contentType)
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v.(string))
		}
	}
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

func SendDingMessage(msg string, sn string) {
	var body = map[string]string{
		"sn":  sn,
		"msg": msg,
	}
	var body_byte, _ = json.Marshal(body)

	http.Post("http://api.qingjuhe.com/robot", "application/json", bytes.NewReader(body_byte))
}

func In(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func StructToMap(obj interface{}) map[string]string {
	m := make(map[string]string)
	j, _ := json.Marshal(obj)
	json.Unmarshal(j, &m)
	return m
}
