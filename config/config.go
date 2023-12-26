package config

import (
	"github.com/go-ini/ini"
	"live_recording_tools/app/utils"
	"runtime"
)

type SysConfig struct {
	Env      string `ini:"env"`
	Debug    bool   `ini:"debug"`
	DBDriver string `ini:"db_driver"`
	DBHost   string `ini:"db_host"`
	DBPort   string `ini:"db_port"`
	DBUser   string `ini:"db_user"`
	DBPass   string `ini:"db_pass"`
	DBName   string `ini:"db_name"`
	DBDebug  bool   `ini:"db_debug"`

	HttpListenPort string `ini:"http_listen_port"`

	RedisHost         string `ini:"redis_host"`
	RedisPwd          string `ini:"redis_pwd"`
	RedisDb           int    `ini:"redis_db"`
	RedisCacheVersion string `ini:"redis_cache_version"`

	IpDbPath string
}

var Configs *SysConfig = &SysConfig{}

// 加载系统配置文件
func Default() {
	appDir := utils.GetAppDir()
	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "windows":
		appDir = utils.GetAppDir()
		break
	default:
		appDir = utils.GetCurrentPath()
	}

	conf, err := ini.Load(appDir + "/config.ini") //加载配置文件
	if err != nil {
		panic(err)
	}
	conf.BlockMode = false
	err = conf.MapTo(&Configs) //解析成结构体
	if err != nil {
		panic(err)
	}
}
