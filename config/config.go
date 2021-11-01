package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

//配置类
type RootConfig struct {
	HeaderConfig `json:"Config"`
}

type HeaderConfig struct {
	BaseUrl           string   `json:"BaseUrl"`
	UserID            string   `json:"user_id"`
	Pasword           string   `json:"'pasword'"`
	MOD_AUTH_CAS      string   `json:"MOD_AUTH_CAS"`
	Address           string   `json:"Address"`
	Cpdaily_Extension string   `json:"Cpdaily-Extension"`
	UserEmail         []string `json:"UserEmail"`
	PhoneModel        string   `json:"phone_model"`
}

//CpdailyExtension类
type CpdailyExtension struct {
	SystemName    string  `json:"systemName"`
	SystemVersion string  `json:"systemVersion"`
	Model         string  `json:"model"`
	DeviceID      string  `json:"deviceId"`
	AppVersion    string  `json:"appVersion"`
	Lon           float64 `json:"lon"`
	Lat           float64 `json:"lat"`
	UserID        string  `json:"userId"`
}

var confing_path = "config.json"

var Phones = []string{
	"Redmi7",
	"SEA-AL10",
	"Redmi Note 8 Pro",
	"Google Pixel",
	"Samsung Galaxy S8",
	"OnePlus 3T",
	"Oppo A57",
	"Oppo A37",
	"HUAWEI P10",
	"HUAWEI P10 Plus",
	"Oppo A59s",
}

func init() {
	log.SetPrefix("<今日校园 author：lilvwei>")
	_, err := os.Stat(confing_path)
	if err != nil {
		os.Create(confing_path)
	}

}

func GetConfig() *RootConfig {

	jsonFile, err := os.OpenFile(confing_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Println(err, "文件打开错误，请勿更改位置")
	}
	defer jsonFile.Close()
	h := new(RootConfig)
	fd, err := ioutil.ReadFile(confing_path)
	if len(fd) == 0 {
		input := bufio.NewScanner(os.Stdin)

		h.BaseUrl = "yibinu.campusphere.net"

		rand.Seed(time.Now().Unix()) //添加随机数种子
		h.PhoneModel = Phones[rand.Intn(len(Phones))]

		color.New(color.FgGreen).Println("<---请输入学号(回车提交,仅用于登录，并不会提交给第三方)")
		input.Scan()
		h.UserID = input.Text()

		color.New(color.FgGreen).Println("<---请输入密码(回车提交，仅用于登录，并不会提交给第三方)")
		input.Scan()
		h.Pasword = input.Text()

		//color.New(color.FgGreen).Println("<---请输入MOD_AUTH_CAS(回车提交)")
		//input.Scan()
		//h.MOD_AUTH_CAS = input.Text()

		color.New(color.FgGreen).Println("<---请输入Address(回车提交，此地址是表单提交时的定位地址)")
		input.Scan()
		h.Address = input.Text()
		color.New(color.FgGreen).Println("<---请输入你的邮箱，用于接收填写成功报告！格式：1965751527@qq.com(回车提交)")
		input.Scan()
		h.UserEmail = append(h.UserEmail, input.Text())
		h_byte, _ := json.Marshal(h)
		h_byter := bytes.Buffer{}
		_ = json.Indent(&h_byter, h_byte, "", "  ")
		jsonFile.WriteString(h_byter.String())
		jsonFile.Sync()
		return h
	} else {
		json.Unmarshal(fd, &h)
		return h
	}
}

func UpdateConfig(con *RootConfig) {
	jsonFile, err := os.OpenFile(confing_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Println(err, "文件打开错误，请勿更改位置")
	}
	defer jsonFile.Close()
	os.Truncate(confing_path, 0) //清空文件内容，再进行更新
	con_byte, _ := json.Marshal(con)
	con_byter := bytes.Buffer{}
	_ = json.Indent(&con_byter, con_byte, "", "  ")
	jsonFile.WriteString(con_byter.String())
	jsonFile.Sync()
}
