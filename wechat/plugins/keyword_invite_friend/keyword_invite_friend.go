/*
 * 关键字邀请入群插件
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/14
 */
package keyword_invite_friend

import (
	"../../logs" // 导入日志包
	"../../wxweb"  // 导入协议包
	//"../../../../../simplemvc/go"
	"fmt"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(handle), "keyword_invite_friend")
}

// 消息处理函数
func handle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if msg.IsGroup {
		return
	}

	for k,v := range session.Config.InviteFriendKeywords{
		//logs.Info("invite friend try",k,v)
		if msg.Content == k {
			if groups := session.Cm.GetContactByName(v); len(groups) > 0 {
				var g = groups[0]
				var friend = session.Cm.GetContactByUserName(msg.FromUserName)
				logs.Info("找到了群,开始邀请 %s 进群 %s",friend.NickName,v)

				if err := wxweb.InviteFriend(session,g,msg.FromUserName); err == nil {
					logs.Info("邀请进群成功")
					session.SendText(fmt.Sprintf(`已邀请您加入了群聊"%s"`,v), session.Bot.UserName, msg.FromUserName)
					return
				}else{
					logs.Info("邀请进群失败",err.Error())
					session.SendText("抱谦,邀请失败了", session.Bot.UserName, msg.FromUserName)
				}
			}
			break
		}
	}
}
