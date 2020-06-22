package socketcore

import (
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

//QuicSocket ...
type QuicSocket interface {
	net.Conn
}

type quicSocketHolder struct {
	Session quic.Session
	Stream  quic.Stream
	destory bool
}

//NewQuicSocket ...
func NewQuicSocket(s quic.Session, socket quic.Stream) QuicSocket {
	return &quicSocketHolder{
		Session: s,
		Stream:  socket,
	}
}

//NewClientSocket ...
func NewClientSocket(s quic.Session, socket quic.Stream) QuicSocket {
	return &quicSocketHolder{
		Session: s,
		Stream:  socket,
		destory: true,
	}
}

func (s *quicSocketHolder) Read(b []byte) (n int, err error) {
	return s.Stream.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (s *quicSocketHolder) Write(b []byte) (n int, err error) {
	return s.Stream.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (s *quicSocketHolder) Close() error {
	if s.Stream != nil {
		s.Stream.Close()
	}
	if s.destory {
		if s.Session != nil {
			s.Session.CloseWithError(0x12, "ok")
		}
	}
	s.Stream = nil
	s.Session = nil
	return nil
}

// LocalAddr returns the local network address.
func (s *quicSocketHolder) LocalAddr() net.Addr {
	if s.Session != nil {
		return s.Session.LocalAddr()
	}

	return &net.TCPAddr{
		IP:   []byte("127.0.0.1"),
		Port: 1111,
	}
}

// RemoteAddr returns the remote network address.
func (s *quicSocketHolder) RemoteAddr() net.Addr {
	if s.Session != nil {
		return s.Session.RemoteAddr()
	}

	return &net.TCPAddr{
		IP:   []byte("127.0.0.1"),
		Port: 1111,
	}
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
//
// Note that if a TCP connection has keep-alive turned on,
// which is the default unless overridden by Dialer.KeepAlive
// or ListenConfig.KeepAlive, then a keep-alive failure may
// also return a timeout error. On Unix systems a keep-alive
// failure on I/O can be detected using
// errors.Is(err, syscall.ETIMEDOUT).
func (s *quicSocketHolder) SetDeadline(t time.Time) error {

	return s.Stream.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (s *quicSocketHolder) SetReadDeadline(t time.Time) error {
	return s.Stream.SetReadDeadline(t)

}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (s *quicSocketHolder) SetWriteDeadline(t time.Time) error {
	return s.Stream.SetWriteDeadline(t)
}
