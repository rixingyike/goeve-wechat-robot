/*
 * 配置文件
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/14
 */
package wxweb

type Config struct {
	Version         int `json:"version"`//版本号
	Name string `json:"name"`//帐号名称

	EnableAutoVefiry bool `json:"enable_auto_vefiry"`//自动通过好友请求
	NewFriendGreetingWords string `json:"new_friend_greeting_words"`//新朋友问候语

	EnableWelcomeNewGroupMember bool `json:"enable_welcome_new_group_member"`//自动欢迎群内新人
	NewGroupMemberGreetingWords map[string]string `json:"new_group_member_greeting_words"`//新人进群欢迎语

	EnableKeywordInviteFriend bool `json:"enable_keyword_invite_friend"`//启用关键字邀请入群
	InviteFriendKeywords map[string]string `json:"invite_friend_keywords"`//关键字入群
}
