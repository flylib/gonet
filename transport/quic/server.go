package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	. "github.com/zjllib/gonet/v3"
	"math/big"
)

//接收端
type server struct {
	ServerIdentify
	ln quic.Listener
}

func (s *server) Start() (err error) {
	s.ln, err = quic.ListenAddr(s.Addr(), generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		conn, err := s.ln.Accept(context.Background())
		if err != nil {
			continue
		}
		go s.newConn(conn)
	}
}

func (s *server) Stop() error {
	// TODO 关闭处理
	return s.ln.Close()
}

//新连接
func (s *server) newConn(conn quic.Connection) {
	ses := newSession(conn)
	ses.recvStreamLoop()
	//ses.recvLoop()
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
