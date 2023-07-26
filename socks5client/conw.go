package socks5client

import (
	"chimney-go/mobile"
	"chimney-go/socketcore"
	"chimney-go/utils"
	"log"
	"net"
)

func buildGeneralSocket(host, network string, tm uint32, profect mobile.ProtectSocket) (con net.Conn, err error) {
	defer utils.Trace("buildGeneralSocket")()

	log.Println("function: ", host, network)

	log.Println("builcConnect: ", host)
	if profect != nil {
		con, err = socketcore.TCPDail(host, profect)
	} else {
		con, err = net.Dial("tcp", host)
	}
	if err == nil {
		socketcore.SetSocketTimeout(con, tm)
	}
	return con, err
}
