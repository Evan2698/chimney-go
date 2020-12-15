package overquic

import (
	"chimney-go/configure"
	"chimney-go/utils"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
)

type Client struct {
	Session     quic.Session
	Tm          int
	Local       string
	SessionLock sync.RWMutex
	Remote      string
	Password    string
}

func NewClient(s *configure.Settings) (*Client, error) {

	addr := net.JoinHostPort(s.Server, strconv.Itoa(int(s.ServerPort)))
	c := &Client{
		Tm:       int(s.Timeout),
		Local:    net.JoinHostPort(s.Local, strconv.Itoa(int(s.LocalPort))),
		Remote:   addr,
		Password: s.Password,
	}

	se, err := c.TryGetSession()
	if err != nil {
		log.Println("create session failed! ", err)
		return nil, err
	}

	c.SessionLock.Lock()
	c.Session = se
	c.SessionLock.Unlock()
	return c, nil
}

func (c *Client) TryGetSession() (quic.Session, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{c.Password},
	}

	session, err := quic.DialAddr(c.Remote, tlsConf, nil)
	if err != nil {
		log.Println("connect server failed!", err)
		return nil, err
	}

	return session, nil
}

func (c *Client) Serve() error {

	l, err := net.Listen("tcp", c.Local)
	if err != nil {
		log.Println("listen failed ", err)
		return err
	}
	log.Println("client listen on: ", c.Local)

	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			log.Println(" accept failed ", err)
			break
		}
		utils.SetSocketTimeout(con, uint32(c.Tm))
		go c.serveOn(con)
	}

	return nil
}

func (c *Client) serveOn(con io.ReadWriteCloser) {
	c.SessionLock.RLock()
	stream, err := c.Session.OpenStreamSync(context.Background())
	c.SessionLock.RUnlock()
	if err != nil {
		log.Println("open remote stream failed,will create new stream", err)
		ss, e := c.TryGetSession()
		if e == nil {
			c.SessionLock.Lock()
			c.Session.CloseWithError(0x1, "error!")
			c.Session = ss
			c.SessionLock.Unlock()
		} else {
			con.Close()
			log.Fatal("can not Dail to server", e)
			return
		}

		c.SessionLock.RLock()
		stream, e = c.Session.OpenStreamSync(context.Background())
		c.SessionLock.RUnlock()
		if e != nil {
			con.Close()
			log.Fatal("can not Dail to server", e)
			return
		}
	}

	defer func() {
		stream.Close()
		con.Close()
	}()

	readTimeout := time.Duration(c.Tm) * time.Second
	v := time.Now().Add(readTimeout)
	stream.SetReadDeadline(v)
	stream.SetWriteDeadline(v)
	stream.SetDeadline(v)

	var wait sync.WaitGroup
	wait.Add(1)

	go func(w *sync.WaitGroup, proxy, c io.ReadWriteCloser) {

		defer w.Done()

		_, e := io.Copy(proxy, c)
		log.Println("io copy in client(1): ", e)

	}(&wait, stream, con)

	_, err = io.Copy(con, stream)
	log.Println("io copy in client(2): ", err)

	wait.Wait()
}

func (c *Client) Close() {
	c.Session.CloseWithError(12, "byte")
}
