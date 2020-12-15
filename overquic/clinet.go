package overquic

import (
	"io"
	"log"
	"sync"

	"github.com/lucas-clemente/quic-go"
)

type Client struct {
	Session quic.Session
	Tm      int
	Local   string 
}

func NewClient(s config.Settings) (*Client, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{s.Password},
	}

	addr := net.JoinHostPort(s.Server, strconv.Itoa(int(s.Port)))
	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {

		log.Println("connect server failed!", err)
		return nil, err 
	}

	return &Client{
		Session: session,
		Tm : s.Timeout
		Local : net.JoinHostPort(s.Local, strconv.Itoa(int(s.LocalPort)))
	}, nil
}


func(c *Client) Serve() error{
	l, err := net.Listen("tcp", c.Local)
	if err != nil {
		log.Println("listen failed ", err)
		return err
	}

	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			log.Println(" accept failed ", err)
			break
		}
		if s.Flag {
			log.Println("EXIT TCP")
			break
		}

		go c.serveOn(con)
	}
}

func(c *Client) serveOn( con io.ReadWriteCloser){
	
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil{
		log.Println("open remote stream failed", err)
		con.Close()
		return
	}

	defer func(){
		stream.Close()
		con.Close()		
	}()

	var wait sync.WaitGroup
	wait.Add(1)
	go func (w *sync.WaitGroup, s, c io.ReadWriteCloser) {
		 
		defer w.Done()

		_, e := io.Copy(s, c)	
		log.Println("io copy in client(1): ", e )	
		
	}(&wait,stream, con)

	_, err = io.Copy(c, s)
	log.Println("io copy in client(2): ", err )	

	wait.Wait()
}

func(c *Client) Close() {
	c.Session.CloseWithError(12, "byte")
}


