package wxweb

import (
	"fmt"
	"github.com/mdp/qrterminal"
	"../config"
	"../logs"
	"../storage"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"encoding/json"
	"net/http/cookiejar"
	"crypto/tls"
	"net"
	"regexp"
	"github.com/skratchdot/open-golang/open"
)

const (
	// WEB_MODE: in this mode CreateSession will return a QRCode image url
	WEB_MODE = iota + 1
	// MINAL_MODE:  CreateSession will output qrcode in terminal
	TERMINAL_MODE
)

var (
	// DefaultCommon: default session config
	DefaultCommon = &Common{
		AppId:      "wx782c26e4c19acffb",
		LoginUrl:   "https://login.weixin.qq.com",
		Lang:       "zh_CN",
		DeviceID:   "e" + GetRandomStringFromNum(15),
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		SyncSrv:    "webpush.wx.qq.com",
		UploadUrl:  "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
		MediaCount: 0,
	}
)

// Session: wechat bot session
type Session struct {
	Client          *http.Client//共享一个请求供,貌似没有多线程问题

	WxWebCommon     *Common
	WxWebXcg        *XmlConfig
	Cookies         []*http.Cookie
	SynKeyList      *SyncKeyList
	Bot             *User
	Cm              *ContactManager
	QrcodePath      string //qrcode path
	QrcodeUUID      string //uuid
	HandlerRegister *HandlerRegister

	Cache *Cache
	Config Config//配置信息
}

// CreateSession: create wechat bot session
// if common is nil, session will be created with default config
// if handlerRegister is nil,  session will create a new HandlerRegister
func CreateSession(common *Common, handlerRegister *HandlerRegister) (*Session, error) {
	if common == nil {
		common = DefaultCommon
	}

	client, err := newClient()
	if err != nil {
		return nil, err
	}

	// get qrcode
	uuid, err := JsLogin(client, common)
	if err != nil {
		return nil, err
	}
	//logs.Info(uuid)

	session := &Session{
		Client:client,
		WxWebCommon: common,
		QrcodeUUID:  uuid,
		WxWebXcg:&XmlConfig{},
		Cm:&ContactManager{},
	}
	session.Cache = &Cache{session:session}
	session.LoadConfig()

	if handlerRegister != nil {
		session.HandlerRegister = handlerRegister
	} else {
		session.HandlerRegister = CreateHandlerRegister()
	}


	return session, nil
}

func newClient() (*http.Client, error) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	transport := http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Minute,
		}).Dial,
		TLSHandshakeTimeout: 1 * time.Minute,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: &transport,
		Jar:       jar,
		Timeout:   1 * time.Minute,
	}

	return client, nil
}

func (s *Session) ShowQrcode(qrmode int) error {
	if qrmode == TERMINAL_MODE {
		qrterminal.Generate("https://login.weixin.qq.com/l/"+s.QrcodeUUID, qrterminal.L, os.Stdout)
	} else if qrmode == WEB_MODE {
		qrcb, err := QrCode(s.Client, s.WxWebCommon, s.QrcodeUUID)
		if err != nil {
			return err
		}
		ls := rrstorage.CreateLocalDiskStorage("./public/")
		if err := ls.Save(qrcb, "qrcode.jpg"); err != nil {
			return err
		}
		s.QrcodePath = "./public/qrcode.jpg"
		logs.Info("qrcode path", s.QrcodePath)
		open.Start(s.QrcodePath)
	}
	return nil
}

func (s *Session) analizeVersion(uri string) {
	u, _ := url.Parse(uri)

	// version may change
	s.WxWebCommon.CgiDomain = u.Scheme + "://" + u.Host
	s.WxWebCommon.CgiUrl = s.WxWebCommon.CgiDomain + "/cgi-bin/mmwebwx-bin"

	if strings.Contains(u.Host, "wx2") {
		// new version
		s.WxWebCommon.SyncSrv = "webpush.wx2.qq.com"
	} else {
		// old version
		s.WxWebCommon.SyncSrv = "webpush.wx.qq.com"
	}
}

func (s *Session) scanWaiter() error {
loop1:
	for {
		select {
		case <-time.After(3 * time.Second):
			redirectUri, err := Login(s.Client, s.WxWebCommon, s.QrcodeUUID, "0")
			if err != nil {
				logs.Error(err)
				if strings.Contains(err.Error(), "window.code=408") {
					return err
				}
			} else {
				s.WxWebCommon.RedirectUri = redirectUri
				s.SetCookiesAfterScanQrcode()
				break loop1
			}
		}
	}
	return nil
}

// 处理成功扫码之后的
func (s *Session) SetCookiesAfterScanQrcode() {
	s.analizeVersion(s.WxWebCommon.RedirectUri)
	var u,_ = url.Parse(s.WxWebCommon.CgiUrl)
	s.Client.Jar.SetCookies(u,s.Cookies)
}

// @param userCache 代表不使用接口登陆,尝试使用cookie缓存,是第二次+登陆时使用
// LoginAndServe: login wechat web and enter message receiving loop
func (s *Session) LoginAndServe(useCache bool) error {
	var err error

	logs.Info("useCache:",useCache)
	if !useCache || s.Cache.Load() != nil {
		//显示终端二维码
		s.ShowQrcode(WEB_MODE)

		if err := s.scanWaiter(); err != nil {
			return err
		}
		if s.Cookies, err = WebNewLoginPage(s.Client, s.WxWebCommon, s.WxWebXcg, s.WxWebCommon.RedirectUri); err != nil {
			return err
		}
	}

	jb, err := WebWxInit(s.Client, s.WxWebCommon, s.WxWebXcg)
	if err != nil {
		return err
	}

	jc, err := rrconfig.LoadJsonConfigFromBytes(jb)
	if err != nil {
		return err
	}

	s.SynKeyList, err = GetSyncKeyListFromJc(jc)
	if err != nil {
		return err
	}
	s.Bot, _ = GetUserInfoFromJc(jc)
	//logs.Info(s.Bot)
	ret, err := WebWxStatusNotify(s.Client, s.WxWebCommon, s.WxWebXcg, s.Bot)
	if err != nil {
		return err
	}
	if ret != 0 {
		return fmt.Errorf("WebWxStatusNotify fail, %d", ret)
	}

	if err := s.RetrieveAllContact(); err != nil {
		// 获取联系人失败,为0,删除cookie缓存
		s.DeleteCookieCache()
		return err
	}

	s.Cache.Write()

	if err := s.serve(); err != nil {
		return err
	}
	return nil
}

//删除缓存使用二维码登陆
func (this *Session) DeleteCookieCache(){
	os.Remove("./cache.json")
}

// 循环拉取联系人,直到seq=0
func (s *Session) RetrieveAllContact() error {
	var seq = 0
	for {
		cb, err := WebWxGetContact(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, seq)
		if err != nil {
			return err
		}
		var cr ContactResponse
		if err := json.Unmarshal(cb, &cr); err != nil {
			return err
		}
		// 处理user type
		for _, u := range cr.MemberList {
			if u.VerifyFlag/8 != 0 {
				u.Type = USER_TYPE_OFFICAL
			} else if strings.HasPrefix(u.UserName, `@@`) {
				u.Type = USER_TYPE_GROUP
			} else {
				u.Type = USER_TYPE_FRIEND
			}
		}

		if s.Cm.cl == nil {
			s.Cm.cl = cr.MemberList
		}else{
			s.Cm.cl = append(s.Cm.cl, cr.MemberList...)
		}

		logs.Warn("已获取一批联系人:",len(s.Cm.cl))
		seq = cr.Seq
		if seq == 0 {
			break
		}
	}
	var n = len(s.Cm.cl)
	logs.Warn("已获取全部联系人:",n)
	if n == 0 {
		return fmt.Errorf("get all contact error")
	}

	return nil
}

func (s *Session) serve() error {
	msg := make(chan []byte, 1000)
	// syncheck
	errChan := make(chan error)
	go s.producer(msg, errChan)
	for {
		select {
		case m := <-msg:
			go s.consumer(m)
		case err := <-errChan:
			// all received message have been consumed
			return err
		}
	}
}
func (s *Session) producer(msg chan []byte, errChan chan error) {
	logs.Info("entering synccheck loop")
loop1:
	for {
		ret, sel, err := SyncCheck(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, s.WxWebCommon.SyncSrv, s.SynKeyList)
		//[I] webpush.wx.qq.com 0 6
		logs.Info(s.WxWebCommon.SyncSrv, ret, sel)

		if err != nil {
			logs.Error(err)
			continue
		}
		if ret == 0 {
			// check success
			err := WebWxSync(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, msg, s.SynKeyList)
			if err != nil {
				logs.Error(err)
			}
			if sel == 2 {
				// new message
				//err := WebWxSync(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, msg, s.SynKeyList)

			} else if sel != 0 && sel != 7 {
				// 此处仅是取了消息变化
				//errChan <- fmt.Errorf("session down, sel %d", sel)
				//break loop1
			}
		}else if ret == 1100 {
			errChan <- fmt.Errorf("失败/登出微信:%d",ret)
			break loop1
		} else if ret == 1101 {
			errChan <- nil
			break loop1
		} else if ret == 1205 {
			errChan <- fmt.Errorf("api blocked, ret:%d", 1205)
			break loop1
		}
	}

}

// 解析同步到的数据
func (s *Session) consumer(msg []byte) {
	// analize message
	jc, _ := rrconfig.LoadJsonConfigFromBytes(msg)
	//logs.Info("consumer.AddMsgCount..",jc)

	//AddMsgCount,ModContactCount,DelContactCount,ModChatRoomMemberCount
	var ModContactCount,_ = jc.GetInt("ModContactCount")//通过了好友请求+1,关注公众号+1,删除了好友+1?
	var DelContactCount,_ = jc.GetInt("DelContactCount")//取关公众号+1,
	var ModChatRoomMemberCount,_ = jc.GetInt("ModChatRoomMemberCount")
	var AddMsgCount, _ = jc.GetInt("AddMsgCount")

	logs.Info("AddMsgCount:",AddMsgCount,"ModContactCount:",ModContactCount,"DelContactCount:",DelContactCount,"ModChatRoomMemberCount:",ModChatRoomMemberCount)

	// 联系人变化
	if ModContactCount > 0 {
		var ModContactList,_ = jc.GetInterfaceSlice("ModContactList")
		//var mcts []map[string]interface{}
		for _, v0 := range ModContactList {
			v := v0.(map[string]interface{})
			//var isNewFriend bool //是否新增好友
			vf, _ := v[`VerifyFlag`].(float64)
			un, _ := v[`UserName`].(string)

			if vf/8 != 0 {
				v[`Type`] = USER_TYPE_OFFICAL
				//mcts = append(mcts, v)
				s.Cm.UpdateContact(v)
			} else if strings.HasPrefix(un, `@@`) {
				v[`Type`] = USER_TYPE_GROUP
				//群组的联系人变化,会走到这里,目前群组成员是临时获取的,所以这里不用处理

			} else {
				v[`Type`] = USER_TYPE_FRIEND
				//mcts = append(mcts, v)
				//logs.Warn("jc",jc)
				// 为什么自己删除好友也会走到这里,走的是修改,并不是删除
				if u,isNewFriend,err := s.Cm.UpdateContact(v); err == nil {//如果是新增用户,派发一个消息
					if isNewFriend {
						var rmsg = &ReceivedMessage{
							MsgType:MSG_NEW_FRIEND,
							Data:u,
						}
						s.HandlerRegister.Runs(s,rmsg)
					}
				}
			}
		}
	}

	// 删除了联系人
	if DelContactCount > 0 {
		var DelContactList,_ = jc.GetInterfaceSlice("DelContactList")
		var mcts []map[string]interface{}
		for _, v0 := range DelContactList {
			v := v0.(map[string]interface{})
			mcts = append(mcts, v)
		}
		s.Cm.RemoveContact(mcts)
	}

	// 有新消息
	// 新消息随着联系人变化现时出现,所以将消息处理放在最后
	if AddMsgCount > 0 {
		msgis, _ := jc.GetInterfaceSlice("AddMsgList")
		for _, v := range msgis {
			rmsg := s.analize(v.(map[string]interface{}))
			//logs.Info("rmsg.MsgType msg:",rmsg.MsgType,rmsg)
			go s.PreparseMessage(rmsg)
			s.HandlerRegister.Runs(s,rmsg)
		}
	}
}

// 解析新增的消息并生成一个对象
func (s *Session) analize(msg map[string]interface{}) *ReceivedMessage {
	rmsg := &ReceivedMessage{
		MsgId:        msg["MsgId"].(string),
		Content:      msg["Content"].(string),
		FromUserName: msg["FromUserName"].(string),
		ToUserName:   msg["ToUserName"].(string),
		MsgType:      int(msg["MsgType"].(float64)),
		OriginalMsg:msg,
	}

	if strings.Contains(rmsg.FromUserName, "@@") {
		rmsg.IsGroup = true
		// group message
		ss := strings.Split(rmsg.Content, ":")
		if len(ss) > 1 {
			rmsg.Who = ss[0]
			rmsg.Content = strings.TrimPrefix(ss[1], "<br/>")
			rmsg.IsSendedByMySelf = rmsg.Who == s.Bot.UserName
		}

		// 检查是否艾特了自己
		if !rmsg.IsSendedByMySelf{
			atme := `@`
			if len(s.Bot.DisplayName) > 0 {
				atme += s.Bot.DisplayName
			} else {
				atme += s.Bot.NickName
			}
			rmsg.IsAtMe = strings.Contains(rmsg.Content, atme)
		}
	}
	return rmsg

}

// SendText: send text msg type 1
func (s *Session) SendText(msg, from, to string) (string, string, error) {
	b, err := WebWxSendMsg(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, from, to, msg)
	if err != nil {
		return "", "", err
	}
	jc, _ := rrconfig.LoadJsonConfigFromBytes(b)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	if ret != 0 {
		errMsg, _ := jc.GetString("BaseResponse.ErrMsg")
		return "", "", fmt.Errorf("WebWxSendMsg Ret=%d, ErrMsg=%s", ret, errMsg)
	}
	msgID, _ := jc.GetString("MsgID")
	localID, _ := jc.GetString("LocalID")
	// 每发过一次消息,暂停200毫秒
	time.Sleep(time.Millisecond * time.Duration(200))
	return msgID, localID, nil
}

// SendImg: send img, upload then send
func (s *Session) SendImg(path, from, to string) {
	ss := strings.Split(path, "/")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return
	}
	mediaId, err := WebWxUploadMedia(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := WebWxSendMsgImg(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
		return
	}
}

// SendImgFromBytes: send image from mem
func (s *Session) SendImgFromBytes(b []byte, path, from, to string) {
	ss := strings.Split(path, "/")
	mediaId, err := WebWxUploadMedia(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := WebWxSendMsgImg(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
		return
	}
}

// GetImg: get img by MsgId
func (s *Session) GetImg(msgId string) ([]byte, error) {
	return WebWxGetMsgImg(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, msgId)
}

// SendEmotionFromPath: send gif, upload then send
func (s *Session) SendEmotionFromPath(path, from, to string) {
	ss := strings.Split(path, "/")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return
	}
	mediaId, err := WebWxUploadMedia(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := WebWxSendEmoticon(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
	}
}

// SendEmotionFromBytes: send gif/emoji from mem
func (s *Session) SendEmotionFromBytes(b []byte, from, to string) {
	mediaId, err := WebWxUploadMedia(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, from+".gif", b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := WebWxSendEmoticon(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
	}
}

// RevokeMsg: revoke message
func (s *Session) RevokeMsg(clientMsgId, svrMsgId, toUserName string) {
	err := WebWxRevokeMsg(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies, clientMsgId, svrMsgId, toUserName)
	if err != nil {
		logs.Error("revoke msg %s failed, %s", clientMsgId+":"+svrMsgId, err)
		return
	}
}

// Logout: logout web wechat

func (s *Session) Logout() error {
	return WebWxLogout(s.Client, s.WxWebCommon, s.WxWebXcg, s.Cookies)
}

// 新增方法

// 接受好友请求
// @param 3=通过好友请求
func (s *Session) AcceptFriendRequest(userName,content,ticket string, opcode int) (err error) {
	//logs.Info("AcceptFriendRequest ticket ",ticket)
	var vefifyUsers = []*VerifyUser{
		&VerifyUser{
			Value:userName,
			VerifyUserTicket:ticket,
		},
	}
	if body,err := WebWxVerifyUser(s.Client, s.WxWebCommon,s.WxWebXcg,s.Cookies,content,vefifyUsers, opcode); err == nil {
		var res BaseResponse
		if err = json.Unmarshal(body, &res); err != nil {
			return err
		}
		if res.Ret != 0 {
			return fmt.Errorf("accept friend request err")
		}
	}else{
		logs.Info("AcceptFriendRequest.err",err.Error())
	}
	return nil
}

// 主动添加好友
// @param 2=加好友
func (s *Session) MakeFriend(userName,verifyContent string) (err error) {
	//logs.Info("AcceptFriendRequest ticket ",ticket)
	var vefifyUsers = []*VerifyUser{
		&VerifyUser{
			Value:userName,
			VerifyUserTicket:"",
		},
	}
	if body,err := WebWxVerifyUser(s.Client, s.WxWebCommon,s.WxWebXcg,s.Cookies,verifyContent,vefifyUsers, FRIEND_VEFIFY_OPCODE_REQUEST); err == nil {
		var res BaseResponse
		if err = json.Unmarshal(body, &res); err != nil {
			return err
		}
		if res.Ret != 0 {
			return fmt.Errorf("MakeFriend err")
		}
	}else{
		logs.Info("MakeFriend.err",err.Error())
	}
	return nil
}

// 预解析消息,生成特定消息事件
func (s *Session) PreparseMessage(msg *ReceivedMessage) {
	switch msg.MsgType {
	case MSG_SYS:
		// 处理新人进群的消息,派发-2
		if names := invitedFriendIntoGroupReg.FindStringSubmatch(msg.Content); len(names) > 2 {
			var memeberName = names[2]
			var group = s.Cm.GetContactByUserName(msg.FromUserName)

			msg.MsgType = MSG_NEW_GROUP_MEMBER
			msg.Data = NewGroupMemberMsgdata{
				NickName:memeberName,
				Group:group,
			}
			s.HandlerRegister.Runs(s,msg)
		}
	default:
		logs.Info("not preparse ",msg.MsgType)
	}
}
var invitedFriendIntoGroupReg = regexp.MustCompile("^(.+)邀请\"(.+)\"加入了群聊$")

//从同目录config.json中读取配置
func (this *Session) LoadConfig(){
	bs, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return
	}
	json.Unmarshal(bs, &this.Config)
}