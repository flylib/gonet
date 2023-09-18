package tcp

import (
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
)

var _ ISession = new(session)

// Socket会话
type session struct {
	//核心会话标志
	SessionIdentify
	//存储功能
	SessionAbility
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn net.Conn
	//缓存数据，用于解决粘包问题
	cache []byte
}

// 新会话
func newSession(c *Context, conn net.Conn) *session {
	is := c.CreateSession()
	s := is.(*session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msg interface{}) error {
	data, err := s.Context.Package(msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(data)
	return err
}

func (s *session) Close() error {
	return s.conn.Close()
}

// 接收循环
func (s *session) recvLoop() {
	for {
		var buf = make([]byte, MTU)
		n, err := s.conn.Read(buf)
		if err != nil {
			s.Context.RecycleSession(s, err)
			return
		}
		if n == 0 {
			continue
		}
		//如果有粘包未处理数据部分，放入本次进行处理
		if len(s.cache) > 0 {
			buf = append(s.cache, buf[:n]...)
			n = len(buf)
			s.cache = nil
		}
		msg, unUsedCount, err := s.UnPackage(buf[:n])
		if err != nil {
			s.cache = nil
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		//存储未使用部分
		if unUsedCount > 0 {
			s.cache = append(s.cache, buf[n-unUsedCount-1:n]...)
		}
		s.Context.PushGlobalMessageQueue(s, msg)
	}
}
