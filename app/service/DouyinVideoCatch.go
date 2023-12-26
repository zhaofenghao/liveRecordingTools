package service

import (
	"encoding/json"
	"ffmpeg_work/app/utils"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type DouyinRoomInfo struct {
	AnchorName string `json:"anchor_name"`
	RoomInfo   struct {
		Room struct {
			//AdminUserIds    []int64  `json:"admin_user_ids"`
			//AdminUserIdsStr []string `json:"admin_user_ids_str"`
			Status    int    `json:"status"`
			StatusStr string `json:"status_str"`
			StreamUrl struct {
				FlvPullUrl struct {
					FULLHD1 string `json:"FULL_HD1"`
					HD1     string `json:"HD1"`
					SD1     string `json:"SD1"`
					SD2     string `json:"SD2"`
				} `json:"flv_pull_url"`
				HlsPullUrl    string `json:"hls_pull_url"`
				HlsPullUrlMap struct {
					FULLHD1 string `json:"FULL_HD1"`
					HD1     string `json:"HD1"`
					SD1     string `json:"SD1"`
					SD2     string `json:"SD2"`
				} `json:"hls_pull_url_map"`
			} `json:"stream_url"`
			Title string `json:"title"`
		} `json:"room"`
	} `json:"roomInfo"`
	Store interface{} `json:"store"`
}

// 获取推流&房间信息
func GetStreamData(url string) DouyinRoomInfo {
	//res := e.Gin{C: c}
	defer func() {
		if er := recover(); er != nil {
			log.Printf("采集直播间数据错误 =====> %s", er)
			return
		}
	}()
	var header = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Referer":         "https://live.douyin.com/",
		"Cookie":          "ttwid=1%7CB1qls3GdnZhUov9o2NxOMxxYS2ff6OSvEWbv0ytbES4%7C1680522049%7C280d802d6d478e3e78d0c807f7c487e7ffec0ae4e5fdd6a0fe74c3c6af149511; my_rd=1; passport_csrf_token=3ab34460fa656183fccfb904b16ff742; passport_csrf_token_default=3ab34460fa656183fccfb904b16ff742; d_ticket=9f562383ac0547d0b561904513229d76c9c21; n_mh=hvnJEQ4Q5eiH74-84kTFUyv4VK8xtSrpRZG1AhCeFNI; store-region=cn-fj; store-region-src=uid; LOGIN_STATUS=1; __security_server_data_status=1; FORCE_LOGIN=%7B%22videoConsumedRemainSeconds%22%3A180%7D; pwa2=%223%7C0%7C3%7C0%22; download_guide=%223%2F20230729%2F0%22; volume_info=%7B%22isUserMute%22%3Afalse%2C%22isMute%22%3Afalse%2C%22volume%22%3A0.6%7D; strategyABtestKey=%221690824679.923%22; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22%3Atrue%2C%5C%22screen_width%5C%22%3A1536%2C%5C%22screen_height%5C%22%3A864%2C%5C%22browser_online%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A8%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downlink%5C%22%3A10%2C%5C%22effective_type%5C%22%3A%5C%224g%5C%22%2C%5C%22round_trip_time%5C%22%3A150%7D%22; VIDEO_FILTER_MEMO_SELECT=%7B%22expireTime%22%3A1691443863751%2C%22type%22%3Anull%7D; home_can_add_dy_2_desktop=%221%22; __live_version__=%221.1.1.2169%22; device_web_cpu_core=8; device_web_memory_size=8; xgplayer_user_id=346045893336; csrf_session_id=2e00356b5cd8544d17a0e66484946f28; odin_tt=724eb4dd23bc6ffaed9a1571ac4c757ef597768a70c75fef695b95845b7ffcd8b1524278c2ac31c2587996d058e03414595f0a4e856c53bd0d5e5f56dc6d82e24004dc77773e6b83ced6f80f1bb70627; __ac_nonce=064caded4009deafd8b89; __ac_signature=_02B4Z6wo00f01HLUuwwAAIDBh6tRkVLvBQBy9L-AAHiHf7; ttcid=2e9619ebbb8449eaa3d5a42d8ce88ec835; webcast_leading_last_show_time=1691016922379; webcast_leading_total_show_times=1; webcast_local_quality=sd; live_can_add_dy_2_desktop=%221%22; msToken=1JDHnVPw_9yTvzIrwb7cQj8dCMNOoesXbA_IooV8cezcOdpe4pzusZE7NB7tZn9TBXPr0ylxmv-KMs5rqbNUBHP4P7VBFUu0ZAht_BEylqrLpzgt3y5ne_38hXDOX8o=; msToken=jV_yeN1IQKUd9PlNtpL7k5vthGKcHo0dEh_QPUQhr8G3cuYv-Jbb4NnIxGDmhVOkZOCSihNpA2kvYtHiTW25XNNX_yrsv5FN8O6zm3qmCIXcEe0LywLn7oBO2gITEeg=; tt_scid=mYfqpfbDjqXrIGJuQ7q-DlQJfUSG51qG.KUdzztuGP83OjuVLXnQHjsz-BRHRJu4e986",
	}
	//connectReq := &http.Request{
	//	Method: "CONNECT",
	//	URL:    &url2.URL{Opaque: originUrl},
	//	Header: header,
	//	Host:   originUrl,
	//}

	r := utils.Request{
		Url:     url,
		Method:  "GET",
		Headers: header,
	}
	resp := r.Send()
	if resp.StatusCode != http.StatusOK {
		//fmt.Println(resp.Reason)
		//return
	}

	content := resp.Text
	//res.Success("ok", content)

	match_json_str := regexp.MustCompile(`(\{\\"state\\":.*?)]\\n"]\)`).FindString(content)
	if match_json_str == "" {
		match_json_str = regexp.MustCompile(`(\{\\"common\\":.*?)]\\n"]\)</script><div hidden`).FindString(content)
	}

	cleaned_string := strings.Replace(match_json_str, "\\", "", -1)
	cleaned_string = strings.Replace(cleaned_string, "u0026", "&", -1)

	re := regexp.MustCompile("\"roomStore\":(.*?),\"linkmicStore\"")
	room_store := re.FindStringSubmatch(cleaned_string)[1]

	//fmt.Println("room_store ==>", room_store)
	re = regexp.MustCompile("\"nickname\":\"(.*?)\",\"avatar_thumb")
	anchor_name := re.FindStringSubmatch(room_store)[1]

	room_store = strings.Split(room_store, ",\"has_commerce_goods\"")[0] + "}}}"
	var roomInfo map[string]interface{}
	json.Unmarshal([]byte(room_store), &roomInfo)

	var roomData DouyinRoomInfo
	json.Unmarshal([]byte(room_store), &roomData)
	roomData.AnchorName = anchor_name
	return roomData
}

type StreamResult struct {
	AnchorName string `json:"anchor_name"`
	IsLive     bool   `json:"is_live"`
	M3u8Url    string `json:"m3u8_url"`
	FlvUrl     string `json:"flv_url"`
	RecordUrl  string `json:"record_url"`
}

// 获取抖音直播源地址
func GetStreamUrl(roomData DouyinRoomInfo) StreamResult {
	anchor_name := roomData.AnchorName
	var streamResult StreamResult
	streamResult.AnchorName = anchor_name
	streamResult.IsLive = false

	status := roomData.RoomInfo.Room.Status
	if status == 2 {
		streamUrl := roomData.RoomInfo.Room.StreamUrl
		flvUrlList := streamUrl.FlvPullUrl
		m3u8UrlList := streamUrl.HlsPullUrlMap

		mapM3u8Data := utils.StructToMap(m3u8UrlList)
		var qualityList []string
		for k, _ := range mapM3u8Data {
			qualityList = append(qualityList, k)
		}
		//if len(qualityList) < 4 {
		//	qualityList = append(qualityList, qualityList[len(qualityList)-1])
		//}

		//"原画": "FULL_HD1",
		//"蓝光": "FULL_HD1",
		//"超清": "HD1",
		//"高清": "SD1",
		//"标清": "SD2",

		videoQulities := map[string]int{
			"原画": 0, "蓝光": 0, "超清": 1, "高清": 2, "标清": 3,
		}
		qualityIndex := videoQulities["原画"] // 先写死原画
		qualityKey := qualityList[qualityIndex]
		m3u8Url := mapM3u8Data[qualityKey]
		if m3u8Url == "" {
			log.Println("未获取到m3u8 url，自动重新获取")
			m3u8Url = getStreamUrl(mapM3u8Data, "")
			log.Println("重新获取m3u8 url ==>" + m3u8Url)
		}

		mapFlvUrlData := utils.StructToMap(flvUrlList)
		flvUrl := mapFlvUrlData[qualityKey]
		if flvUrl == "" {
			log.Println("未获取到flv url，自动重新获取")
			flvUrl = getStreamUrl(mapFlvUrlData, "")
			log.Println("重新获取flv url ==>" + flvUrl)
		}

		streamResult.M3u8Url = m3u8Url
		streamResult.FlvUrl = flvUrl
		streamResult.IsLive = true
		streamResult.RecordUrl = m3u8Url // 使用 m3u8 链接进行录制
	}

	return streamResult
}

func getStreamUrl(data map[string]string, want string) string {
	if data[want] != "" {
		return data[want]
	} else {
		if data["FULL_HD1"] != "" {
			return data["FULL_HD1"]
		}
		if data["HD1"] != "" {
			return data["HD1"]
		}
		if data["SD1"] != "" {
			return data["SD1"]
		}
		if data["SD2"] != "" {
			return data["SD2"]
		}
	}

	return ""
}
