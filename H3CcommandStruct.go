package main

import "time"

//H3cCommand 用户发送来的json字符串
//这个结构体同时被Server和client使用，如需更改结构体请同步。
type H3cCommand struct {
	SwitchUsername  string        `json:"switch_username"`
	SwitchPassword  string        `json:"switch_password"`
	SwitchCmdLevel  int           `json:"switch_cmd_level"` //   level1 的命令只用于display ; level 2 的命令是在sy模式下运行的
	SwitchCommand   string        `json:"switch_command"`
	SwitchIPAndPort string        `json:"switch_ipandport"`
	SwitchTimeout   time.Duration `json:"SwitchTimeout"` // second
}
