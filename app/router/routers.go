package router

import (
	"ffmpeg_work/app/controller"
	"github.com/gin-gonic/gin"
)

func InitAllRouter(r *gin.Engine) {
	r.GET("/run", controller.Run)
}
