package socks5server

import (
	"chimney-go/socketcore"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

func (s *serverHolder) runQuicServer() {
	log.Println("quic server address: ", s.ServerAddress)
	listener, err := quic.ListenAddr(s.ServerAddress, generateTLSConfig(), nil)
	if err != nil {
		log.Println("Create quick socket failed", s.ServerAddress, err)
		return
	}

	for {
		_, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("quic accept session failed ", err)
			break
		}
		go s.serveQuicSession()
	}

}

// GenerateTLSConfig ...
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
		NextProtos:   []string{socketcore.QuicProtocolName},
	}
}

func (s *serverHolder) serveQuicSession() {
	// defer func() {
	// 	session.CloseWithError(0x12, "ok")
	// }()
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Println(" fatal error on serveQuicSession: ", err)
	// 	}
	// }()

	// for {
	// 	stream, err := session.AcceptStream(context.Background())
	// 	if err != nil {
	// 		log.Println("quic accept stream failed", session.LocalAddr().String())
	// 		break
	// 	}
	// 	v := socketcore.NewQuicSocket(session, stream)
	// 	go s.serveOn(v)
	// }
}
