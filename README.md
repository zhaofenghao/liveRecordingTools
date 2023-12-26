# 直播间录制工具
# liveRecording 当前版本【1.0.0】
基于 Gin框架，ffmpeg的直播录屏工具，目前支持windows平台上抖音直播录制。
目前功能较少，后续会持续更新~

### 如何使用

```
go mod init live_recording_tools
go mod tidy
go run main.go
```
#### 修改config/config.json配置文件，填入你需要录制的直播间配置
```
{"save_path":"","check_frequency": 30,"config": [{"is_monitor": true,"status": 1,"url": "https://live.douyin.com/80017709309", "sort": 0},{"is_monitor": false,"status": 1,"url": "https://live.douyin.com/574697141527", "sort": 0}]}
```
#### 启动监控
```
http://localhost:8075/run
```


