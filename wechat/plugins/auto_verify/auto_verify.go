/*
 * 自动通过好友的插件,仅是自动通过好友请求,别的什么也不做
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/12
 */
package auto_verify

import (
	"../../logs" // 导入日志包
	"../../wxweb"  // 导入协议包
	//"../../../../../simplemvc/go"
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_FRIEND_REQUEST, wxweb.Handler(handle), "auto_verify")
}

// 消息处理函数
func handle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// 可选:可以用contact manager来过滤, 过滤掉没有保存到通讯录的群
	//&lt;msg fromusername="wxid_2gonc4xr40px22" encryptusername="v1_600f84cccfe971f580eff125cc3b0d66d2b7a8af91c73ab550397a1a8dab982a81c56ffb0b0041da96d6d01d4d0d98ab@stranger" fromnickname="石桥码农" content="我是石桥码农"  shortpy="SQMN" imagestatus="3" scene="14" country="CN" province="" city="" sign="日行一刻CEO，在行互联网＆篆刻行家，罗辑思维铁杆会员" percard="1" sex="1" alias="cibeiyouli" weibo="" weibonickname="" albumflag="0" albumstyle="0" albumbgimgid="" snsflag="1" snsbgimgid="http://shmmsns.qpic.cn/mmsns/INk4JvWfe8VDibkwaSrjIT0J3icJIynocvhOK2Kn5HMYudFd98Otpnvq1Uz1X3w9AEScOiaajD7pwM/0" snsbgobjectid="12534037318213841103" mhash="48039e64882e0568dee8460344e46d4a" mfullhash="48039e64882e0568dee8460344e46d4a" bigheadimgurl="http://wx.qlogo.cn/mmhead/ver_1/w5O1lVCiaYjHnxiaAFTicoz8DLLIkHnPicthvbH2SOBYJkFDemohyoAibCKIWDE8hXwbqicicbSWBDaMhgX2Jcib4ZBaicVNJHp0uErlvHgGqjaNUibsU/0" smallheadimgurl="http://wx.qlogo.cn/mmhead/ver_1/w5O1lVCiaYjHnxiaAFTicoz8DLLIkHnPicthvbH2SOBYJkFDemohyoAibCKIWDE8hXwbqicicbSWBDaMhgX2Jcib4ZBaicVNJHp0uErlvHgGqjaNUibsU/96" ticket="v2_8e8814938dd6034150c5f473d7ae31cef8e5afdce23c33807e7f5c38fea3a578dd527b17d331aea73bd5f57274b5eaa4fd5c423a0b7a63bdaa6e51245624c61c@stranger" opcode="2" googlecontact="" qrticket="" chatroomusername="7017477994@chatroom" sourceusername="" sourcenickname=""&gt;&lt;brandlist count="0" ver="655370693"&gt;&lt;/brandlist&gt;&lt;/msg&gt;

	// RecommendInfo:map[Signature:若你喜欢怪人，其实我很美 Ticket:v2_272f6da7a851a0324261ecac4ccbfbadab5a697673e3e64bf0bc69003e0554b032398a33f17fd8a64b92ae300d7619037cf31b3922fa50152934929a87ddd702@stranger QQNum:0 VerifyFlag:0 Alias: NickName:阿朱 Content:我是阿朱 OpCode:2 Scene:14 AttrStatus:2.147733861e+09 City:崇明 UserName:@3452460bb646ab2cda2c3e422777e100028df62a2e4529f90b5e53a37fd62b7b Province:上海 Sex:2]
	var rinfo = msg.OriginalMsg[`RecommendInfo`].(map[string]interface{})
	var userName = rinfo["UserName"].(string)
	var nickName = rinfo["NickName"].(string)
	var content = rinfo["Content"].(string)
	var ticket = rinfo["Ticket"].(string)
	//logs.Info("msg.OriginalMsg[Ticket]", msg.OriginalMsg["Ticket"] )

	if err := session.AcceptFriendRequest(userName,content,ticket,wxweb.FRIEND_VEFIFY_OPCODE_ACCEPT); err == nil {
		logs.Info("通过了 %s 好友请求",nickName)
	}else{
		logs.Info("通过 %s 好友请求失败 %s",nickName,err.Error())
	}

	//logs.Info( sim.ToJsonObject(msg.OriginalMsg) )
}