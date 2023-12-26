package main

import (
	"ffmpeg_work/config"
	"ffmpeg_work/daemon"
)

func main() {
	config.Default()
	daemon.RunService()
}
