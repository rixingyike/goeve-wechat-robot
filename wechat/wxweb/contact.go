package wxweb

import (
	"encoding/json"
	"strings"
	"../logs"
)

var (
	// SpecialContact: special contacts map
	SpecialContact = map[string]bool{
		"filehelper":            true,
		"newsapp":               true,
		"fmessage":              true,
		"weibo":                 true,
		"qqmail":                true,
		"tmessage":              true,
		"qmessage":              true,
		"qqsync":                true,
		"floatbottle":           true,
		"lbsapp":                true,
		"shakeapp":              true,
		"medianote":             true,
		"qqfriend":              true,
		"readerapp":             true,
		"blogapp":               true,
		"facebookapp":           true,
		"masssendapp":           true,
		"meishiapp":             true,
		"feedsapp":              true,
		"voip":                  true,
		"blogappweixin":         true,
		"weixin":                true,
		"brandsessionholder":    true,
		"weixinreminder":        true,
		"officialaccounts":      true,
		"wxitil":                true,
		"userexperience_alarm":  true,
		"notification_messages": true,
	}
)

// ContactManager: contact manager
type ContactManager struct {
	cl []*User //contact list
}

// CreateContactManagerFromBytes: create contact maanger from bytes
func CreateContactManagerFromBytes(cb []byte) (*ContactManager, error) {
	var cr ContactResponse
	if err := json.Unmarshal(cb, &cr); err != nil {
		return nil, err
	}
	cm := &ContactManager{
		cl: cr.MemberList,
	}
	return cm, nil
}

// 移除用户列表
func (s *ContactManager) RemoveContact(mcts []map[string]interface{}) error {
	//logs.Info("RemoveContact.mcts",mcts)
	var users = make([]*User,0)
	if data,err := json.Marshal(mcts); err != nil {
		return err
	}else{
		if err := json.Unmarshal(data, &users); err != nil {
			return err
		}
	}

	for _, u := range users {
		for k, v := range s.cl {
			if v.UserName == u.UserName {
				// 如果原来已经存在,更新之
				s.cl = append(s.cl[:k],s.cl[k+1:]...)
				logs.Info("移除用户", u.NickName, u.DisplayName, u.UserName)
				break
			}
		}
	}

	return nil
}

// 更新用户列表
func (s *ContactManager) UpdateContact(mct map[string]interface{}) (u User, newFriend bool ,err error) {
	//logs.Info("UpdateContact.mcts",mcts)

	if data,err1 := json.Marshal(mct); err1 != nil {
		err = err1
		return
	}else{
		if err1 := json.Unmarshal(data, &u); err1 != nil {
			err = err1
			return
		}
	}

	newFriend = true
	for k, v := range s.cl {
		if v.UserName == u.UserName {
			// 如果原来已经存在,更新之
			s.cl = append(append(s.cl[:k], &u),s.cl[k+1:]...)
			newFriend = false
			logs.Info("修改用户", u.NickName, u.DisplayName, u.UserName)
			break
		}
	}
	if newFriend {
		// 如果是新的,添加
		s.cl = append(s.cl, &u)
		logs.Info("新增用户", u.NickName, u.DisplayName, u.UserName)
	}

	return
}

// AddConactFromBytes
// upate contact manager from bytes
func (s *ContactManager) AddConactFromBytes(cb []byte) error {
	var cr ContactResponse
	if err := json.Unmarshal(cb, &cr); err != nil {
		return err
	}
	s.cl = append(s.cl, cr.MemberList...)
	return nil
}

// GetContactByUserName
// get contact by UserName
func (s *ContactManager) GetContactByUserName(un string) *User {
	for _, v := range s.cl {
		if v.UserName == un {
			return v
		}
	}
	return nil
}

// GetGroupContact: get group contacts
func (s *ContactManager) GetGroupContact() []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if strings.Contains(v.UserName, "@@") {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

// GetStrangers: not group and not StarFriend
func (s *ContactManager) GetStrangers() []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if !strings.Contains(v.UserName, "@@") &&
			v.StarFriend == 0 &&
			v.VerifyFlag&8 == 0 &&
			!SpecialContact[v.UserName] {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

// GetContactByName: get contacts by name
func (s *ContactManager) GetContactByName(sig string) []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if v.NickName == sig || v.RemarkName == sig {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

// GetContactByQuanPin: get contact by User.QuanPin
func (s *ContactManager) GetContactByQuanPin(sig string) *User {
	for _, v := range s.cl {
		if v.PYQuanPin == sig || v.RemarkPYQuanPin == sig {
			return v
		}
	}
	return nil
}

// GetAll: get all contacts
func (s *ContactManager) GetAll() []*User {
	return s.cl
}
