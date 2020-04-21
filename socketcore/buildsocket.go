package socketcore

import (
	"chimney-go/mobile"
	"errors"
	"log"
	"net"
	"os"
	"syscall"
)

// TCPDail for create tcp connection
func TCPDail(host string, pFun mobile.ProtectSocket) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {

		log.Println("parse tcp address failed!", host, err)
		return nil, err
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		log.Println("create tcp socket failed!!!", err)
		return nil, err
	}
	defer syscall.Close(fd)
	outcon, err := connectSocketCoreBase(fd, tcpAddr.IP, tcpAddr.Port, pFun, func(f int) error {

		err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TOS, 128)
		if err != nil {
			log.Println("set socket attributes ", err)
			return err
		}
		return nil
	})

	return outcon, nil
}

func connectSocketCoreBase(fd int,
	ip net.IP, port int,
	pFun mobile.ProtectSocket,
	funAttr func(fd int) error) (net.Conn, error) {

	sa, err := buildSocketAddress(ip, port)
	if err != nil {
		log.Println("construct network address failed!", err)
		return nil, err
	}

	if pFun != nil {
		ret := pFun.Protect(fd)
		log.Println("protect socket: ", ret)
	}

	if funAttr != nil {
		err = funAttr(fd)
		if err != nil {
			log.Println("set socket attrbuites failed!! ", err)
			return nil, err
		}
	}

	err = syscall.Connect(fd, sa)
	if err != nil {
		log.Println("connect remote end failed!!!!", err)
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "")
	defer file.Close()

	outcon, err := net.FileConn(file)
	if err != nil {
		log.Println("convert to FileConn failed:", err)
		return nil, err
	}
	return outcon, nil
}

func buildSocketAddress(ip net.IP, port int) (sa syscall.Sockaddr, err error) {

	if ip == nil {
		return nil, errors.New("none address")
	}

	if ip.To4() == nil {
		ipa := ip.To16()
		sa = &syscall.SockaddrInet6{
			Port: port,
			Addr: [16]byte{ipa[0], ipa[1], ipa[2], ipa[3],
				ipa[4], ipa[5], ipa[6],
				ipa[7], ipa[8], ipa[9],
				ipa[10], ipa[11], ipa[12],
				ipa[13], ipa[14], ipa[15]},
		}
	} else {
		ipa := ip.To4()
		sa = &syscall.SockaddrInet4{
			Port: port,
			Addr: [4]byte{ipa[0], ipa[1], ipa[2], ipa[3]},
		}
	}

	return sa, err
}

// UDPDail for android
func UDPDail(host string, pFun mobile.ProtectSocket) (net.Conn, error) {

	tcpAddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		log.Print("parse tcp address failed: ", err)
		return nil, err
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Println("create udp socket failed!!!", err)
		return nil, err
	}
	defer syscall.Close(fd)

	outcon, err := connectSocketCoreBase(fd, tcpAddr.IP, tcpAddr.Port, pFun, nil)

	return outcon, err
}
