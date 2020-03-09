package socks5server

import (
	"chimney-go/socketcore"
	"context"
	"log"

	"github.com/lucas-clemente/quic-go"
)

func (s *serverHolder) runQuicServer() {
	log.Println("quic server address: ", s.ServerAddress)
	listener, err := quic.ListenAddr(s.ServerAddress, socketcore.GenerateTLSConfig(), nil)
	if err != nil {
		log.Println("Create quick socket failed", s.ServerAddress, err)
		return
	}

	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("quic accept session failed ", err)
			break
		}
		go s.serveQuicSession(session)
	}

}

func (s *serverHolder) serveQuicSession(session quic.Session) {
	defer func() {
		session.Close()
	}()

	for {
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println("quic accept stream failed", session.LocalAddr().String())
			break
		}
		v := socketcore.NewQuicSocket(session, stream)
		go s.serveOn(v)
	}
}
