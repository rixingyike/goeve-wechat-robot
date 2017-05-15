/*
 * 文件暂无名,或许又是坑
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/14
 */
package new_friend

import (
	"../../logs" // 导入日志包
	"../../wxweb"  // 导入协议包
	//"../../../../../simplemvc/go"
	"fmt"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_NEW_FRIEND, wxweb.Handler(handle), "new_friend")
}

// 消息处理函数
func handle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	var u = msg.Data.(wxweb.User)
	logs.Info("新增了好友",u.NickName)
	var sex = u.Sex
	//[抱拳][握手][玫瑰]
	var words = "[握手]"
	if sex == wxweb.USER_SEX_MALE {
		words = "[抱拳]"
	}else if sex == wxweb.USER_SEX_FEMALE {
		words = "[玫瑰]"
	}
	words = fmt.Sprintf("%s\n\n%s", words, session.Config.NewFriendGreetingWords)
	session.SendText(words, session.Bot.UserName, u.UserName)
}
