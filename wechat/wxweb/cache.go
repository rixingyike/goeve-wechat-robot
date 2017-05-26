/*
 * 缓存cookie等信息到本地,实现自动恢复登陆
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/5/12
 */
package wxweb

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../logs"
	"os"
	//"net/url"
)

const CACHE_FILE_PATH = "./cache.json"

type Cache struct {
	session *Session 	`json:"-"`

	//需要序列化的对象
	WxWebXcg        XmlConfig
	Cookies         []*http.Cookie
	RedirectUri string
}

// 装载缓存,从本地文件或数据库
func (this *Cache) Load() (err error) {
	this.Cookies = make([]*http.Cookie,0)
	logs.Info("load cache..")

	bs, err := ioutil.ReadFile(CACHE_FILE_PATH)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bs, this);err == nil {
		this.session.WxWebXcg = &this.WxWebXcg
		this.session.Cookies = this.Cookies
		this.session.WxWebCommon.RedirectUri = this.RedirectUri
		this.session.SetCookiesAfterScanQrcode()
	}else{
		logs.Info("load err",err.Error())
	}
	return
}

// 写入缓存
func (this *Cache) Write() (err error) {
	this.WxWebXcg = *this.session.WxWebXcg
	this.Cookies = this.session.Cookies
	this.RedirectUri = this.session.WxWebCommon.RedirectUri

	b, err := json.Marshal(this)
	if err != nil {
		logs.Info(`write/marshal cache error: %v`, err)
		return
	}

	oflag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	file, err := os.OpenFile(CACHE_FILE_PATH, oflag, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(b)
	return
}