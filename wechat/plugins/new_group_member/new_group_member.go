/*
 * 文件暂无名,或许又是坑
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/14
 */
package new_group_member

import (
	"../../logs" // 导入日志包
	"../../wxweb"  // 导入协议包
	//"../../../../../simplemvc/go"
	"fmt"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_NEW_GROUP_MEMBER, wxweb.Handler(handle), "new_group_member")
}

// 消息处理函数
func handle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	var data = msg.Data.(wxweb.NewGroupMemberMsgdata)
	logs.Info("有新朋友进群:",data.NickName)
	//[抱拳][握手][玫瑰]
	var words = fmt.Sprintf("欢迎 @%s",data.NickName)
	var newMemberGreetingWords = ""
	if data.Group != nil {
		for k,v := range session.Config.NewGroupMemberGreetingWords {
			if k == data.Group.NickName {
				newMemberGreetingWords = v
				break
			}
		}
	}
	words = fmt.Sprintf("%s\n\n%s",words,newMemberGreetingWords)

	session.SendText(words, session.Bot.UserName, msg.FromUserName)
}