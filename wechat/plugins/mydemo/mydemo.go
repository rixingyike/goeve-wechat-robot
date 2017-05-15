/*
 * 这是一个插件的示例
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/12
 */
package mydemo

import (
	"../../logs" // 导入日志包
	"../../wxweb"  // 导入协议包
	"strings"
	"regexp"
	"fmt"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(handle), "mydemo")

	// 可以多个消息类型使用同一个处理函数，也可以分开
	//session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(demo), "imgdemo")
}

// 消息处理函数
func handle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	// 可选:可以用contact manager来过滤, 过滤掉没有保存到通讯录的群
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("ignore the messages from", msg.FromUserName)
		return
	}

	// 可选: 过滤消息类型
	//if msg.MsgType == wxweb.MSG_IMG {
	//	return
	//}

	// 群消息
	if msg.IsGroup {
		// 艾特自己主动加好友
		if msg.IsAtMe {
			//如果不是好友,查不到好友
			var friend = session.Cm.GetContactByUserName(msg.Who)
			var words string
			if friend != nil {
				words = fmt.Sprintf("@%s 私聊~",friend.NickName)
			}else{
				// 如果不是好友,通过群组对象查找联系人信息
				var memebeNicName string
				if group := session.Cm.GetContactByUserName(msg.FromUserName);group != nil {
					if mm,err := wxweb.CreateMemberManagerFromGroupContact(session,group); err == nil {
						if member := mm.GetContactByUserName(msg.Who);member != nil{
							memebeNicName = member.NickName
						}
					}
				}

				if memebeNicName != "" {
					words = fmt.Sprintf("@%s 加你了,私聊~",memebeNicName)
				}else{
					words = fmt.Sprintf("加你了,私聊~")
				}

				var group = session.Cm.GetContactByUserName(msg.FromUserName);
				var verifyContent = fmt.Sprintf(`我是群聊"%s"的%s`,group.NickName,session.Bot.NickName)
				if err := session.MakeFriend(msg.Who,verifyContent); err == nil {
					logs.Info("向%s发送好友请求成功",msg.Who)
				}else{
					logs.Info("向%s发送好友请求失败",msg.Who)
				}
			}
			session.SendText(words, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
			return
		}
		// 皮球踢人的技法
		if strings.HasPrefix(msg.Content,"@") {
			var kickReg = regexp.MustCompile(`^@(.+)? \[足球\]`)//[足球]
			if names := kickReg.FindStringSubmatch(msg.Content); len(names) > 1 {
				var memberDisplayName = names[1]
				logs.Info("开始踢人:", memberDisplayName)
				if group := session.Cm.GetContactByUserName(msg.FromUserName);group != nil {
					if mm,err := wxweb.CreateMemberManagerFromGroupContact(session,group); err == nil {
						if members := mm.GetContactsByDisplayName(memberDisplayName); len(members)>0{
							var member = members[0]
							if err := wxweb.KickGroupMember(session,group,member.UserName); err == nil {
								var words = fmt.Sprintf("已将 @%s 移出群聊", memberDisplayName)
								session.SendText(words, session.Bot.UserName, msg.FromUserName)
								logs.Info("从群聊 %s 中移除 %s 成功",group.NickName, memberDisplayName)
								return
							}else{
								logs.Info("踢人失败",err.Error())
							}
						}else{
							logs.Info("未找到目标")
						}
					}
				}
				logs.Info("从群聊中移除 %s 失败", memberDisplayName)
			}
		}


		return
	}

	// 取出收到的内容
	// 取text
	logs.Info(msg.Content)
	//// 取img
	//if b, err := session.GetImg(msg.MsgId); err == nil {
	//	logs.Debug(string(b))
	//}

	// anything

	var txt = msg.Content;

	if txt == "查看群组" {
		groups := session.Cm.GetGroupContact()
		var message = ""
		for _, g := range groups {
			message += "\n" + g.DisplayName + "|" + g.NickName
		}
		session.SendText(message, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
	} else if strings.HasPrefix(txt,"查看群") {
		if txt == "查看群" {
			groups := session.Cm.GetContactByName("「河南」日行一刻修行公社")
			if len(groups) == 0 {
				logs.Info("not group found")
				return
			}
			group := groups[0]
			if mm,err := wxweb.CreateMemberManagerFromGroupContact(session,group); err == nil {
				logs.Info("mm group name %s", mm.Group.DisplayName)
				members := mm.Group.MemberList
				logs.Info("group members count %d", len(members))

				var message = ""
				for _, m := range members {
					message += "\n" + m.NickName + "|" + m.RemarkName + "|"+m.DisplayName
				}
				session.SendText(message, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
			}
		}else{
			var viewGroupReg = regexp.MustCompile(`^查看群 (.+)`)
			if subs := viewGroupReg.FindStringSubmatch(msg.Content); len(subs) > 1 {
				var friendName = subs[1]
				if groups := session.Cm.GetContactByName(friendName); len(groups) > 0 {
					var g = groups[0]
					var message = "\n" + g.DisplayName + "|" + g.NickName + "|" + g.RemarkName+ "|" + g.UserName
					session.SendText(message, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
				}
			}
		}
	} else if strings.HasPrefix(txt,"查看好友") {
		if txt == "查看好友" {
			friends := session.Cm.GetStrangers()
			var message = ""
			for _, f := range friends {
				message += "\n" + f.DisplayName + "|" + f.NickName
			}
			session.SendText(message, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		}else{
			var viewFriendReg = regexp.MustCompile(`^查看好友 (.+)`)
			if subs := viewFriendReg.FindStringSubmatch(msg.Content); len(subs) > 1 {
				var friendName = subs[1]
				if friends := session.Cm.GetContactByName(friendName); len(friends) > 0 {
					var f = friends[0]
					var message = "\n" + f.DisplayName + "|" + f.NickName + "|" + f.RemarkName+ "|" + f.UserName
					session.SendText(message, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
				}
			}
		}

	} else if strings.HasPrefix(txt,"邀请进群") {
		if txt != "邀请进群" {
			var inviteReg = regexp.MustCompile(`^邀请进群 (.+)`)
			if subs := inviteReg.FindStringSubmatch(msg.Content); len(subs) > 1 {
				var groupName = subs[1]
				if groups := session.Cm.GetContactByName(groupName); len(groups) > 0 {
					var g = groups[0]

					var friend = session.Cm.GetContactByUserName(msg.FromUserName)
					logs.Info("找到了群,开始邀请 %s 进群 %s",friend.NickName,groupName)

					if err := wxweb.InviteFriend(session,g,msg.FromUserName); err == nil {
						logs.Info("邀请进群成功")
					}else{
						logs.Info("邀请进群失败",err.Error())
					}
				}
			}
		}

	}else {
		// 回复消息
		// 第一个参数: 回复的内容
		// 第二个参数: 机器人ID
		// 第三个参数: 联系人/群组/特殊账号ID
		//session.SendText("plugin demo", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		// 回复图片和gif 参见wxweb/session.go
	}

}