package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"todayschool/config"
)

type ProcessingList struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		TotalSize  int `json:"totalSize"`
		PageSize   int `json:"pageSize"`
		PageNumber int `json:"pageNumber"`
		Rows       []struct {
			Wid            string `json:"wid"`
			FormWid        string `json:"formWid"`
			Priority       string `json:"priority"`
			Subject        string `json:"subject"`
			Content        string `json:"content"`
			SenderUserName string `json:"senderUserName"`
			CreateTime     string `json:"createTime"`
			StartTime      string `json:"startTime"`
			EndTime        string `json:"endTime"`
			CurrentTime    string `json:"currentTime"`
			IsHandled      int    `json:"isHandled"`
			IsRead         int    `json:"isRead"`
		} `json:"rows"`
	} `json:"datas"`
}

type ProcessJson struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
}

var base_querycontroller = "/wec-counselor-collector-apps/stu/collector/queryCollectorProcessingList"

func PostProcessingList(con *config.RootConfig) (string, int) {
	pro := ProcessJson{PageNumber: 1, PageSize: 20}
	jsons, _ := json.Marshal(pro)
	post_controller := "https://" + con.BaseUrl + base_querycontroller
	log.Println("请求路径", post_controller)
	req, err := http.NewRequest("POST", post_controller, bytes.NewReader(jsons))

	if err != nil {
		log.Println("请求出错", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10;"+con.PhoneModel+" Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.101 Mobile Safari/537.36 cpdaily/9.0.12 wisedu/9.0.12")
	req.Header.Set("Cookie", "MOD_AUTH_CAS="+con.MOD_AUTH_CAS)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("HOST", con.BaseUrl)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("请求出错:", err)
	}
	defer resp.Body.Close()
	list := new(ProcessingList)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &list)

	if len(list.Datas.Rows) > 0 {
		for k, v := range list.Datas.Rows {
			log.Println(k+1, "条待处理")
			log.Println("发布人：", v.SenderUserName)
			log.Println("创建时间：", v.CreateTime)
			log.Println("开始时间：", v.StartTime)
			log.Println("结束时间：", v.EndTime)
			log.Println("当前时间：", v.CurrentTime)
			log.Println("需要wid：", v.Wid)
			log.Println("需要formWid：", v.FormWid)
			log.Println("是否提交过了（0：没有提交，1：提交过了）：", v.IsHandled)
			return v.Wid, v.IsHandled
			//PostDetail(v.Wid)
		}
	} else {
		log.Println("没有需要处理的表单......忽略，等待下一次执行")
	}
	return "", 1
}
