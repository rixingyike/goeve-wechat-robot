## goeve-wechat-robot

_这是一个简单的智能助理微信机器人,主要有以下功能:_

1. 自动通过好友申请，主动发送问候语（自定义）
1. 回复关键字自动拉入微信群（自定义）
1. 群内新人加入，主动发送欢迎语（自定义）

有go语言开发经验的同学,可以自行clone源码、编译。没有经验的同学,可以前往[知乎live](https://www.zhihu.com/lives/846360223609413632),购买分享课程,参照语音说明配置使用。
价格不贵,也就是请作者喝一怀咖啡的价格。

### 最新版本

![v1版本截图](./wechat/v1.jpeg)


### 如何使用?

前往网盘[https://pan.baidu.com/s/1jIuqoF8](https://pan.baidu.com/s/1jIuqoF8) 下载以下文件,密码在live内:

1）下载执行文件
windows系统请下载sim-robot_win32.exe或sim-robot_win64.exe（xp7使用win32，win7以上使用win64位）
mac系统请下载sim-robot_darwin

2）下载配置文件
无论是什么系统，都需要下载config.json，与执行文件放在同一目录下。

3）参照live说明，配置好json文件，单击运行
详细配置说明见live：[零编程打造一款私人智能助理](https://www.zhihu.com/lives/846360223609413632)

### 听live分享

查看原理、配置及源码使用说明,请前往知乎live:
[零编程打造一款私人智能助理](https://www.zhihu.com/lives/846360223609413632)

不想下载执行文件的同学,可自行搭建环境从源码编译,也是一样的。

### 参与讨论

我还在开发更多的功能,扫描下方二维码,加阿娟微信,回复"智能",她将拉你进智能机器人开发讨论群:
![二维码](./wechat/qrcode.png)

马上与大牛们们开始热情的互动吧~

(注:本项目源码基于[qrterminal](github.com/mdp/qrterminal)、[wechat-go](https://github.com/songtianyi/wechat-go)、[beego](https://github.com/beego/bee)等类库修改,在此一并感谢。)
