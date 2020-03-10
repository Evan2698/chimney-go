package socks5client

import (
	"chimney-go/socketcore"
	"chimney-go/utils"
	"context"
	"crypto/tls"
	"log"
	"net"
	"strings"

	quic "github.com/lucas-clemente/quic-go"
)

func buildGeneralSocket(host, network string, tm uint32) (con net.Conn, err error) {
	defer utils.Trace("buildGeneralSocket")()

	log.Println("function: ", host, network)

	if strings.Contains("quic", network) {
		return buildQuicSocket(host, network, tm)
	}

	log.Println("builcConnect: ", host)
	con, err = net.Dial("tcp", host)
	if err == nil {
		socketcore.SetSocketTimeout(con, tm)
	}
	return con, err
}

func buildQuicSocket(host, network string, tm uint32) (con net.Conn, err error) {
	log.Println("quic connect:  ", host, network)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{socketcore.QuicProtocolName},
	}

	session, err := quic.DialAddr(host, tlsConf, &quic.Config{})
	if err != nil {
		log.Println("create quick socket session failed", err)
		return nil, err
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		session.CloseWithError(0, "failed")
		log.Println("create quick socket stream failed", err)
		return nil, err
	}
	log.Print("create socket(quic) socket success!")
	v := socketcore.NewClientSocket(session, stream)
	socketcore.SetSocketTimeout(v, tm)
	return v, nil
}
