/*
 * 群内管理员踢人插件
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/14
 */
package kick_member

import (
	"../../logs" // 导入日志包
	"../../wxweb"  // 导入协议包
	//"../../../../../simplemvc/go"
	"regexp"
	"strings"
	"fmt"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(handle), "kick_member")
}

var adminReg = regexp.MustCompile(`(管理员)+`)
var kickReg = regexp.MustCompile(`^@(.+)\S+\[菜刀\]`)//[足球]

// 消息处理函数
func handle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if !msg.IsGroup || msg.IsAtMe {
		return
	}

	if strings.HasPrefix(msg.Content,"@") {
		// 管理员必须是好友,且备注为管理员字样
		if !msg.IsSendedByMySelf{//机器人自己有权限,此处标注管理员的也有权限
			if admin := session.Cm.GetContactByUserName(msg.Who); admin == nil {
				logs.Info("通讯录未找到管理员")
				return
			}else if !adminReg.MatchString(admin.RemarkName) {
				logs.Info("%s没有管理员权限",admin.NickName)
				var words = fmt.Sprintf("@%s Sorry,你没有管理员权限",admin.DisplayName)
				session.SendText(words, session.Bot.UserName, msg.FromUserName)
				return
			}
		}

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


}
