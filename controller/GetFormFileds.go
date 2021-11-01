package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"todayschool/config"
)

type Forms struct {
	Wid           string        `json:"wid"`
	FormWid       string        `json:"formWid"`
	FieldType     string        `json:"fieldType"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	IsRequired    bool          `json:"isRequired"`
	HasOtherItems bool          `json:"hasOtherItems"`
	Sort          int           `json:"sort"`
	ColName       string        `json:"colName"`
	FieldItems    [] FieldItems `json:"fieldItems"`
	ScoringRule   string        `json:"scoringRule"`
	Score         interface{}   `json:"score"`
	AnswerContent interface{}   `json:"answerContent"`
	BasicConfig   struct {
		MinValue       interface{} `json:"minValue"`
		MaxValue       interface{} `json:"maxValue"`
		Decimals       interface{} `json:"decimals"`
		DateFormatType interface{} `json:"dateFormatType"`
		DateType       interface{} `json:"dateType"`
		PointTime      interface{} `json:"pointTime"`
		MinTime        interface{} `json:"minTime"`
		MaxTime        interface{} `json:"maxTime"`
		ExtremeLabel   interface{} `json:"extremeLabel"`
		Level          interface{} `json:"level"`
		IsReversal     interface{} `json:"isReversal"`
		LeftValue      interface{} `json:"leftValue"`
		RightValue     interface{} `json:"rightValue"`
		VoteType       interface{} `json:"voteType"`
		VoteMinValue   interface{} `json:"voteMinValue"`
		VoteMaxValue   interface{} `json:"voteMaxValue"`
	} `json:"basicConfig"`
	LogicWid        int         `json:"logicWid"`
	Value           string      `json:"value"`
	Show            interface{} `json:"show"`
	FormType        string      `json:"formType"`
	SortNum         string      `json:"sortNum"`
	LogicShowConfig struct {
	} `json:"logicShowConfig"`
}

type FieldItems struct {
	ItemWid       string      `json:"itemWid"`
	Content       string      `json:"content"`
	ImageURL      interface{} `json:"imageUrl"`
	IsOtherItems  bool        `json:"isOtherItems"`
	ContentExtend string      `json:"contentExtend"`
	OtherItemType string      `json:"otherItemType"`
	BasicConfig   interface{} `json:"basicConfig"`
	ShowLogic     string      `json:"showLogic"`
	IsAnswer      bool        `json:"isAnswer"`
	Score         interface{} `json:"score"`
	IsSelected    interface{} `json:"isSelected"`
	SelectCount   interface{} `json:"selectCount"`
	TotalCount    int         `json:"totalCount"`
}

type GetForm struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		TotalSize int     `json:"totalSize"`
		ExistData int     `json:"existData"`
		Rows      []Forms `json:"rows"`
	} `json:"datas"`
}

type ToForm struct {
	PageNumber   int    `json:"pageNumber"`
	PageSize     int    `json:"pageSize"`
	FormWid      string `json:"formWid"`
	CollectorWid string `json:"collectorWid"`
}

var base_formfiledscontroller = "/wec-counselor-collector-apps/stu/collector/getFormFields"

func GetFormFileds(formwid, collectorwid string, con *config.RootConfig) *GetForm {
	str := bytes.Buffer{}
	send_form := ToForm{PageNumber: 1, PageSize: 999, FormWid: formwid, CollectorWid: collectorwid}
	jsons, _ := json.Marshal(send_form)
	post_controller := "https://" + con.BaseUrl + base_formfiledscontroller
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
	list := new(GetForm)
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Indent(&str, body, "", "   ")
	json.Unmarshal(body, &list)
	//fmt.Println(str.String())

	_, err = os.Stat("lilvwei.json")
	if err != nil {
		os.Create("lilvwei.json")
	}

	configfile, err := os.OpenFile("lilvwei.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Println("文件创建失败", err)
	}
	defer configfile.Close()
	fd, err := ioutil.ReadFile("lilvwei.json")
	if err != nil {
		log.Println("读取失败!", err)
		return nil
	}

	if len(fd) == 0 {
		color.New(color.FgHiYellow).Add(color.Underline).Println("没有内容开始写入配置......")
		var write bytes.Buffer
		var write_byte []byte
		newform := new(GetForm)
		input := bufio.NewScanner(os.Stdin) //标准输入
		for k, v := range list.Datas.Rows {
			if v.IsRequired {
				if k == 0 {
					color.New(color.FgHiCyan).Println("必填字段(格式:四川省/成都市/郫都区)：", v.Title)
				} else {
					color.New(color.FgHiCyan).Println("必填字段：", v.Title)
				}
				color.New(color.FgHiCyan).Println("内容选择(只需输入序号即可！)：")
				for _, b := range v.FieldItems {
					color.New(color.FgHiCyan).Println("      ", b.ItemWid, "<---------->", b.Content)
				}
				if v.BasicConfig.MaxValue == nil && v.BasicConfig.MinValue == nil {
					fmt.Println("请输入(回车提交)：")
				} else {
					fmt.Println("请输入(回车提交)最大值:", v.BasicConfig.MaxValue, ",最小值:", v.BasicConfig.MinValue, ":")
				}
				input.Scan()
				list.Datas.Rows[k].Value = input.Text()
				for f, c := range v.FieldItems {
					if strings.Compare(input.Text(), c.ItemWid) == 0 {
						color.New(color.FgHiGreen).Println("你选择了：", c.Content)
						list.Datas.Rows[k].FieldItems[f].ItemWid = input.Text()
						list.Datas.Rows[k].FieldItems[f].IsSelected = 1
						list.Datas.Rows[k].Value = ""
						list.Datas.Rows[k].FieldItems = list.Datas.Rows[k].FieldItems[f : f+1]
					}
				}
			} else {
				list.Datas.Rows[k].FieldItems = list.Datas.Rows[k].FieldItems[0:0] //多余数组元素切掉
			}
			newform.Datas.Rows = append(newform.Datas.Rows, list.Datas.Rows[k])
			newform.Datas.Rows[k].Show = true
			newform.Datas.Rows[k].FormType = "0"
			newform.Datas.Rows[k].SortNum = strconv.Itoa(k + 1)
		}
		newform.Message = list.Message
		newform.Code = list.Code
		write_byte, _ = json.Marshal(newform)
		_ = json.Indent(&write, write_byte, "", "  ") //格式化写入
		configfile.WriteString(write.String())
		configfile.Sync()
		return newform
	} else {
		ready_list := new(GetForm)
		json.Unmarshal(fd, &ready_list)
		for k, v := range list.Datas.Rows {
			ready_list.Datas.Rows[k].FormWid = v.FormWid
			ready_list.Datas.Rows[k].Wid = v.Wid
			for a, b := range v.FieldItems {
				for e, _ := range ready_list.Datas.Rows[k].FieldItems {
					if strings.Compare(ready_list.Datas.Rows[k].FieldItems[e].Content, list.Datas.Rows[k].FieldItems[a].Content) == 0 {
						ready_list.Datas.Rows[k].FieldItems[e].ItemWid = b.ItemWid
					}
				}
			}
		}
		return ready_list
	}
}
