package router

import (
	"github.com/gin-gonic/gin"
	"live_recording_tools/app/controller"
)

func InitAllRouter(r *gin.Engine) {
	r.GET("/run", controller.Run)
}
