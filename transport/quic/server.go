package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/flylib/gonet"
	"github.com/quic-go/quic-go"
	"math/big"
	"reflect"
)

// 接收端
type server struct {
	gonet.PeerIdentify
	ln *quic.Listener
}

func NewServer(ctx *gonet.Context) gonet.IServer {
	s := &server{}
	s.WithContext(ctx)
	return s
}

func (s *server) Listen(url string) (err error) {
	s.ln, err = quic.ListenAddr(url, s.generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	s.SetAddr(url)

	for {
		conn, err := s.ln.Accept(context.Background())
		if err != nil {
			continue
		}
		ses := newSession(s.Context, conn)
		go ses.acceptStreamLoop()
	}
}

func (s *server) Close() error {
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
