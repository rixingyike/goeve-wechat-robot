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

func main() {
	// create session
	session, err := wxweb.CreateSession(nil, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	// load plugins for this session

	if session.Config.EnableAutoVefiry {
		auto_verify.Register(session)
		session.HandlerRegister.EnableByName("auto_verify")
		new_friend.Register(session)
		session.HandlerRegister.EnableByName("new_friend")
	}
	if session.Config.EnableWelcomeNewGroupMember{
		new_group_member.Register(session)
		session.HandlerRegister.EnableByName("new_group_member")
	}
	if session.Config.EnableKeywordInviteFriend{
		keyword_invite_friend.Register(session)
		session.HandlerRegister.EnableByName("keyword_invite_friend")
	}

	for {
		if err := session.LoginAndServe(true); err != nil {
			logs.Error("session exit, %s", err)
			logs.Info("trying re-login with qrcode")
			if err := session.LoginAndServe(false); err != nil {
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