/*
 * 群人新人进群消息数据
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/14
 */
package wxweb

type NewGroupMemberMsgdata struct {
	NickName string
	Group *User
}
