package validator

import (
	"ffmpeg_work/pkg/e"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

//var json struct {
//	Id       int    `json:"id" binding:"required"`
//	Username string `json:"username" binding:"required"`
//}

func UserAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := e.Gin{C: c}
		uid := c.DefaultPostForm("uid", "")
		userName := c.DefaultPostForm("user_name", "")
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 请求前
		//c.Next()

		// 请求后
		//if err := c.Bind(&json); err != nil {
		//	res.Fail(-3, err.Error(), nil)
		//	c.Abort()
		//}

		latency := time.Since(t)
		if uid == "" {
			res.Fail(-1, "no uid!", nil)
			c.Abort()
		}

		if userName == "" {
			res.Fail(-1, "no user_name!", nil)
			c.Abort()
		}
		log.Print(latency)

		// 获取发送的 status
		status := c.Writer.Status()
		log.Println(status)
	}
}
