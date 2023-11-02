package tcp

import (
	"github.com/flylib/gonet"
	"log"
	"net"
)

// Socket会话
type Session struct {
	//核心会话标志
	gonet.SessionIdentify
	//存储功能
	gonet.SessionAbility
	//累计收消息总数
	recvCount uint64
	//raw conn
	conn net.Conn
	//缓存数据，用于解决粘包问题
	cache []byte
}

// 新会话
func newSession(c *gonet.Context, conn net.Conn) *Session {
	is := c.CreateSession()
	s := is.(*Session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Session) Send(msgID uint32, msg any) error {
	buf, err := s.Context.Package(s, msgID, msg)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(buf)
	return err
}

func (s *Session) Close() error {
	return s.conn.Close()
}

// 接收循环
func (s *Session) recvLoop() {
	var buf = make([]byte, gonet.MTU)
	for {
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
		msg, unUsedCount, err := s.UnPackage(s, buf[:n])
		if err != nil {
			s.cache = nil
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		//存储未使用部分
		if unUsedCount > 0 {
			s.cache = append(s.cache, buf[n-unUsedCount-1:n]...)
		}
		s.Context.PushGlobalMessageQueue(msg)
	}
}
