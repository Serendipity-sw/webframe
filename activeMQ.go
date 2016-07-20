package main

import (
	"fmt"
	"github.com/go-stomp/stomp"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"sync"
)

var (
	conn *stomp.Conn
)

/**
初始化构造函数
创建人:邵炜
创建时间:2016年7月20日11:38:34
*/
func init() {
	deferinit.AddInit(mqConntion, closeActiveMQ, 991)
	deferinit.AddRoutine(mqMessageReceive)

}

/**
acitveMQ 服务连接请求
创建人:邵炜
创建时间:2016年7月20日10:31:06
*/
func mqConntion() {
	var (
		err error
	)
	conn, err = stomp.Dial("tcp", mqAddr)
	if err != nil {
		glog.Error("activeMQ conntion error! err: %s \n", err.Error())
		return
	}
}

/**
activeMQ服务消息发送
创建人:邵炜
创建时间:2016年7月20日11:17:31
*/
func mqMessageSend(content []byte) {
	err := conn.Send(fmt.Sprintf("/queue/%s", queueResult), "text/plain", content)
	if err != nil {
		glog.Error("mqMessageSend send message error! messageContent: %s err: %s \n", string(content), err.Error())
	}
}

/**
MQ消息接收
创建人:邵炜
创建时间:2016年7月19日15:53:57
*/
func mqMessageReceive(ch chan struct{}, wg *sync.WaitGroup) {
	sub, err := conn.Subscribe(fmt.Sprintf("/queue/%s", queue), stomp.AckAuto)
	if err != nil {
		glog.Error("mqMessageReceive activeMQ subscribe error! queueResult: %s err: %s \n", queueResult, err.Error())
		return
	}
	go func() {
		<-ch
		err = sub.Unsubscribe()
		if err != nil {
			glog.Error("mqMessageReceive activeMQ unsubscribe error! queueResult: %s err: %s \n", queueResult, err.Error())
		}
		wg.Done()
	}()
	var msg *stomp.Message
	for {
		msg = <-sub.C
		_ = msg.Body //订阅消息接收处理
	}
}

/**
activeMQ服务关闭
创建人:邵炜
创建时间:2016年7月20日11:28:19
*/
func closeActiveMQ() {
	err := conn.Disconnect()
	if err != nil {
		glog.Error("closeActiveMQ close activeMQ error! err: %s \n", err.Error())
	}
}
