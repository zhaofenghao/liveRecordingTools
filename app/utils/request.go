package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 请求结构体
type Request struct {
	Method  string            `json:"method"`  // 请求方法
	Url     string            `json:"url"`     // 请求url
	Params  map[string]string `json:"params"`  // Query参数
	Headers map[string]string `json:"headers"` // 请求头
	Cookies map[string]string `json:"cookies"` // todo 处理 Cookies
	Data    map[string]string `json:"data"`    // 表单格式请求数据
	Json    map[string]string `json:"json"`    // JSON格式请求数据 todo 多层 嵌套
	Files   map[string]string `json:"files"`   // todo 处理 Files
	Raw     string            `json:"raw"`     // 原始请求数据
}

// 响应结构体
type Response struct {
	StatusCode int               `json:"status_code"` // 状态码
	Reason     string            `json:"reason"`      // 状态码说明
	Elapsed    float64           `json:"elapsed"`     // 请求耗时(秒)
	Content    []byte            `json:"content"`     // 响应二进制内容
	Text       string            `json:"text"`        // 响应文本
	Headers    map[string]string `json:"headers"`     // 响应头
	Cookies    map[string]string `json:"cookies"`     // todo 添加响应Cookies
	Request    *Request          `json:"request"`     // 原始请求
}

// 处理请求方法
func (r *Request) getMethod() string {
	return strings.ToUpper(r.Method) // 必须转为全部大写
}

// 组装URL
func (r *Request) getUrl() string {
	if r.Params != nil {
		urlValues := url.Values{}
		Url, _ := url.Parse(r.Url) // todo 处理err
		for key, value := range r.Params {
			urlValues.Set(key, value)
		}
		Url.RawQuery = urlValues.Encode()
		return Url.String()
	}
	return r.Url
}

// 组装请求数据
func (r *Request) getData() io.Reader {
	var reqBody string
	if r.Raw != "" {
		reqBody = r.Raw
	} else if r.Data != nil {
		urlValues := url.Values{}
		for key, value := range r.Data {
			urlValues.Add(key, value)
		}
		reqBody = urlValues.Encode()
		r.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else if r.Json != nil {
		bytesData, _ := json.Marshal(r.Json)
		reqBody = string(bytesData)
		r.Headers["Content-Type"] = "application/json"
	}
	return strings.NewReader(reqBody)
}

// 添加请求头-需要在getData后使用
func (r *Request) addHeaders(req *http.Request) {
	if r.Headers != nil {
		for key, value := range r.Headers {
			req.Header.Add(key, value)
		}
	}
}

// 准备请求
func (r *Request) prepare() *http.Request {
	Method := r.getMethod()
	Url := r.getUrl()
	Data := r.getData()
	req, _ := http.NewRequest(Method, Url, Data)
	r.addHeaders(req)
	return req
}

// 组装响应对象
func (r *Request) packResponse(res *http.Response, elapsed float64) Response {
	var resp Response
	resBody, _ := ioutil.ReadAll(res.Body)
	resp.Content = resBody
	resp.Text = string(resBody)
	resp.StatusCode = res.StatusCode
	resp.Reason = strings.Split(res.Status, " ")[1]
	resp.Elapsed = elapsed
	resp.Headers = map[string]string{}
	for key, value := range res.Header {
		resp.Headers[key] = strings.Join(value, ";")
	}
	return resp
}

// 发送请求
func (r *Request) Send() Response {
	req := r.prepare()
	client := &http.Client{}
	start := time.Now()
	res, _ := client.Do(req) // todo 处理err
	defer res.Body.Close()
	elapsed := time.Since(start).Seconds()
	return r.packResponse(res, elapsed)
}
