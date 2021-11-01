package controller

import (
	"bytes"
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"todayschool/config"
	"todayschool/mail"
)

//提交表单
type SubmitForm struct {
	FormWid       string  `json:"formWid"`
	Address       string  `json:"address"`
	CollectWid    string  `json:"collectWid"`
	SchoolTaskWid string  `json:"schoolTaskWid"`
	Form          []Forms `json:"form"`
	UaIsCpadaily  bool    `json:"uaIsCpadaily"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

//提交响应
type SubSuccessResp struct {
	Code              string      `json:"code"`
	Message           string      `json:"message"`
	Wid               interface{} `json:"wid"`
	HasWindowLocation int         `json:"hasWindowLocation"`
	WindowLocation    string      `json:"windowLocation"`
	Score             interface{} `json:"score"`
}

//加密aes
type SubmitFormAes struct {
	AppVersion    string  `json:"appVersion"`
	SystemName    string  `json:"systemName"`
	BodyString    string  `json:"bodyString"`
	Sign          string  `json:"sign"` //appVersion=9.0.12&bodyString=bodyString&deviceId=deviceId&lat=lat&lon=lon&model=OPPO R11 Plus&systemName=android&systemVersion=4.4.4&userId=username&ytUQ7l2ZZu8mLvJZ
	Model         string  `json:"model"`
	Lon           float64 `json:"lon"`
	CalVersion    string  `json:"calVersion"`
	SystemVersion string  `json:"systemVersion"`
	DeviceID      string  `json:"deviceId"`
	UserID        string  `json:"userId"`
	Version       string  `json:"version"`
	Lat           float64 `json:"lat"`
}

var base_submitformscontroller = "/wec-counselor-collector-apps/stu/collector/submitForm"

func SubmitFormFileds(form *SubmitFormAes, str string, con *config.RootConfig) {
	jsons, _ := json.Marshal(form)
	post_controller := "https://" + con.BaseUrl + base_submitformscontroller
	log.Println("请求路径", post_controller)
	req, err := http.NewRequest("POST", post_controller, bytes.NewReader(jsons))

	if err != nil {
		log.Println("请求出错", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; "+con.PhoneModel+" Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.101 Mobile Safari/537.36 cpdaily/9.0.12 wisedu/9.0.12")
	req.Header.Set("Cookie", "MOD_AUTH_CAS="+con.MOD_AUTH_CAS)
	req.Header.Set("Cpdaily-Extension", con.Cpdaily_Extension)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("HOST", con.BaseUrl)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("请求出错:", err)
	}
	defer resp.Body.Close()
	ret := new(SubSuccessResp)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &ret)

	title := "来自不知名靓仔"
	subject := "靓仔提醒你"
	html := "<html><body ><div  style=\"width:100%;height:800px;color:#60599a;background:-webkit-linear-gradient(323deg,#70f7fe,#fbd7c6,#fdefac,#bfb5dd,#bed5f5);\"><h2 style=\"text-align:center\">你好，你的今日校园填写成功啦！！！</h2><br><div style=\"color:red\"><h3 >" + str + "记住，请勿登录今日校园！！！</h3><h3 >记住，请勿登录今日校园！！！</h3><h3 >记住，请勿登录今日校园！！！</h3></div>重要事情说三遍！！！</div></body></html>"

	config := config.GetConfig()

	if strings.Compare(ret.Message, "SUCCESS") == 0 {
		for _, v := range config.UserEmail {
			err := mail.SendMail(title, v, subject, html, "html")
			if err != nil {
				log.Println(ret.Message, "----", v, "邮件发送失败！", err)
				color.New(color.BgHiRed).Println( "填写成功!!!----", v, "邮件发送失败！", err)
			}
			log.Println(ret.Message, "----", v, "邮件发送成功！！")
			color.New(color.FgHiGreen).Println("填写成功!!!----", v, "邮件发送成功！！")
		}
		color.New(color.FgHiYellow).Println("提交成功！")

	} else {
		html = "<html><body ><div  style=\"width:100%;height:800px;color:#60599a;background: -webkit-linear-gradient(323deg, #70f7fe, #fbd7c6, #fdefac, #bfb5dd, #bed5f5);\"><h2 style=\"color:red;\">今日校园填写失败了！！！请确保你没有在使用程序过程中登录了今日校园，否则请删除配置文件重试:</h2><br>" + ret.Message + " </div></body></html>"
		for _, v := range config.UserEmail {
			mail.SendMail(title, v, subject, html, "html")
			if err != nil {
				log.Println(ret.Message, "----", v, "邮件发送失败！", err)
				color.New(color.BgHiRed).Println( "填写失败!!!----", v, "邮件发送失败！", err)
			}
			log.Println(ret.Message, "----", v, "邮件发送成功！！")
			color.New(color.FgHiGreen).Println( "填写失败!!!----", v, "邮件发送成功！！")
		}
		color.New(color.FgHiYellow).Println("提交失败！")
	}
	color.New(color.FgHiGreen).Println("正在运行............")
}
