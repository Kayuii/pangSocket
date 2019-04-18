package vsocket

import (
	"log"
	"net"
)

type socketTypes interface {
	ConnHandle(server *serv, sess *session)
	Pack(data []byte) []byte
}
type serv struct {
	EventPool     *routerMap
	SessionMaster *sessionManager
	SocketType    socketTypes
}

// NewService create
func NewService(socketType socketTypes) *serv {
	ser := &serv{
		EventPool:  NewRouterMap(),
		SocketType: socketType,
	}
	ser.SessionMaster = newSessonManager(ser)
	return ser
}

func (s *serv) Listening(address string) {
	tcpListen, err := net.Listen("tcp", address)

	if err != nil {
		panic(err)
	}
	go s.SessionMaster.HeartBeat(2)
	fd := uint32(0)
	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//调用握手事件
		if s.EventPool.OnHandel(fd, conn) == false {
			continue
		}
		s.SessionMaster.SetSession(fd, conn)
		go s.SocketType.ConnHandle(s, s.SessionMaster.GetSessionByID(fd))
		fd++
	}
}

func (s *serv) Hook(fd uint32, requestData map[string]string) bool {
	//调用接收消息事件
	if s.EventPool.OnMessage(fd, requestData) == false {
		return false
	}
	//requestData["fd"] = fmt.Sprintf("%d", fd)
	//路由
	if actionName, exit := requestData["action"]; exit {
		if s.EventPool.HookAction(actionName, fd, requestData) == false {
			return false
		}
	} else {
		if s.EventPool.HookModule(requestData["module"], requestData["method"], fd, requestData) == false {
			return false
		}
	}
	return true
}
