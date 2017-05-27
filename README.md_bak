SwitchConfigApi
---------------------

简介
-----------
1. 本程序为H3C交换机封装了一层北向接口， 用户可以通过RESTful的方式提交JSON给本接口， 本程序会解析JSON内容去交换机通过模拟登录的方式执行相应的命令。         
2. 本程序只演示向交换机写入和删除路由的操作。


Fire Up
---------------
```
启动接口服务
 ./server.sh
ListenPort:0.0.0.0:8083,AuthKey:admin:admin2,ClientAuthKey:YWRh-WsuYWRh-Wsm,LimitListener:1    // 默认监听8083端口，默认客户端使用YWRh-WsuYWRh-Wsm进行头部验证，默认限制并发为1


client 调用方法
$  curl  -d@'setanddel.json'  http://127.0.0.1:8083  -H "Authorization: Basic YWRh-WsuYWRh-Wsm"   //插入路由
200

```

client提交的JSON举例
-----------------
```

[setanddel.json]        
{
	"switch_username": "yihf",            // ssh登录交换机所用的用户名
	"switch_password": "MTIzNDU2Nzg=",         // 如果你的交换机密码为12345678，通过命令计算出base64转码的结果 ：   echo -n "12345678" |  base64
	"switch_Cmd_level": 2,                   //选择2为执行命令，选择其它对应的处理函数会有不同可自行开发
	"switch_command": "ip route-static 10.201.88.0 255.255.255.0 10.10.88.129;undo ip route-static 10.201.88.0 255.255.255.0 10.10.88.129",   //多条命令时，使用分号分割, 这里的举例是写入了一条路由又删了一条
	"switch_ipandport": "10.10.100.20:22",      // 交换机IP+port
	"switch_timeout": 10                        // 超时时间设置
}
```


状态返回码
--------------
1：  静态路由配置中包含查询的路由
0:   静态路由配置中不包含查询的路由
200：  命令没有报错。
400:   命令不识别
401：  操作的路由不存在
404 :  body 或者认证头部不存在
511:    头部认证失败
300 ：  json串解析失败
301 ：  命令不合规, 包含了level1 和level2 两种fail
302 :   未知的cmdlevel数值
303 :   switch timeout , user or password fail!
304-309 :   cmd_exec 函数执行过程中的错误。



CLI正则验证，保障“安全”的命令才能被执行
------------------------------
1. 为了确保只有安全的命令才会被提交，目前只允许以下5种命令格式，如果需要放宽限制，请修改checkCmd函数

```
2. switch_Cmd_level":2 时，   只放行这两种命令，不符合正则表达式的命令不会被执行。
      * * ip route-static  10.XXX.XXX.XXX 255.255.255.0   XXX.XXX.XXX.XXX       增加路由命令
      ＊＊ undo ip route-static  10.XXX.XXX.XXX 255.255.255.0  XXX.XXX.XXX.XXX   删除路由命令
```  


安全机制
----------------
1. 客户端需要有合法的Http头部验证，这个验证字符串是在本程序启动的时候生产，客户端必须使用这个特定字符串才能正确调用接口，可通过参数自定义头验证。
2. 客户端JSON中包含的交换机密码， 目前采用：base64(真实密码)加密。  
3. 服务端启动接口的时候可以定义并发数量，控制同一时间操作交换机的client数量。
4. 客户端提交的CLI命令，必须通过正则验证， 避免误操作。



centos7创建服务
-------------
```
vi /usr/lib/systemd/system/hongswitchapi.service
[Unit]
Description=api of switch
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/home/workspace/src/SwitchConfigApi/SwitchConfigApi -log_dir=/home/workspace/src/SwitchConfigApi/log         
Restart=on-failure

[Install]
WantedBy=multi-user.target
```


启动服务并设置为开机启动
------------------
```
systemctl start hongswitchapi.service
systemctl enable hongswitchapi.service
```


TODO
-------------
1.  服务端可以对client的HOSTIP进行验证，防止非法ip调用。



## 开发环境
golang 1.8

## 作者介绍
yihongfei  QQ:413999317   MAIL:yihf@liepin.com

CCIE 38649


## 寄语
为网络自动化运维尽绵薄之力，每一个网工都可以成为NetDevOps
