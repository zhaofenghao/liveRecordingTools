package daemon

import (
	Middlewares "ffmpeg_work/app/middleware"
	"ffmpeg_work/app/router"
	"ffmpeg_work/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	_ "github.com/mkevac/debugcharts"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const version = "1.0.0"

func RunService() {
	//服务信息
	options := make(service.KeyValue)
	options["LimitNOFILE"] = 1000000
	svcConfig := &service.Config{
		Name:        "proxyChecker",
		DisplayName: "proxyChecker",
		Description: "proxyChecker",
		Option:      options,
	}
	prg := &BaseService{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Printf("%s 启动失败: %s", svcConfig.DisplayName, err)
		return
	}
	//监听指令
	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			log.Printf("【ServerInstalled】 Service %s install.\n", svcConfig.DisplayName)
		case "uninstall":
			err = s.Uninstall()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			log.Printf("【ServerUninstall】 Service %s uninstall.\n", svcConfig.DisplayName)
		case "run":
			_ = s.Run()
		case "start":
			err = s.Start()
			//err = service.Control(s, os.Args[1])
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			log.Printf("【ServerStart】 Service %s started.\n", svcConfig.DisplayName)
		case "restart":
			err = s.Restart()
			if err != nil {
				fmt.Printf("Failed to restart: %s\n", err)
				return
			}
			log.Printf("【ServerRestart】 Service %s restarted.\n", svcConfig.DisplayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			log.Printf("【ServerStopped】 Service %s stop.\n", svcConfig.DisplayName)
		case "status":
			sta, err := s.Status()
			if err != nil {
				fmt.Printf("Failed to status: %s\n", err)
				return
			}
			var status = "StatusUnknown"
			if sta == service.StatusRunning {
				status = "Running"
			} else if sta == service.StatusStopped {
				status = "Stopped"
			}
			log.Printf("【ServerStatus】 Service %s  status=%s \n", svcConfig.DisplayName, status)
		case "v":
			log.Printf("【ServerStatus】 Service %s  version=%s \n", svcConfig.DisplayName, version)
		}
		return
	} else {
		log.Printf("【ServerRun】 服务 %s 启动\n", svcConfig.DisplayName)
		var err = s.Run()
		log.Printf("【ServerRun】 服务 %s 启动成功\n", svcConfig.DisplayName)
		if err != nil {
			fmt.Println("启动失败", err.Error())
		}
	}
}

type BaseService struct {
}

func (svr *BaseService) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	if service.Interactive() {
		log.Printf("【Start】 Running in terminal")
	} else {
		log.Printf("【Start】 Running under service manager")
	}
	log.Printf("【Start】 启动服务")
	go svr.run()
	log.Printf("【Start】 启动服务成功")
	return nil
}

func (svr *BaseService) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	_, _ = s.Status()
	log.Printf("【Clean】监听程序")
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (svr *BaseService) run() (err error) {
	setup()
	if config.Configs.Env == "prd" {
		gin.SetMode(gin.ReleaseMode)
	}
	InitRouters()
	return nil
}

func InitRouters() {
	r := gin.Default()
	r.Use(Middlewares.Cors())
	router.InitAllRouter(r)
	//appDir := utils.GetAppDir()
	//if runtime.GOOS != "windows" {
	//	appDir = utils.GetCurrentPath()
	//}

	//r.LoadHTMLGlob(appDir + "/app/templates/**/*")
	//r.Static("/statics", appDir+"/app/templates/statics")
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	r.Run(":" + config.Configs.HttpListenPort)
}

func setup() {
	config.Default()
	//mysql.Set_up()
	//redis.Default()
}
