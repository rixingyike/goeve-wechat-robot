package wxweb

import (
	"encoding/xml"
	"strconv"
	"strings"
)

const (
	// custom msg types
	MSG_NEW_FRIEND = -1 //新增好友
	MSG_NEW_GROUP_MEMBER = -2 //新人进群消息

	// msg types
	MSG_TEXT        = 1     // text message,文本,表情
	MSG_IMG         = 3     // image message
	MSG_VOICE       = 34    // voice message
	MSG_FRIEND_REQUEST = 37    // friend verification message 有人添加好友
	MSG_PF          = 40    // POSSIBLEFRIEND_MSG
	MSG_SCC         = 42    // shared contact card
	MSG_VIDEO       = 43    // video message
	MSG_EMOTION     = 47    // gif
	MSG_LOCATION    = 48    // location message
	MSG_LINK        = 49    // shared link message
	MSG_VOIP        = 50    // VOIPMSG
	// 这里可能是未读消息
	MSG_INIT        = 51    // wechat init message,貌似是登陆之后发的离线消息,通知客户端拉取,有一捆username,包括群username
	MSG_VOIPNOTIFY  = 52    // VOIPNOTIFY
	MSG_VOIPINVITE  = 53    // VOIPINVITE
	MSG_SHORT_VIDEO = 62    // short video message
	MSG_SYSNOTICE   = 9999  // SYSNOTICE
	// 新人请求加我好友,我通过:AddMsgCount: 1 ModContactCount: 1 DelContactCount: 0 ModChatRoomMemberCount: 0
	// 提示语:你已添加了寻找有缘人，现在可以开始聊天了
	// 我加别人好友,同样:AddMsgCount: 1 ModContactCount: 1 DelContactCount: 0 ModChatRoomMemberCount: 0
	// 提示语:我通过了你的朋友验证请求，现在我们可以开始聊天了
	MSG_SYS         = 10000 // system message
	MSG_WITHDRAW    = 10002 // withdraw notification message

	// verify opcode
	FRIEND_VEFIFY_OPCODE_DELETE int = 1//删除好友?
	FRIEND_VEFIFY_OPCODE_REQUEST int = 2//主动添加好友
	FRIEND_VEFIFY_OPCODE_ACCEPT int = 3//通过好友请求

	// 用户性别,与微信返回的值契合
	USER_SEX_NONE = 0
	USER_SEX_MALE = 1 //男
	USER_SEX_FEMALE = 2 //女

	// 联系人类似
	USER_TYPE_OFFICAL = 0//公众号
	USER_TYPE_FRIEND = 2//好友
	USER_TYPE_GROUP = 4//群组
	USER_TYPE_MEMBER = 8//群组成员
)

// Common: session config
type Common struct {
	AppId       string
	LoginUrl    string
	Lang        string
	DeviceID    string
	UserAgent   string
	CgiUrl      string//baseurl
	CgiDomain   string
	SyncSrv     string
	UploadUrl   string
	MediaCount  uint32
	RedirectUri string
}

// InitReqBody: common http request body struct
type InitReqBody struct {
	BaseRequest        *BaseRequest
	Msg                interface{}
	SyncKey            *SyncKeyList
	rr                 int
	Code               int
	FromUserName       string
	ToUserName         string
	ClientMsgId        int
	ClientMediaId      int
	TotalLen           int
	StartPos           int
	DataLen            int
	MediaType          int
	Scene              int
	Count              int
	List               []*User
	Opcode             int
	SceneList          []int
	SceneListCount     int
	VerifyContent      string
	VerifyUserList     []*VerifyUser
	VerifyUserListSize int
	skey               string
	MemberCount        int
	MemberList         []*User
	Topic              string

	//ChatRoomName 		string //群username
	//InviteMemberList string //进群后人数达到40人使用
	//AddMemberList string //进群后人数不到40人使用
}

// RevokeReqBody: revoke message api http request body
type RevokeReqBody struct {
	BaseRequest *BaseRequest
	ClientMsgId string
	SvrMsgId    string
	ToUserName  string
}

// LogoutReqBody: logout api http request body
type LogoutReqBody struct {
	sid string
	uin string
}

// BaseRequest: http request body BaseRequest
type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

// XmlConfig: web api xml response struct
type XmlConfig struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`
}

// SyncKey: struct for synccheck
type SyncKey struct {
	Key int
	Val int
}

// SyncKeyList: list of synckey
type SyncKeyList struct {
	Count int
	List  []SyncKey
}

// s.String output synckey list in string
func (s *SyncKeyList) String() string {
	strs := make([]string, 0)
	for _, v := range s.List {
		strs = append(strs, strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(strs, "|")
}

// User: contact struct
type User struct {
	Uin               int
	UserName          string
	NickName          string
	HeadImgUrl        string
	ContactFlag       int
	MemberCount       int
	MemberList        []*User
	RemarkName        string
	PYInitial         string
	PYQuanPin         string
	RemarkPYInitial   string
	RemarkPYQuanPin   string
	HideInputBarFlag  int
	StarFriend        int
	Sex               int
	Signature         string
	AppAccountFlag    int
	Statues           int
	AttrStatus        uint32
	Province          string
	City              string
	Alias             string
	VerifyFlag        int
	OwnerUin          int
	WebWxPluginSwitch int
	HeadImgFlag       int
	SnsFlag           int
	UniFriend         int
	DisplayName       string
	ChatRoomId        int
	KeyWord           string
	EncryChatRoomId   string
	IsOwner           int
	MemberStatus      int

	Type int //类别
}

// 群组成员在群内显示名称是displayName
func (this *User) GetDisplayName() string {
	if len(this.DisplayName) > 0 {
		return this.DisplayName
	}
	return this.NickName
}

// TextMessage: text message struct
type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

// MediaMessage
type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
	MediaId      string
}

// EmotionMessage: gif/emoji message struct
type EmotionMessage struct {
	ClientMsgId  int
	EmojiFlag    int
	FromUserName string
	LocalID      int
	MediaId      string
	ToUserName   string
	Type         int
}

// BaseResponse: web api http response body BaseResponse struct
type BaseResponse struct {
	Ret    int
	ErrMsg string
}

// 接口调用是否成功
func (this *BaseResponse) IsSuccess() bool {
	return this.Ret == 0
}

// ContactResponse: get contact response struct
type ContactResponse struct {
	BaseResponse *BaseResponse
	MemberCount  int
	MemberList   []*User
	Seq          int
}

// GroupContactResponse: get batch contact response struct
type GroupContactResponse struct {
	BaseResponse *BaseResponse
	Count        int
	ContactList  []*User
}

// VerifyUser: verify user request body struct
type VerifyUser struct {
	Value            string //username
	VerifyUserTicket string //ticket
}

// ReceivedMessage: for received message
type ReceivedMessage struct {
	IsGroup      bool
	MsgId        string
	Content      string
	FromUserName string
	ToUserName   string
	Who          string //群自是谁在发消息,其username
	MsgType      int

	// 新增
	OriginalMsg      map[string]interface{} // 原始消息数据
	Data interface{} // 附带的进一步解析的特定消息类型的数据
	IsSendedByMySelf bool //是否为自己所发,在群内处理
	IsAtMe bool //
}