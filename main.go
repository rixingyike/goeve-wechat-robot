package main

import (
	"./wechat/logs"
	"./wechat/wxweb"
	"time"
	"./wechat/plugins/auto_verify"
	"github.com/andlabs/ui"
)

var SLEEP_TIME = 5
var eveRobot *wxweb.Session
var loop = true

func main() {
	uiMain()
	//StartRobot()
}

func uiMain()  {
	err := ui.Main(func() {
		button := ui.NewButton("Start")
		greeting := ui.NewLabel("感谢使用伊夫")
		box := ui.NewVerticalBox()
		box.Append(greeting, false)
		box.Append(button, false)

		window := ui.NewWindow("日行一刻伊夫微智助理", 300, 200, false)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			if button.Text() == "Start" {
				greeting.SetText("伊夫运行中..")
				button.SetText("Stop")
				loop = true
				go StartRobot()
			}else{
				greeting.SetText("伊夫已停止")
				button.SetText("Start")
				loop = false
				//eveRobot.Logout()
			}
			//greeting.SetText("Hello, " + name.Text() + "!")
		})
		window.OnClosing(func(*ui.Window) bool {
			eveRobot.Logout()
			loop = false
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		//panic(err)
	}
}

func StartRobot()  {
	// create session
	eveRobot, err := wxweb.CreateSession(nil, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	// load plugins for this session

	if eveRobot.Config.EnableAutoVefiry {
		auto_verify.Register(eveRobot)
		eveRobot.HandlerRegister.EnableByName("auto_verify")
	}

	for {
		if !loop {
			break
		}
		if err := eveRobot.LoginAndServe(true); err != nil {
			logs.Error("session exit, %s", err)
			logs.Info("trying re-login with qrcode")
			if err := eveRobot.LoginAndServe(false); err != nil {
				logs.Error("re-login error, %s", err)
			}
			time.Sleep(time.Duration(SLEEP_TIME) * time.Second)
			SLEEP_TIME += SLEEP_TIME
		} else {
			logs.Info("closed by user")
			break
		}
	}
}