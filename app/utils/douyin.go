package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 初始化
func InitConfig(path string) io.Writer {
	CheckFile(path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	// 写入默认配置
	var tmp = make([]byte, 128)
	n, err := file.Read(tmp)
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	}
	if n == 0 {
		var content = "{\"save_path\":\"\",\"check_frequency\": 30,\"config\": [{\"is_monitor\": false,\"status\": 1,\"url\": \"\", \"sort\": 0}]}"
		_, err := file.Write([]byte(content))
		if err != nil {
			log.Fatalln("Init config json failed ======>", err)
		}
	}
	return io.Writer(file)
}

func CheckFile(filePath string) {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		path := filepath.Dir(filePath)
		MkDir(path)
	case os.IsPermission(err):
		log.Fatalf("permission:%v", err)
	}
}

func MkDir(filePath string) {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func FileGetContents(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	f.Close()
	return data, err
}

type DouyinConfig struct {
	SavePath       string             `json:"save_path"`
	CheckFrequency int                `json:"check_frequency"`
	Config         []DouyunRoomConfig `json:"config"`
}

type DouyunRoomConfig struct {
	IsMonitor bool   `json:"is_monitor"`
	Status    int    `json:"status"`
	Url       string `json:"url"`
	Sort      int    `json:"sort"`
}

func SaveConfigs(config DouyinConfig) error {
	f, err := os.OpenFile(config.SavePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	} else {
		data, err := json.Marshal(&config)
		if err == nil {
			_, err := f.Write(data)
			if err == nil {
				err = f.Close()
				return err
			} else {
				return err
			}
		} else {
			return err
		}
	}
}

// 读取配置文件
func ReadConfigList(path string) (config DouyinConfig, err error) {
	var content []byte
	_ = InitConfig(path)
	content, err = FileGetContents(path)
	if err != nil {
		fmt.Println("open file error: " + err.Error())
		return
	}
	config = DouyinConfig{}
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	return
}
