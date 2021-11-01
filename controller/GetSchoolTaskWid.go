package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"todayschool/config"
)

type DetailList struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Collector struct {
			Wid                 string      `json:"wid"`
			FormWid             string      `json:"formWid"`
			Priority            string      `json:"priority"`
			EndTime             string      `json:"endTime"`
			CurrentTime         string      `json:"currentTime"`
			SchoolTaskWid       string      `json:"schoolTaskWid"`
			IsConfirmed         int         `json:"isConfirmed"`
			SenderUserName      string      `json:"senderUserName"`
			CreateTime          string      `json:"createTime"`
			AttachmentUrls      interface{} `json:"attachmentUrls"`
			AttachmentNames     interface{} `json:"attachmentNames"`
			AttachmentSizes     interface{} `json:"attachmentSizes"`
			IsUserSubmit        int         `json:"isUserSubmit"`
			FetchStuLocation    bool        `json:"fetchStuLocation"`
			IsLocationFailedSub bool        `json:"isLocationFailedSub"`
			Address             interface{} `json:"address"`
		} `json:"collector"`
		Form struct {
			Wid            string        `json:"wid"`
			FormType       string        `json:"formType"`
			FormTitle      string        `json:"formTitle"`
			ExamTime       interface{}   `json:"examTime"`
			FormContent    string        `json:"formContent"`
			BackReason     interface{}   `json:"backReason"`
			IsBack         int           `json:"isBack"`
			Attachments    []interface{} `json:"attachments"`
			Score          float64       `json:"score"`
			StuScore       int           `json:"stuScore"`
			ConfirmDesc    string        `json:"confirmDesc"`
			IsshowOrdernum int           `json:"isshowOrdernum"`
			IsAnonymous    int           `json:"isAnonymous"`
			IsallowUpdated int           `json:"isallowUpdated"`
			IsshowScore    int           `json:"isshowScore"`
			IsshowResult   int           `json:"isshowResult"`
		} `json:"form"`
	} `json:"datas"`
}

type Wid struct {
	CollectorWid string `json:"collectorWid"`
}

var base_detailcontroller = "/wec-counselor-collector-apps/stu/collector/detailCollector"

func PostDetail(wid string, con *config.RootConfig) (string, string) {
	wid_send := Wid{CollectorWid: wid}
	jsons, _ := json.Marshal(wid_send)
	post_controller := "https://" + con.BaseUrl + base_detailcontroller
	log.Println("请求路径", post_controller)
	req, err := http.NewRequest("POST", post_controller, bytes.NewReader(jsons))

	if err != nil {
		log.Println("请求出错", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; "+con.PhoneModel+" Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.101 Mobile Safari/537.36 cpdaily/9.0.12 wisedu/9.0.12")
	req.Header.Set("Cookie", "MOD_AUTH_CAS="+con.MOD_AUTH_CAS)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("HOST", con.BaseUrl)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("请求出错:", err)
	}
	defer resp.Body.Close()
	list := new(DetailList)
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &list)
	//fmt.Println(string(body))
	log.Println("表单标题：", list.Datas.Form.FormTitle)
	log.Println("发布人：", list.Datas.Collector.SenderUserName)
	log.Println("需要SchoolTaskWid：", list.Datas.Collector.SchoolTaskWid)
	log.Println("需要FormWid：", list.Datas.Collector.FormWid)
	return list.Datas.Collector.FormWid, list.Datas.Collector.SchoolTaskWid
}
