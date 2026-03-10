package tcp

import (
	"github.com/flylib/gonet"
	"net"
)

// Session is the TCP connection session.
type Session struct {
	gonet.SessionCommon
	recvCount uint64
	conn      net.Conn
	cache     []byte
}

// newSession gets an idle session from the pool and attaches conn.
// Returns nil if the max session limit has been reached.
func newSession(c *gonet.AppContext[*Session], conn net.Conn) *Session {
	s, ok := c.GetIdleSession()
	if !ok {
		return nil
	}
	s.conn = conn
	c.GetEventHandler().OnConnect(s)
	return s
}

func (s *Session) RemoteAddr() net.Addr { return s.conn.RemoteAddr() }

// Send is safe to call from multiple goroutines.
func (s *Session) Send(msgID uint32, msg any) error {
	buf, err := s.GetContext().Package(s, msgID, msg)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	_, err = s.conn.Write(buf)
	return err
}

func (s *Session) Close() error {
	return s.conn.Close()
}

// Clear resets the session state so it can be reused from the pool.
func (s *Session) Clear() {
	s.SessionCommon.Clear()
	s.conn = nil
	s.cache = s.cache[:0]
	s.recvCount = 0
}

// recvLoop reads from the TCP connection, handling sticky packets.
func (s *Session) recvLoop() {
	buf := make([]byte, gonet.MTU)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			s.GetContext().GetEventHandler().OnClose(s, err)
			s.GetContext().RecycleSession(s)
			return
		}
		if n == 0 {
			continue
		}
		s.recvCount++

		// merge with any cached (incomplete) data from prior read
		var data []byte
		if len(s.cache) > 0 {
			data = append(s.cache, buf[:n]...)
			s.cache = s.cache[:0]
		} else {
			data = buf[:n]
		}

		// parse all complete packets from data
		for len(data) > 0 {
			msg, unused, err := s.GetContext().UnPackage(s, data)
			if err != nil {
				// incomplete packet — cache remaining bytes for next read
				s.cache = append(s.cache[:0], data...)
				break
			}
			s.GetContext().PushGlobalMessageQueue(msg)
			if unused <= 0 {
				break
			}
			data = data[len(data)-unused:]
		}
	}
}
