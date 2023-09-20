package proxyserver

import (
	"net"
	"time"
)

type ProxyCon interface {
	net.Conn
}

type ProxyHub struct {
	tunnel net.Conn
}

func (t *ProxyHub) Read(buf []byte) (n int, err error) {
	return t.tunnel.Read(buf)
}

func (t *ProxyHub) Write(b []byte) (n int, err error) {

	return t.tunnel.Write(b)
}

func (t *ProxyHub) Close() error {

	return t.tunnel.Close()
}

// LocalAddr returns the local network address, if known.
func (t *ProxyHub) LocalAddr() net.Addr {
	return t.tunnel.LocalAddr()
}

// RemoteAddr returns the remote network address, if known.
func (t *ProxyHub) RemoteAddr() net.Addr {
	return t.tunnel.RemoteAddr()
}

func (t *ProxyHub) SetDeadline(tt time.Time) error {
	return t.tunnel.SetDeadline(tt)
}

// SetReadDeadline sets the deadline for future Read calls

func (t *ProxyHub) SetReadDeadline(tt time.Time) error {
	return t.tunnel.SetReadDeadline(tt)
}

// SetWriteDeadline sets the deadline for future Write calls
func (t *ProxyHub) SetWriteDeadline(tt time.Time) error {
	return t.tunnel.SetWriteDeadline(tt)
}
