package socks5client

import (
	"chimney-go/socketcore"
	"chimney-go/utils"
	"context"
	"log"
	"net"
	"strings"
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

	session, err := getQuicInstance(host)
	if err != nil {
		log.Println("create instance failed:  ", host, err)
		return nil, err
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		destoryQuicInstance()
		log.Println("create quick socket stream failed", err)
		return nil, err
	}
	log.Print("create socket(quic) socket success!")
	v := socketcore.NewQuicSocket(session, stream)
	socketcore.SetSocketTimeout(v, tm)
	return v, nil
}
