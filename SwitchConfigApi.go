package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	log "glog"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	limmit "golang.org/x/net/netutil"
)

const (
	//Base64Table  64个字符的顺序不同会产生不同的BASE64编码效果，起到轻量级的加密作用。
	Base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXY+-/0123456789abcdefghijklmnopqrstuvwxyz"
)

// AuthString  全局变量，用于存储server端用户自定义的key，通过全局变量传递进入handle函数。
var AuthString string

// HongCoder 自定义的base64编码表，可以作为一种简单的加密方法
var HongCoder = base64.NewEncoding(Base64Table)

func main() {

	//  参数部分使用flag包实现
	var InputIPPort = flag.String("IpPort", "0.0.0.0:8083", "本地监听的IP和端口")
	var InputAuthKey = flag.String("AuthKey", "admin:admin2", "客户端需要在http头中包含这个key的特殊编码作为认证，以证明是合法客户端请求") //可扩展为从中提取出用户角色信息，思路来源于docker仓库
	var InputHelp = flag.Bool("help", false, "help info")
	var InputVersion = flag.Bool("version", false, "version information")
	var InputLimit = flag.Int("Limit", 1, "同一时间访问该接口的client端的数量，越低对交换机cpu影响越小")
	//var InputLogDir = flag.String("log_dir", "", "log file dir")
	flag.Parse()

	//判断日志文件夹是否存在， 不存在则创建一个
	//fmt.Println(*InputLogDir)

	//fmt.Println(flag.NArg())    以上指定的参数都不计数
	if (flag.NArg() > 0) || (*InputHelp) {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]...\nExample:%s\n注意：log_dir指定的需要是已存在的目录。 \n\n", os.Args[0], " ./SwitchConfigApi  -AuthKey=admin:admin2 -IpPort=0.0.0.0:8083   --log_dir=./log   -alsologtostderr ")
		flag.PrintDefaults()

		os.Exit(0)

	}

	if *InputVersion { //如果用户是想看版本信息，则显示完就退出
		fmt.Fprintln(os.Stderr, "version 1.0 ,Copyright@hongfeio.o@163.com")
		os.Exit(0)
	}

	AuthString = *InputAuthKey //复制给全局变量，传入hander

	// 以下回显中包含，用户自定义的ip＋端口，认证的Key，括弧内是特殊编码后的字符串，客户端需要把这段字符串包含在http头中进行API的验证， 还显示了http server的并发数限制。
	log.Infof("ListenPort:%s,AuthKey:%s(ClientAuthKey:%s),LimitListener:%d\n", *InputIPPort, *InputAuthKey, HongCoder.EncodeToString([]byte(*InputAuthKey)), *InputLimit)
	fmt.Printf("ListenPort:%s,AuthKey:%s,ClientAuthKey:%s,LimitListener:%d\n", *InputIPPort, *InputAuthKey, HongCoder.EncodeToString([]byte(*InputAuthKey)), *InputLimit)

	// http  监听函数
	l, err := net.Listen("tcp", *InputIPPort)
	if err != nil {
		log.Warning(err.Error())

	}
	defer l.Close()
	defer log.Flush() // 确保日志都写入文件中。

	l = limmit.LimitListener(l, *InputLimit) //  接口并行的数量默认限制为1，即同一个时间只有一个客户端可以访问接口，1 防止多个客户端同时操作一台交换机  ，2 保证日志的可读性

	http.HandleFunc("/", handler) //设置访问的路由
	http.Serve(l, nil)            //这个函数是一个死循环， 一直监听客户端的请求。

}

//handler  主要函数
func handler(w http.ResponseWriter, r *http.Request) {

	// 使手机端可以跨域请求
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin,Content-Type, Accept,User-Agent,Authorization,Accept-Encoding")
	defer r.Body.Close()

	// 只接受POST和GET方法， 因为手机端会发送出OPTIONS的请求，过滤掉
	if (r.Method == "POST") || (r.Method == "GET") {

		body, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(body)

		// 日志显示客户端发起的请求信息
		log.Infof("<--begin\n Body:%s\n Content-Type:%s\n Authorization:%s\n Content-Length:%d\n User-Agent:%s\n Method:%s\n Host:%s\n", bodyStr, r.Header["Content-Type"], r.Header["Authorization"], r.ContentLength, r.UserAgent(), r.Method, r.Host)

		// 判断参数是否完整，包括json串和认证字符串，任意为空返回404
		if (bodyStr == "") || (r.Header["Authorization"] == nil) {
			log.Errorf(" ErrorCode:%d  no json or no Authorization!\n", 404)
			fmt.Fprintf(w, "%d\n", 404)
			return
		}

		// 判断http中的认证是否正确
		AuthStr := r.Header["Authorization"]
		userpwdBase64encode := strings.Split(AuthStr[0], " ")[1]
		uDec, _ := HongCoder.DecodeString(userpwdBase64encode) //解密客户端提供上来的字符串

		log.Infof(" ServerAuthKey: %s  VS ClientAuthKeyDecode: %s\n", AuthString, string(uDec)) //服务器端定义的AuthKey和客户端上报的AuthKey对比显示。

		if string(uDec) != AuthString { //   AuthString是用户设定的参数， client端请求API的时候需要把这个字符串Hong_base64自定义运算，放入http的头部。

			log.Errorf(" ErrorCode:%d  Authorization Fail!\n", 511)
			fmt.Fprintf(w, "%d\n", 511)
			return

		}

		//  开始解析json串
		var h3c H3cCommand

		if err := json.Unmarshal(body, &h3c); err != nil {

			log.Warningf(" WarningCode:%d  Json Decode Fail! %s\n", 300, err)
			fmt.Fprintf(w, "%d\n", 300)
			return
		}

		// 根据不同的cmdlevel丢给不同的函数处理。
		switch {
		// case h3c.SwitchCmdLevel == 1:
		//
		// 	if !checkCmd(h3c.SwitchCommand, 1) {
		// 		log.Warningf(" WarningCode:%d  CheckCmd level 1 Fail!\n", 301)
		// 		fmt.Fprintf(w, "%d\n", 301)
		// 		return
		// 	}
		// 	routeString, retnum := H3cCommandViewRoute(h3c)
		// 	if retnum != 200 {
		// 		log.Warningf(" WarningCode:%d H3cCommandViewRoute !\n", retnum)
		// 		fmt.Fprintf(w, "%d\n", retnum)
		// 		return
		// 	}
		// 	fmt.Fprintf(w, "%s\n", strings.Join(routeString, ";"))
		// 	return

		case h3c.SwitchCmdLevel == 2:
			//  判断命令是否合规的函数，合规返回true，不合规返回false
			if !checkCmd(h3c.SwitchCommand, 2) {
				log.Warningf(" WarningCode:%d  CheckCmd level 2 Fail!\n", 301)
				fmt.Fprintf(w, "%d\n", 301)
				return
			}
			fmt.Fprintf(w, "%d\n", H3cCommandExec(h3c)) // 200表示正常执行，返回400表示命令有错误(^) ，返回值可以是一个json串，可以参考vservermap2的代码
			return

		default:
			log.Warningf(" WarningCode:%d  cmd level is unknown!\n", 302)
			fmt.Fprintf(w, "%d\n", 302)
			return

		}
	}

}
