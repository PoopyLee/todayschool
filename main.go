package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-basic/uuid"
	"github.com/robfig/cron"
	_ "gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
	"time"
	"todayschool/browser"
	"todayschool/config"
	"todayschool/controller"
	"todayschool/untils"
)

func init() {
	c := color.New(color.FgMagenta).Add(color.Bold)
	c.Println(` ___       ___  ___       ___      ___ ___       __   _______   ___     
|\  \     |\  \|\  \     |\  \    /  /|\  \     |\  \|\  ___ \ |\  \    
\ \  \    \ \  \ \  \    \ \  \  /  / | \  \    \ \  \ \   __/|\ \  \   
 \ \  \    \ \  \ \  \    \ \  \/  / / \ \  \  __\ \  \ \  \_|/_\ \  \  
  \ \  \____\ \  \ \  \____\ \    / /   \ \  \|\__\_\  \ \  \_|\ \ \  \ 
   \ \_______\ \__\ \_______\ \__/ /     \ \____________\ \_______\ \__\
    \|_______|\|__|\|_______|\|__|/       \|____________|\|_______|\|__|
                                                                        
                                                                        
                                                                        `)

	//----------------------------------------
	// Create a new color object
	c = color.New(color.FgHiYellow)
	c.Println("如果输入错误，请查看本程序所在目录下的两个文件，分别是config.json和lilvwei.json，删除即可重新配置！")
	color.New(color.FgHiGreen).Println("程序开始运行了，等待整点运行---->每一小时运行一次")

	config.GetConfig() //初始化配置文件
}

func main() {
	SetlogFile()
	todayschool()
	StartTimeFunc(SetlogFile)
	StartTimeFunc(todayschool)
}

func todayschool() {
	toform := new(controller.SubmitForm)
	con := config.GetConfig()
	if len(con.MOD_AUTH_CAS) == 0 {
		log.Println("config.json文件中的MOD_AUTH_CAS为空,打开浏览器获取中....")
		color.New(color.BgHiRed).Println("config.json文件中的MOD_AUTH_CAS为空,打开浏览器获取中....")
		browser.OpenBrowser(con) //打开浏览器获取MOD_AUTH_CAS
		config.UpdateConfig(con) //更新配置
	} else if len(con.Address) == 0 {
		log.Println("config.json文件中的Address为空，请填写！")
		color.New(color.BgHiRed).Println("config.json文件中的Address为空，请填写！")
	}
	wid, ok := controller.PostProcessingList(con)
	if ok == 0 {
		formid, schoolid := controller.PostDetail(wid, con)
		form := controller.GetFormFileds(formid, wid, con)
		toform.SchoolTaskWid = schoolid
		toform.CollectWid = wid
		toform.FormWid = formid
		toform.Address = con.Address
		toform.UaIsCpadaily = true
		for k, v := range form.Datas.Rows {
			if k == 0 {
				str := strings.Split(v.Value, "/")
				toform.Latitude, toform.Longitude = untils.GetCityToLocation(str[1])
				if len(con.Cpdaily_Extension) == 0 {
					UpdateCpdailyConfig(str[1])
				}
			}
			toform.Form = append(toform.Form, v)
		}
		//form获取成功后开始进行组装加密信息
		fom_byte, _ := json.Marshal(toform)
		submitAes := new(controller.SubmitFormAes)
		submitAes.AppVersion = "9.0.12"
		submitAes.BodyString = base64.StdEncoding.EncodeToString(untils.EncryptAES(fom_byte))
		submitAes.Lat = toform.Latitude
		submitAes.Lon = toform.Longitude
		submitAes.CalVersion = "firstv"
		submitAes.DeviceID = untils.MD5_16(con.PhoneModel) + con.PhoneModel
		submitAes.Version = "first_v2"
		submitAes.SystemName = "android"
		submitAes.SystemVersion = "10"
		submitAes.Model = con.PhoneModel
		submitAes.UserID = con.UserID
		signs := []string{"appversion=" + submitAes.AppVersion, "bodyString=" + submitAes.BodyString, "deviceId=" + submitAes.DeviceID, "lat=" + fmt.Sprintf("%.6f", submitAes.Lat), "lon=" + fmt.Sprintf("%.6f", submitAes.Lon),
			"model=" + submitAes.Model, "systemName=" + submitAes.SystemName, "systemVersion=" + submitAes.SystemVersion, "userId=" + submitAes.UserID, "ytUQ7l2ZZu8mLvJZ"}
		submitAes.Sign = untils.MD5_32(strings.Join(signs, "&"))

		//submitAes_byte, _ := json.MarshalIndent(submitAes, "", " ")
		//fmt.Println(string(submitAes_byte))

		str := ""
		for _, v := range form.Datas.Rows {
			if v.IsRequired {
				//color.New(color.FgHiCyan).Println("必填字段：", v.Title)
				str += v.Title + ":<br>" + v.Value + "<br>"
				for _, c := range v.FieldItems {
					str += c.Content + "<br>"
				}

			}
		}
		//fmt.Println(str)
		controller.SubmitFormFileds(submitAes, str, con)
	} else {
		log.Println("无需填写...等待下一次运行")
		color.New(color.FgHiGreen).Println("无需填写...等待下一次运行")
	}
}

func UpdateCpdailyConfig(add string) {
	con := config.GetConfig()
	var confing_path = "config.json"
	jsonFile, err := os.OpenFile(confing_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Println(err, "文件打开错误，请勿更改位置")
	}
	defer jsonFile.Close()
	os.Truncate(confing_path, 0) //清空文件内容，再进行更新
	h := new(config.RootConfig)
	cpad := new(config.CpdailyExtension)
	cpad.SystemName = "android"
	cpad.SystemVersion = "10"
	cpad.Model = con.PhoneModel
	cpad.DeviceID = uuid.New()
	cpad.AppVersion = "9.0.12"
	cpad.UserID = con.UserID
	cpad.Lat, cpad.Lon = untils.GetCityToLocation(add)
	cpad_byte, _ := json.Marshal(&cpad)

	h = con
	h.Cpdaily_Extension = base64.StdEncoding.EncodeToString(untils.EncryptAES(cpad_byte)) //加密

	h_byte, _ := json.Marshal(h)
	h_byter := bytes.Buffer{}
	_ = json.Indent(&h_byter, h_byte, "", "  ")
	jsonFile.WriteString(h_byter.String())
	jsonFile.Sync()
}

func StartTimeFunc(f func()) {
	spec := "0 0 */1 * * *" //每1小时执行一次
	//spec:="* * * * * *"	//每小时执行一次
	c := cron.New()
	c.AddFunc(spec, f)
	c.Start()
	defer c.Stop()
	select {}
}

func SetlogFile() {
	//设置每日一个日志文件
	filename := time.Now().Format("2006-01-02") + ".log"
	_, err := os.Stat(filename)
	if err != nil {
		os.Create(filename)
	}
	logfile, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	log.SetOutput(logfile)
}
