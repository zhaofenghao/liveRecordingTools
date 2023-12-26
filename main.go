package main

import (
	"live_recording_tools/config"
	"live_recording_tools/daemon"
)

func main() {
	config.Default()
	daemon.RunService()
}
