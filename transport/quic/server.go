package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/quic-go/quic-go"
	. "github.com/zjllib/gonet/v3"
	"math/big"
	"reflect"
)

// 接收端
type server struct {
	ServerIdentify
	ln *quic.Listener
}

func NewServer(addr string) *server {
	s := &server{}
	s.SetAddr(addr)
	return s
}

func (s *server) Listen() (err error) {
	s.ln, err = quic.ListenAddr(s.Addr(), s.generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		conn, err := s.ln.Accept(context.Background())
		if err != nil {
			continue
		}
		s := newSession(s.Context, conn)
		go s.acceptStream()
	}
}

func (s *server) Stop() error {
	return s.ln.Close()
}

func (s *server) SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}

// setup a bare-bones TLS config for the server
func (s *server) generateTLSConfig() *tls.Config {
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
