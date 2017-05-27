package main

import (
	"encoding/base64"
	log "glog"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// //H3c6800Exec   尝试netconf操作
// func H3c6800Exec(h3cStruct H3cCommand) int {
// 	var CMDRETNUM = 200 //默认执行返回值200，表示正常
//
// 	//  base64解密交换机密码
// 	XXswitchRealPassword, _ := base64.URLEncoding.DecodeString(h3cStruct.SwitchPassword)
// 	//  解密后的密码除去前10位之后的才是真正的密码
// 	switchRealPassword := string(XXswitchRealPassword)[10:]
//
// 	//组装ssh的配置文件
//
// 	config := &ssh.ClientConfig{
// 		User: h3cStruct.SwitchUsername,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(switchRealPassword),
// 		},
// 		Timeout: h3cStruct.SwitchTimeout * time.Second,
// 		Config: ssh.Config{
// 			Ciphers: []string{"aes128-cbc"},
// 		},
// 	}
// 	// config.Config.Ciphers = append(config.Config.Ciphers, "aes128-cbc")
// 	clinet, err := ssh.Dial("tcp", h3cStruct.SwitchIPAndPort, config)
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 303, err)
// 		return 303
// 	}
//
// 	session, err := clinet.NewSession()
// 	defer session.Close()
//
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 304, err)
// 		return 304
// 	}
//
// 	modes := ssh.TerminalModes{
// 		ssh.ECHO:          1,     // disable echoing
// 		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
// 		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
// 	}
//
// 	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
//
// 		log.Errorf(" ErrorCode:%d , %s", 305, err)
// 		return 305
// 	}
//
// 	w, err := session.StdinPipe()
// 	if err != nil {
//
// 		log.Errorf(" ErrorCode:%d , %s", 306, err)
// 		return 306
// 	}
// 	r, err := session.StdoutPipe()
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 307, err)
// 		return 307
// 	}
// 	e, err := session.StderrPipe()
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 308, err)
// 		return 308
// 	}
//
// 	in, out := MuxShell(w, r, e)
//
// 	if err := session.Shell(); err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 309, err)
// 		return 309
// 	}
// 	<-out // 第一次输出为登陆后的设备提示信息， Copyright (c) 2004-2013 Hangzhou H3C Tech. Co.,  可以不打印
// 	in <- "xml"
// 	//log.Infoln(<-out) //进入系统模式
// 	fmt.Println(<-out)
//
// 	in <- `
// 	<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
// 	 <capabilities>
// 	    <capability>
// 	            urn:ietf:params:netconf:base:1.0
// 	    </capability>
// 	  </capabilities>
// 	</hello>]]>]]>
// 	`
//
// 	in <- `
// 	<rpc message-id="103" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
// 	         <close-session/>
// 	</rpc>]]>]]>
// 	`
//
// 	in <- "quit"
//
// 	for resultStr := range out { //  循环从out管道中遍历出最后两条的输出结果
// 		log.Infoln(resultStr)
// 		fmt.Println(<-out)
// 	}
//
// 	session.Wait()
// 	log.Infof("%s\n", "<--ExecClose")
//
// 	return CMDRETNUM
//
// }

//H3cCommandViewRoute   本函数作用在于查看某条路由是否存在
// func H3cCommandViewRoute(h3cStruct H3cCommand) ([]string, int) {
//
// 	var routeSlice []string //准备存放过滤出来的静态路由
//
// 	var CMDRETNUM = 200 //默认执行返回值0，
//
// 	//  base64解密交换机密码
// 	XXswitchRealPassword, _ := base64.URLEncoding.DecodeString(h3cStruct.SwitchPassword)
// 	//  解密后的密码除去前10位之后的才是真正的密码
// 	switchRealPassword := string(XXswitchRealPassword)[10:]
//
// 	//组装ssh的配置文件
//
// 	config := &ssh.ClientConfig{
// 		User: h3cStruct.SwitchUsername,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(switchRealPassword),
// 		},
// 		Timeout: h3cStruct.SwitchTimeout * time.Second,
// 		Config: ssh.Config{
// 			Ciphers: []string{"aes128-cbc"},
// 		},
// 	}
// 	// config.Config.Ciphers = append(config.Config.Ciphers, "aes128-cbc")
// 	clinet, err := ssh.Dial("tcp", h3cStruct.SwitchIPAndPort, config)
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 303, err)
// 		return routeSlice, 303
// 	}
// 	defer clinet.Close()
//
// 	session, err := clinet.NewSession()
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 304, err)
// 		return routeSlice, 304
// 	}
// 	defer session.Close()
//
//
// 	modes := ssh.TerminalModes{
// 		ssh.ECHO:          1,     // disable echoing
// 		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
// 		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
// 	}
//
// 	// 定义窗口
// 	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
//
// 		log.Errorf(" ErrorCode:%d , %s", 305, err)
// 		return routeSlice, 305
// 	}
//
// 	w, err := session.StdinPipe()
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 306, err)
// 		return routeSlice, 306
// 	}
// 	r, err := session.StdoutPipe()
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 307, err)
// 		return routeSlice, 307
// 	}
// 	e, err := session.StderrPipe()
// 	if err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 308, err)
// 		return routeSlice, 308
// 	}
//
// 	in, out := MuxShell(w, r, e)
//
// 	if err := session.Shell(); err != nil {
// 		log.Errorf(" ErrorCode:%d , %s", 309, err)
// 		return routeSlice, 309
// 	}
// 	<-out // 第一次输出为登陆后的设备提示信息， Copyright (c) 2004-2013 Hangzhou H3C Tech. Co.,  丢弃显示
//
// 	in <- h3cStruct.SwitchCommand
//
// 	in <- "quit" //  这个一输入就断开连接了， 得不到out的输出
//
// 	for resultStr := range out { //  循环从out管道中便利出执行结果
// 		log.Infoln(resultStr)
//
// 		// 用于过滤的正则表达式也定制为只显示10.2XX类的路由，防止误伤到其他交换机路由。
// 		re, _ := regexp.Compile("ip route-static 10\\.2[0-9]{2}\\.[0-9]{1,3}\\.[0-9]{1,3} .* [0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
//
// 		allinone := re.FindAll([]byte(resultStr), -1)
//
// 		for _, u := range allinone {
// 			routeSlice = append(routeSlice, string(u)) // 填充到事先准备好的slice中
//
// 		}
//
// 		if strings.Contains(resultStr, "^") { //如果有别的错误类型也可以在这里添加返回码
// 			CMDRETNUM = 400 //  命令不能识别
// 		}
//
// 	}
//
// 	session.Wait()
//
// 	log.Infof("%s\n", "<--ViewClose")
//
// 	return routeSlice, CMDRETNUM
//
// }

//H3cCommandExec   本函数作用执行一条添加路由的操作
func H3cCommandExec(h3cStruct H3cCommand) int {

	var CMDRETNUM = 200 //默认执行返回值200，表示正常

	//  base64解密交换机密码
	XXswitchRealPassword, _ := base64.URLEncoding.DecodeString(h3cStruct.SwitchPassword)
	//  解密后的密码除去前10位之后的才是真正的密码
	switchRealPassword := string(XXswitchRealPassword)

	//组装ssh的配置文件

	config := &ssh.ClientConfig{
		User: h3cStruct.SwitchUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(switchRealPassword),
		},
		Timeout: h3cStruct.SwitchTimeout * time.Second,
		Config: ssh.Config{
			Ciphers: []string{"aes128-cbc"},
		},
	}
	// config.Config.Ciphers = append(config.Config.Ciphers, "aes128-cbc")
	clinet, err := ssh.Dial("tcp", h3cStruct.SwitchIPAndPort, config)
	if err != nil {
		log.Errorf(" ErrorCode:%d , %s", 303, err)
		return 303
	}
	defer clinet.Close()

	session, err := clinet.NewSession()

	if err != nil {
		log.Errorf(" ErrorCode:%d , %s", 304, err)
		return 304
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {

		log.Errorf(" ErrorCode:%d , %s", 305, err)
		return 305
	}

	w, err := session.StdinPipe()
	if err != nil {

		log.Errorf(" ErrorCode:%d , %s", 306, err)
		return 306
	}
	r, err := session.StdoutPipe()
	if err != nil {
		log.Errorf(" ErrorCode:%d , %s", 307, err)
		return 307
	}
	e, err := session.StderrPipe()
	if err != nil {
		log.Errorf(" ErrorCode:%d , %s", 308, err)
		return 308
	}

	in, out := MuxShell(w, r, e)

	if err := session.Shell(); err != nil {
		log.Errorf(" ErrorCode:%d , %s", 309, err)
		return 309
	}
	<-out // 第一次输出为登陆后的设备提示信息， Copyright (c) 2004-2013 Hangzhou H3C Tech. Co.,  可以不打印
	in <- "sy"
	log.Infoln(<-out) //进入系统模式

	for _, _u := range strings.Split(h3cStruct.SwitchCommand, ";") {

		in <- _u //执行合规的那四种命令
		resultStr := <-out
		log.Infoln(resultStr)

		if strings.Contains(resultStr, "Route doesn't exist") {
			CMDRETNUM = 401 //操作的路由不存在
		}
		if strings.Contains(resultStr, "^") { //如果已经用正则表达式严谨的匹配过ip应该不会出现这里的错误
			CMDRETNUM = 400 //  命令不能识别
		}

	}

	//当执行的是多条命令时， 只有所有的命令都没有报错，才会执行save force ，只要有一条命令有问题，则不执行save动作。
	if CMDRETNUM == 200 { // 如果之前的命令都没有执行错误，则保存配置，否则不保存直接退出。<
		in <- "save  force"
	}

	in <- "quit" //退出到普通用户模式。
	in <- "quit" //  这个一输入就断开连接了， 得不到out的输出

	for resultStr := range out { //  循环从out管道中遍历出最后两条的输出结果
		log.Infoln(resultStr)
	}

	session.Wait()

	log.Infof("%s\n", "<--ExecClose")

	return CMDRETNUM

}

func checkCmd(cmdStr string, cmdLevel int) bool { //是否包含以下这些掩码，包含的话说明威胁不大，如果不包含则不运行这条命令；
	InitStatus := true

	for i, _u := range strings.Split(cmdStr, ";") {
		// 如果命令行中有分号，则需要

		log.Infof("chechcmd %d:%s\n", i, _u)

		//如果行尾多加了一个分号
		if len(_u) < 10 {
			log.Warning("attention: len checkcmd < 10 ,too short")
		}

		//根据cmdlevel的不通，实施不同的匹配规则
		switch {

		case cmdLevel == 1:
			if m, _ := regexp.MatchString("^display .*", _u); m {
				//fmt.Println("---display----->", m, mm)
				continue
			}
			// 一组命令中只要有一条命令不合规，则返回false，此处不做详细的提示。
			log.Warning("checkcmd error," + _u)
			InitStatus = false
			break

		case cmdLevel == 2:
			if m, _ := regexp.MatchString("ip route-static 10\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3} 255\\.255\\.255\\.0 [0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}", _u); m { //操作的路由是10.2XX.X.X/24，以防止增删不当
				continue
			}

			// 一组命令中只要有一条命令不合规，则返回false，此处不做详细的提示。
			log.Warning("checkcmd error," + _u)
			InitStatus = false
			break

		default:
			log.Warning("checkcmd error, unknown cmdLevel")
			InitStatus = false
		}

	}
	//fmt.Println("for oever")
	return InitStatus

}
