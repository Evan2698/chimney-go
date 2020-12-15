package overquic

import (
	"bytes"
	"chimney-go/common"
	"chimney-go/utils"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/lucas-clemente/quic-go"
)

const (
	socks5Version          uint8 = 0x5
	socks5NoAuth           uint8 = 0x0
	socks5AuthWithUserPass uint8 = 0x2
	socks5ReplySuccess     uint8 = 0x0
)

const (
	socks5CMDConnect uint8 = 0x1
	socks5CMDBind    uint8 = 0x2
	socks5CMDUDP     uint8 = 0x3
)

const (
	socks5AddressIPV4   uint8 = 0x1
	socks5AddressIPV6   uint8 = 0x4
	socks5AddressDomain uint8 = 0x3
)

func LaunchServer(address, password string) error {

	listener, err := quic.ListenAddr(address, common.GenerateTLSConfig(password), nil)
	if err != nil {
		return err
	}

	defer listener.Close()
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("session listen failed!!!", err)
			break
		}
		go handleSession(sess)
	}

	log.Println("server exit!!!")

	return nil
}

func handleSession(s quic.Session) {

	defer s.CloseWithError(0x34, "Error ocurred!!!")
	for {
		stream, err := s.AcceptStream(context.Background())
		if err != nil {
			log.Println("Accept Stream failed!!!", err)
			break
		}

		go handleServe(stream)
	}

}

func handleServe(s quic.Stream) {

	defer s.Close()

	if echoHello(s) != nil {

		return
	}

	target, err := handleConnectCommand(s)
	if err != nil || target == nil {
		return
	}

	defer target.Close()

	var wait sync.WaitGroup
	wait.Add(1)
	go func(w *sync.WaitGroup, s, c io.ReadWriteCloser) {

		defer w.Done()

		_, e := io.Copy(s, c)
		log.Println("io copy in server(1): ", e)

	}(&wait, s, target)

	_, err = io.Copy(target, s)
	log.Println("io copy in server(2): ", err)

	wait.Wait()

}

func echoHello(conn io.ReadWriteCloser) error {
	defer utils.Trace("echoHello")()
	tmpBuffer := common.Alloc()
	defer common.Free(tmpBuffer)
	n, err := conn.Read(tmpBuffer)
	if err != nil {
		log.Print("read hello failed:", err)
		return err
	}
	log.Println(" F: ", tmpBuffer[:n])

	if n < 2 || tmpBuffer[0] != socks5Version {
		log.Println("server protocol format is incorrect : ", tmpBuffer[:n])
		res := []byte{socks5Version, 0xff}
		conn.Write(res)
		return errors.New("server protocol format is incorrect")
	}

	welcome := []byte{socks5Version, socks5NoAuth}
	log.Println(" Welcome: ", welcome)
	_, err = conn.Write(welcome)
	log.Println("reply to remote no auth:", err)
	return err
}

func handleConnectCommand(conn io.ReadWriteCloser) (io.ReadWriteCloser, error) {
	tmpBuffer := common.Alloc()
	defer common.Free(tmpBuffer)

	n, err := conn.Read(tmpBuffer)
	if err != nil {
		conn.Write([]byte{0x05, 0x0A, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		log.Println("read connect command failed", err)
		return nil, err
	}
	log.Println("cmd: ", tmpBuffer[:n])

	cmd := tmpBuffer[:n]
	if len(cmd) < 4 {
		conn.Write([]byte{0x05, 0x0A, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		log.Println("cmd length is too short!!")
		return nil, errors.New("cmd length is too short")
	}
	if tmpBuffer[0] != socks5Version {
		conn.Write([]byte{0x05, 0x0A, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		log.Println("cmd protocol is incorrect")
		return nil, errors.New("cmd protocol is incorrect")
	}
	var target io.ReadWriteCloser
	switch cmd[1] {
	case socks5CMDConnect:
		target, err = responseCommandConnect(conn, cmd)
		if err != nil {
			conn.Write([]byte{0x05, 0x0A, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			log.Println("handleConnect", err)
			return nil, errors.New("connect failed")
		}

	case socks5CMDBind:
		fallthrough
	case socks5CMDUDP:
		fallthrough
	default:
		conn.Write([]byte{0x05, 0x0B, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		log.Println("Not Support this CMD", cmd)
		return nil, errors.New("Not Support this CMD")
	}

	return target, nil

}

func responseCommandConnect(conn io.ReadWriteCloser, cmd []byte) (io.ReadWriteCloser, error) {

	content := cmd[4:]
	port := utils.Bytes2Uint16(content[len(content)-2:])
	var addr string

	if cmd[3] == socks5AddressIPV4 || cmd[3] == socks5AddressIPV6 {
		addr = net.IP(content[:len(content)-2]).String()
	} else if cmd[3] == socks5AddressDomain {
		addr = string(content[1 : len(content)-2])
	} else {
		log.Println("connect command is incorrect", cmd)
		return nil, errors.New("connect command is incorrect")
	}

	host := net.JoinHostPort(addr, strconv.Itoa(int(port)))

	log.Println("Connect Address:", host)

	target, err := net.Dial("tcp", host)
	if err != nil {
		log.Println("Dial host failed! ", host)
		return nil, err
	}

	peerAddress := target.LocalAddr().String()
	s, p, err := net.SplitHostPort(peerAddress)
	if err != nil {
		target.Close()
		return nil, err
	}

	np, _ := strconv.Atoi(p)
	port = uint16(np)
	atype := 0x1
	IPvX := net.ParseIP(s)
	if IPvX == nil {
		atype = 0x3
	}
	if IPvX.To4() == nil {
		atype = 0x4
	}

	var op bytes.Buffer
	op.Write([]byte{socks5Version, socks5ReplySuccess, 0x00, byte(atype)})
	if atype == 0x3 {
		op.WriteByte(byte(len(s)))
	} else if atype == 0x1 {
		op.Write(IPvX.To4())
	} else {
		op.Write(IPvX)
	}
	op.Write(utils.Port2Bytes(port))
	conn.Write(op.Bytes())
	return target, nil
}
