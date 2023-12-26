package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"live_recording_tools/app/service"
	"live_recording_tools/app/utils"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"
)

var monitoring = 0                             // 监控个数
var recording []string                         // 监控的直播间名
var runningList = new(sync.Map)                // 监控的直播间url
var recordingTimeList = make(map[string]int64) // 监控的直播间开始时间
var rstr = "[\\/\\\\:\\*\\?\"\\<\\>\\|&u]"
var splitVideoByTime = true // 是否自动切片
var ctxMap = make(map[string]*context.CancelFunc)
var configPath = utils.GetAppDir() + "/config/config.json"
var FfmpegExe = utils.GetAppDir() + "\\ffmpeg"
var displayRepet = false

const RefreshTime = 10 // 直播间刷新间隔

func NowTime() string {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location()).Format("2006-01-02_15-04-05")
	return newTime
}

func Run(c *gin.Context) {
	go Recording()
	Display()
}

func Recording() {
	for {
		config, err := utils.ReadConfigList(configPath)
		if err != nil {
			log.Println("读取配置文件错误，请检查文件配置：", err.Error())
			return
		}

		var configUrls []string
		// 动态添加直播间监听
		for i, config := range config.Config {
			if config.Url == "" {
				continue
			}
			configUrls = append(configUrls, config.Url)
			v, _ := runningList.Load(config.Url)
			if v == nil {
				if config.Status != 1 {
					continue
				}
				ctx, cancel := context.WithCancel(context.Background())
				go SaveDouyinLiveVideo(ctx, config.Url, i)
				monitoring++
				ctxMap[config.Url] = &cancel
				runningList.Store(config.Url, config)
			} else {
				c := v.(utils.DouyunRoomConfig)
				c.Status = config.Status
				runningList.Store(config.Url, c)
			}
		}

		// 动态关闭直播间监听
		runningList.Range(func(_, v any) bool {
			p := v.(utils.DouyunRoomConfig)
			if !utils.In(p.Url, configUrls) || p.Status != 1 {
				cancelFunc := *ctxMap[p.Url]
				cancelFunc() // 调用 cancel 函数
				runningList.Delete(p.Url)
				monitoring--
			}
			return true
		})
		time.Sleep(time.Second * RefreshTime)
	}
}

// 日志打印
func Display() {
	if displayRepet == true {
		return
	}
	displayRepet = true
	for {
		time.Sleep(10 * time.Second)
		var printStr = fmt.Sprintf("\n共监测%d直播间 \n", monitoring)
		if len(recording) == 0 {
			printStr += "没有正在录制的直播 \n"
			log.Println(printStr)
			continue
		} else {
			//printStr += fmt.Sprintf("录制视频质量：%s \n", videoQuality)
			//log.Println(fmt.Sprintf("录制视频质量：%s", videoQuality))
			var records []string
			for _, r := range recording {
				if !utils.In(r, records) {
					records = append(records, r)
				}
			}
			printStr += fmt.Sprintf("正在录制%d个直播 ：\n", len(records))
			//log.Println(fmt.Sprintf("正在录制%d个直播: ", len(records)))
			for _, recordingLive := range records {
				haveRecordTime := time.Now().Unix() - recordingTimeList[recordingLive]
				//log.Println(fmt.Sprintf("%s 正在录制中，已录制：%d秒", recordingLive, haveRecordTime))
				printStr += fmt.Sprintf("  %s 正在录制中，已录制：%d秒 \n", recordingLive, haveRecordTime)
			}
			log.Println(printStr)
		}
	}
}

// 获取抖音直播录屏
func SaveDouyinLiveVideo(ctx context.Context, url string, countVariable int) {
	// 死循环，一直查询直播间状态
	var ffmpegCmd *exec.Cmd
	var noError = true
	var recordName string
	var recordFinished2 = false
	for {
		select {
		case <-ctx.Done():
			log.Printf("2关闭此直播间监听: %s ", recordName)
			FilterRecording(recordName)
			return
		default:
			roomData := service.GetStreamData(url)
			portInfo := service.GetStreamUrl(roomData)

			var warningCount = 0
			if portInfo.AnchorName == "" {
				log.Printf("序号%v网址内容获取失败,进行重试中...获取失败的地址是:%s", countVariable, url)
				warningCount++
			} else {

				re := regexp.MustCompile(rstr)
				anchorName := re.ReplaceAllString(portInfo.AnchorName, "_")
				recordName = fmt.Sprintf("序号%v %s", countVariable, anchorName)

				v, _ := runningList.Load(url)
				if v == nil {
					config := v.(utils.DouyunRoomConfig)
					if config.Status == 1 {
						return
					}
				}

				if utils.In(recordName, recording) {
					log.Printf("新增直播间%s已经存在，此条略过", anchorName)
					return
				}

				if !portInfo.IsLive {
					log.Printf("%s 等待直播... ", recordName)
					FilterRecording(recordName)
					return
				} else {
					if utils.In(recordName, recording) {
						return
					}

					content := fmt.Sprintf("%s  正在直播中...", recordName)
					log.Printf(content)

					// 推送@TODO
					realUrl := portInfo.RecordUrl
					fullPath := utils.GetAppDir() + "\\video\\" + anchorName

					utils.CheckFile(fullPath)
					if realUrl != "" {
						_, err := os.Stat(fullPath)
						if os.IsNotExist(err) {
							err = os.Mkdir(fullPath, 0755)
							if err != nil {
								log.Printf("创建文件夹失败：%s, err: ", fullPath, err.Error())
								return
							}
						}

						// 先录制MP4
						userAgent := "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36"
						ffmpegCommand := []string{
							"-y",
							"-v", "verbose",
							"-rw_timeout", "15000000", // 15s
							"-loglevel", "error",
							"-hide_banner",
							"-user_agent", userAgent,
							"-protocol_whitelist", "rtmp,crypto,file,http,https,tcp,tls,udp,rtp",
							"-thread_queue_size", "1024",
							"-analyzeduration", "2147483647",
							"-probesize", "2147483647",
							"-fflags", "+discardcorrupt",
							"-i", realUrl,
							"-bufsize", "10000k",
							"-sn", "-dn",
							"-reconnect_delay_max", "30",
							"-reconnect_streamed", "-reconnect_at_eof",
							"-max_muxing_queue_size", "64",
							"-correct_ts_overflow", "1",
						}

						// 代理添加 @TODO
						recording = append(recording, recordName)
						recordingTimeList[recordName] = time.Now().Unix()
						recInfo := fmt.Sprintf("%s 录制视频中: %s", anchorName, fullPath)

						filename := anchorName + "_" + NowTime() + ".ts"
						log.Printf("%s/%s", recInfo, filename)
						saveFilePath := fullPath + "\\" + filename
						var command []string
						if splitVideoByTime {
							saveFilePath = fmt.Sprintf("%s\\%s_%s_%%03d.ts", fullPath, anchorName, NowTime())
							command = []string{
								"-c:v", "copy",
								"-c:a", "aac",
								"-map", "0",
								"-f", "segment",
								"-segment_time", "3600", // 切片时间间隔
								"-segment_format", "ts",
								"-movflags", "+faststart",
								"-reset_timestamps", "1",
								saveFilePath,
							}
						} else {
							command = []string{
								"-map", "0",
								"-c:v", "copy",
								"-c:a", "copy",
								//"-f", "mp4",
								"-f", "mpegts",
								saveFilePath,
							}
						}

						ffmpegCommand = append(ffmpegCommand, command...)
						ffmpegCmd = exec.CommandContext(ctx, FfmpegExe, ffmpegCommand...) //

						// 设置标准输出和标准错误输出
						//ffmpegCmd.Stdout = os.Stdout
						//ffmpegCmd.Stderr = os.Stderr

						_, e := ffmpegCmd.CombinedOutput() // 阻塞
						if e != nil {
							fmt.Println("CombinedOutput error", e)
							recordFinished2 = true
						}

						// 将ts文件转为mp4
						//time.Sleep(1 * time.Second)
						//mp4Path := fullPath + "\\" + anchorName + "_" + NowTime() + ".mp4"
						//TransVideoType(saveFilePath, mp4Path)
					}

					if recordFinished2 {
						FilterRecording(recordName)
						if noError {
							//runningList.Delete(url)
							log.Printf("%s 直播录制完成", anchorName)
						} else {
							log.Printf("直播录制出错,请检查网络")
						}
						recordFinished2 = false
					}
				}

			}
		}

		// 刷新直播间状态
		time.Sleep(RefreshTime * time.Second)
	}

}

// 转换视频格式
func TransVideoType(tsPath, mp4Path string) error {
	// 将ts文件转为mp4
	time.Sleep(1 * time.Second)
	cmd := exec.Command(FfmpegExe, "-i", tsPath, "-c:v", "copy", "-c:a", "copy", mp4Path)

	// 设置输出和错误流
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return err
	} else {
		//删除源ts文件
		err = os.Remove(tsPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// 插入的视频也有声音
//var filter = fmt.Sprintf("%s%s%s%s%s%s%s",
//	"[1:v]scale=300:-1[overlay_image];",
//	"[2:v]scale=300:-1[overlay_video];",
//	"[0:v][overlay_image]overlay=10:H-310[v1];",
//	"[v1]drawtext=fontcolor=0xFF0000:fontfile=D:/work/live_recording_tools/font/simfang.ttf:fontsize=94:x=(W-text_w)/2:text='买买买买买！！！':y=(H-text_h)/2,",
//	"drawtext=fontcolor=#00BFFF:fontfile=D:/work/live_recording_tools/font/simkai.ttf:fontsize=94:x=(W-text_w)/2:text='抢抢抢！！！':y=(H-text_h)/2-100[v2];",
//	//"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10:eof_action=endall[v]",
//	"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10[v];",
//	"[0:a][2:a]amix=inputs=2[aout]",
//)
//
//var commond = []string{
//	"-y",
//	"-i", videoPath,
//	"-i", imgPath,
//	"-i", inputPath,
//	//"-filter_complex", "overlay=W-w-10:H-h-10,drawtext=text='Hello, World':fontfile=D:\\ttf\\MicrosoftYaHei.ttf:fontsize=44:x=50:y=100:fontcolor=black",
//	//"-filter_complex", "overlay=W-w-10:H-h-10,drawtext=text='Hello, World':fontfile=D:\\ttf\\MicrosoftYaHei.ttf:fontsize=64:x=(W-text_w)/2:y=(H-text_h)/2:fontcolor=0xFF0000",
//	"-filter_complex", filter,
//	"-c:a", "aac",
//	"-c:v", "libx264",
//	"-map", "[v]",
//	"-map", "[aout]",
//	outPutPath,
//}

// ===========================================
// 主视频有声音
//var filter = fmt.Sprintf("%s%s%s%s%s%s",
//	"[1:v]scale=300:-1[overlay_image];",
//	"[2:v]scale=300:-1[overlay_video];",
//	"[0:v][overlay_image]overlay=10:H-310[v1];",
//	"[v1]drawtext=fontcolor=0xFF0000:fontfile=D:/work/live_recording_tools/font/simfang.ttf:fontsize=94:x=(W-text_w)/2:text='买买买买买！！！':y=(H-text_h)/2,",
//	"drawtext=fontcolor=#00BFFF:fontfile=D:/work/live_recording_tools/font/simkai.ttf:fontsize=94:x=(W-text_w)/2:text='抢抢抢！！！':y=(H-text_h)/2-100[v2];",
//	"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10:eof_action=endall[v]",
//)

//var commond = []string{
//	"-y",
//	"-i", videoPath,
//	"-i", imgPath,
//	"-i", inputPath,
//	//"-filter_complex", "overlay=W-w-10:H-h-10,drawtext=text='Hello, World':fontfile=D:\\ttf\\MicrosoftYaHei.ttf:fontsize=44:x=50:y=100:fontcolor=black",
//	//"-filter_complex", "overlay=W-w-10:H-h-10,drawtext=text='Hello, World':fontfile=D:\\ttf\\MicrosoftYaHei.ttf:fontsize=64:x=(W-text_w)/2:y=(H-text_h)/2:fontcolor=0xFF0000",
//	"-filter_complex", filter,
//	"-c:a", "copy",
//	"-map", "[v]",
//	"-map", "0:a", //在之前的命令中，-map 0:a 指定要将输入文件中的第一个音频流（0:a 表示第一个音频流，0 是输入文件的索引）包含在输出文件中。这是为了确保输出文件包含了主视频的音频流。
//	outPutPath,
//}

// 向视频中添加文字、图片
func AddEleToVideo() {
	videoPath := ""  // 原视频路径
	imgPath := ""    // 插入图片路径
	inputPath := ""  // 插入视频路径
	outPutPath := "" // 输出视频路径

	var filter = fmt.Sprintf("%s%s%s%s%s%s",
		"[1:v]scale=300:-1[overlay_image];",
		"[2:v]scale=300:-1[overlay_video];",
		"[0:v][overlay_image]overlay=10:H-310[v1];",
		"[v1]drawtext=fontcolor=0xFF0000:fontfile=D:/work/live_recording_tools/font/simfang.ttf:fontsize=94:x=(W-text_w)/2:text='买买买买买！！！':y=(H-text_h)/2,",
		"drawtext=fontcolor=#00BFFF:fontfile=D:/work/live_recording_tools/font/simkai.ttf:fontsize=94:x=(W-text_w)/2:text='抢抢抢！！！':y=(H-text_h)/2-100[v2];",
		"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10:eof_action=endall[v]",
	)

	//var filter = fmt.Sprintf("%s%s%s%s%s%s",
	//	"[1:v]scale=300:-1[overlay_image];",
	//	"[2:v]scale=300:-1[overlay_video];",
	//	"[0:v][overlay_image]overlay=10:H-310[v1];",
	//	"[v1]drawtext=fontcolor=0xFF0000:fontfile=D:/work/live_recording_tools/font/simfang.ttf:fontsize=94:x=(W-text_w)/2:text='买买买买买！！！':y=(H-text_h)/2,",
	//	"drawtext=fontcolor=#00BFFF:fontfile=D:/work/live_recording_tools/font/simkai.ttf:fontsize=94:x=(W-text_w)/2:text='抢抢抢！！！':y=(H-text_h)/2-100[v2];",
	//	//"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10:eof_action=endall[v]",
	//	"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10:eof_action=endall[v]",
	//)

	//
	//var filter = fmt.Sprintf("%s%s%s%s%s",
	//	"[1:v]scale=300:-1[overlay_image];",
	//	"[2:v]scale=300:-1[overlay_video];",
	//	"[0:v][overlay_image]overlay=10:H-310[v1];",
	//	"[v1]drawtext=fontcolor=0xFF0000:fontfile=D:/work/live_recording_tools/font/simfang.ttf:fontsize=94:x=(W-text_w)/2:text='买买买买买！！！':y=(H-text_h)/2[v2];",
	//	"[v2][overlay_video]overlay=W-overlay_w-10:H-overlay_h-10:eof_action=endall[v]",
	//)
	var commond = []string{
		"-y",
		"-i", videoPath,
		"-i", imgPath,
		"-i", inputPath,
		"-filter_complex", filter,
		"-c:a", "copy",
		"-map", "[v]",
		"-map", "0:a", //在之前的命令中，-map 0:a 指定要将输入文件中的第一个音频流（0:a 表示第一个音频流，0 是输入文件的索引）包含在输出文件中。这是为了确保输出文件包含了主视频的音频流。
		outPutPath,
	}
	cmd := exec.Command(FfmpegExe, commond...)

	// 设置输出和错误流
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	cmd.Run()
}

func FilterRecording(recordName string) {
	for i, s := range recording {
		if s == recordName {
			recording = append(recording[:i], recording[i+1:]...)
		}
	}
}
