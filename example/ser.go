package main

import (
	"log"
	"net"
	// pangsocket "github.com/Kayuii/pangSocket"
)

var ser = pangsocket.newService(&ps.service.tcpSocket)

//框架事件
//----------------------------------------------------------------------------------------------------------------------
type event struct {
}

//客户端握手成功事件
func (e event) OnHandel(fd uint32, conn net.Conn) bool {
	log.Println(fd, "链接成功类")
	return true
}

//断开连接事件
func (e event) OnClose(fd uint32) {
	log.Println(fd, "链接断开类")
}

//接收到消息事件
func (e event) OnMessage(fd uint32, msg map[string]string) bool {
	log.Println("这个是接受消息事件", msg)
	return true
}

// Test ---
//框架业务逻辑
type Test struct {
}

// Default : Default
func (t Test) Default(fd uint32, data map[string]string) bool {
	log.Println("default")
	return true
}

// BeforeRequest : BeforeRequest
func (t Test) BeforeRequest(fd uint32, data map[string]string) bool {
	log.Println("before")
	return true
}

// AfterRequest : AfterRequest
func (t Test) AfterRequest(fd uint32, data map[string]string) bool {
	log.Println("after")
	return true
}

// Hello : Hello
func (t Test) Hello(fd uint32, data map[string]string) bool {
	log.Println("收到消息了")
	log.Println(data)
	ser.SessionMaster.WriteByid(fd, []byte("hehehehehehehehe"))
	return true
}

// //-----------------------------------------
// func main() {
// 	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Llongfile)
// 	ser.EventPool.RegisterEvent(&event{})
// 	ser.EventPool.RegisterStructFun("test", &Test{})
// 	ser.Listening(":8565")
// }
