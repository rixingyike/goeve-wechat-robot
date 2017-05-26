package main

import (
	"./wechat/logs"
	"./wechat/wxweb"
	"time"
	"./wechat/plugins/auto_verify"
	"./wechat/plugins/new_friend"
	"./wechat/plugins/new_group_member"
	"./wechat/plugins/keyword_invite_friend"
)

var SLEEP_TIME = 5
var eveRobot *wxweb.Session

func main()  {
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
		new_friend.Register(eveRobot)
		eveRobot.HandlerRegister.EnableByName("new_friend")
	}
	if eveRobot.Config.EnableWelcomeNewGroupMember{
		new_group_member.Register(eveRobot)
		eveRobot.HandlerRegister.EnableByName("new_group_member")
	}
	if eveRobot.Config.EnableKeywordInviteFriend{
		keyword_invite_friend.Register(eveRobot)
		eveRobot.HandlerRegister.EnableByName("keyword_invite_friend")
	}

	for {
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