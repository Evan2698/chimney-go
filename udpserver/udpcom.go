package udpserver

import (
	"bytes"
	"chimney-go/socketcore"
	"chimney-go/utils"
)

// UDPCom ...
type UDPCom struct {
	src  socketcore.Socks5Address
	dst  socketcore.Socks5Address
	cmd  uint8
	data []byte
}

//| 1 cmd| 2(len) | 1 type|  ip(domain) target | 2(len) 1 type| ip(domain) src| (3072)data|

// ParseData ..
func ParseData(in []byte) (*UDPCom, error) {
	v := &UDPCom{}

	op := bytes.NewBuffer(in)
	tmp := op.Next(1)
	v.cmd = tmp[0]
	tmp = op.Next(2)
	ll := utils.Bytes2Uint16(tmp)
	tmp = op.Next(1)
	t := tmp[0]
	ip1 := op.Next(int(ll))
	port := utils.Bytes2Uint16(ip1[len(ip1)-2:])
	tmp = ip1[:len(ip1)-2]
	vv := &socketcore.Socks5Addr{
		AddressType: t,
		Port:        port,
	}
	if t == 1 || t == 4 {
		vv.IPvX = tmp
	} else if t == 3 {
		vv.Domain = string(tmp)
	}
	v.dst = vv
	tmp = op.Next(2)
	ll = utils.Bytes2Uint16(tmp)
	tmp = op.Next(1)
	t = tmp[0]
	ip1 = op.Next(int(ll))
	port = utils.Bytes2Uint16(ip1[len(ip1)-2:])
	tmp = ip1[:len(ip1)-2]

	vv = &socketcore.Socks5Addr{
		AddressType: t,
		Port:        port,
	}
	if t == 1 || t == 4 {
		vv.IPvX = tmp
	} else if t == 3 {
		vv.Domain = string(tmp)
	}
	v.src = vv

	v.data = op.Next(op.Len())

	return v, nil
}

// 1 answer |2(len) | 1 type|  ip(domain) target | 2(len) 1 type| ip(domain) src| data(3072)

// ToAnswer ..
func ToAnswer(n *UDPCom) []byte {
	var buffer bytes.Buffer

	buffer.WriteByte(n.cmd)
	l := n.dst.GetAddressRawBytes()
	buffer.Write(utils.Uint162Bytes(uint16(len(l))))
	buffer.WriteByte(n.dst.GetAddressType())
	buffer.Write(l)
	l = n.src.GetAddressRawBytes()
	buffer.Write(utils.Uint162Bytes(uint16(len(l))))
	buffer.WriteByte(n.src.GetAddressType())
	buffer.Write(l)
	buffer.Write(n.data)

	return buffer.Bytes()
}
