package browser

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
	"log"
	"strings"
	"time"
	"todayschool/config"
)

type GetCookies struct {
	Cookies []struct {
		Name         string `json:"name"`
		Value        string `json:"value"`
		Domain       string `json:"domain"`
		Path         string `json:"path"`
		Expires      int    `json:"expires"`
		Size         int    `json:"size"`
		HTTPOnly     bool   `json:"httpOnly"`
		Secure       bool   `json:"secure"`
		Session      bool   `json:"session"`
		Priority     string `json:"priority"`
		SameParty    bool   `json:"sameParty"`
		SourceScheme string `json:"sourceScheme"`
		SourcePort   int    `json:"sourcePort"`
	} `json:"cookies"`
}

func OpenBrowser(con *config.RootConfig) {
	ctx, _ := chromedp.NewExecAllocator(
		context.Background(),

		// 以默认配置的数组为基础，覆写headless参数
		// 当然也可以根据自己的需要进行修改，这个flag是浏览器的设置
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...,
	)
	// 创建新的chromedp上下文对象，超时时间的设置不分先后
	// 注意第二个返回的参数是cancel()，只是我省略了
	ctx, _ = context.WithTimeout(ctx, 30*time.Second)
	ctx, _ = chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)

	// 执行我们自定义的任务 -函数在第4步
	if err := chromedp.Run(ctx, BrowserLogin(con)); err != nil {
		log.Fatal(err)
		return
	}
	chromedp.Cancel(ctx)
}

func BrowserLogin(con *config.RootConfig) chromedp.Tasks {
	return chromedp.Tasks{
		// 打开登陆页面
		chromedp.Navigate("http://authserver.yibinu.edu.cn/authserver/login?service=https%3A%2F%2Fyibinu.campusphere.net%2Fportal%2Flogin"),
		chromedp.SetValue("#username", con.UserID, chromedp.ByID),
		chromedp.SetValue("#password", con.Pasword, chromedp.ByID),
		chromedp.Click("#casLoginForm > p:nth-child(6) > button"),
		saveCookies(con),
	}
}

// 保存Cookies
func saveCookies(con *config.RootConfig) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		for i:=5;i>0;i-- {
			time.Sleep(time.Second*1)
			color.New(color.BgHiRed).Println("浏览器将在",i,"秒后关闭")
			log.Println("浏览器将在",i,"秒后关闭")
		}
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return
		}
		// 2. 序列化
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}
		cookiesall := new(GetCookies)
		_ = json.Unmarshal(cookiesData, &cookiesall)
		for _, v := range cookiesall.Cookies {
			if strings.Compare(v.Name, "MOD_AUTH_CAS") == 0 {
				con.MOD_AUTH_CAS = v.Value
			}
		}
		//fmt.Println(string(cookiesData))
		return
	}
}
