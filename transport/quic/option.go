package quic

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"time"
)

type Option func(*option)

type option struct {
	//specifies the duration for the handshake to complete.Default is 5 second
	HandshakeTimeout time.Duration

	//Channel id modulo
	modulo uint32
	mtu    int
}

func WithHandshakeTimeout(t time.Duration) Option {
	return func(option *option) {
		option.HandshakeTimeout = t
	}
}

// Channel id modulo
func WithChannelIdModulo(mod uint32) Option {
	return func(option *option) {
		if mod == 0 {
			return
		}
		option.modulo = mod
	}
}

// set maximum transmission unit
func WithMaximumTransmissionUnit(unit int) Option {
	return func(o *option) {
		o.mtu = unit
	}
}

// setup a bare-bones TLS config for the server
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
