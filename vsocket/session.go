package vsocket

import (
	"log"
	"net"
	"sync"
	"time"
)

// 一个session代表一个连接
// session for connect
type session struct {
	ID    uint32
	Con   net.Conn
	times int64
	lock  sync.Mutex
}

// Session : New Session for session
func newSession(id uint32, con net.Conn) *session {
	return &session{
		ID:    id,
		Con:   con,
		times: time.Now().Unix(),
	}
}

func (s *session) write(msg string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, errs := s.Con.Write([]byte(msg))
	return errs
}

func (s *session) close() {
	s.Con.Close()
}

func (s *session) updateTime() {
	s.times = time.Now().Unix()
}

// --SESSION管理类---

type sessionManager struct {
	isWebSocket bool
	ser         *service
	sessions    sync.Map
}

// SessonManager : new SessonManager
func newSessonManager(server *service) *sessionManager {
	if server == nil {
		return nil
	}
	return &sessionManager{
		ser: server,
	}
}

func (sm *sessionManager) GetSessionByID(id uint32) *session {
	tem, exit := sm.sessions.Load(id)
	if exit {
		if sess, ok := tem.(*session); ok {
			return sess
		}
	}
	return nil
}

func (sm *sessionManager) SetSession(fd uint32, conn net.Conn) {
	sess := newSession(fd, conn)
	sm.sessions.Store(fd, sess)
}

//关闭连接并删除
func (sm *sessionManager) DelSessionByID(id uint32) {
	tem, exit := sm.sessions.Load(id)
	if exit {
		if sess, ok := tem.(*session); ok {
			sess.close()
		}
	}
	sm.sessions.Delete(id)
}

//向所有客户端发送消息
func (sm *sessionManager) WriteToAll(msg []byte) {
	msg = sm.ser.SocketType.Pack(msg)
	sm.sessions.Range(func(key, val interface{}) bool {
		if val, ok := val.(*session); ok {
			if err := val.write(string(msg)); err != nil {
				sm.DelSessionByID(key.(uint32))
				log.Println(err)
			}
		}
		return true
	})
}

//向单个客户端发送信息
func (sm *sessionManager) WriteByID(id uint32, msg []byte) bool {
	//把消息打包
	msg = sm.ser.SocketType.Pack(msg)

	tem, exit := sm.sessions.Load(id)
	if exit {
		if sess, ok := tem.(*session); ok {
			if err := sess.write(string(msg)); err == nil {
				return true
			}
		}
	}
	sm.DelSessionByID(id)
	return false
}

//心跳检测   每秒遍历一次 查看所有sess 上次接收消息时间  如果超过 num 就删除该 sess
func (sm *sessionManager) HeartBeat(num int64) {
	for {
		time.Sleep(time.Second)
		sm.sessions.Range(func(key, val interface{}) bool {
			tem, ok := val.(*session)
			if !ok {
				return true
			}

			if time.Now().Unix()-tem.times > num {
				sm.DelSessionByID(key.(uint32))
			}
			return true
		})

	}
}
