package socketcore

import (
	"chimney-go/privacy"
	"chimney-go/utils"
	"errors"
	"log"
	"net"
)

// ISocket ...
type ISocket interface {
	Read() ([]byte, error)
	Write(b []byte) error
	Close()
}

type iSocketHolder struct {
	EChannel net.Conn
	Key      []byte
	I        privacy.EncryptThings
	inBuffer []byte
}

const (
	offset = 512
)

// NewISocket ...
func NewISocket(con net.Conn, i privacy.EncryptThings, key []byte) ISocket {
	return &iSocketHolder{
		EChannel: con,
		Key:      key,
		I:        i,
		inBuffer: Alloc(),
	}
}

func readbytesfromraw(bytes uint32, buffer []byte, con net.Conn) ([]byte, error) {

	if bytes <= 0 {
		return nil, errors.New("0 bytes can not read! ")
	}

	var index uint32
	var err error
	var n int
	for {
		n, err = con.Read(buffer[index:])
		log.Println("read from socket size: ", n, err)
		index = index + uint32(n)
		if err != nil {
			log.Println("error on read_bytes_from_socket ", n, err)
			break
		}

		if index >= bytes && index > 0 {
			log.Println("read count for output ", index, err)
			break
		}

	}

	if index < bytes && index != 0 {
		log.Println("can not run here!!!!!")
	}

	log.Println("read result size: ", index, err)
	return buffer[:bytes], err
}

func (s *iSocketHolder) Close() {
	log.Println("CLOSE*********")
	if s.inBuffer != nil {
		Free(s.inBuffer)
		s.inBuffer = nil
	}
}

func (s *iSocketHolder) Read() ([]byte, error) {
	defer utils.Trace("iSocketHolder.Read")()

	buffer, err := readbytesfromraw(4, s.inBuffer[:4], s.EChannel)
	if err != nil {
		log.Println("read raw content failed", err)
		return nil, err
	}

	vLen := utils.Bytes2Int(buffer)
	log.Println("read: ", buffer, vLen)

	if vLen > pageSize {
		log.Println("content length is too long", vLen)
		return nil, errors.New("Length is too long")
	}

	buffer, err = readbytesfromraw(vLen, s.inBuffer[:vLen], s.EChannel)
	if err != nil {
		log.Println("read content failed: ", err)
		return nil, err
	}
	out, err := s.I.Uncompress(buffer, s.Key)
	if err != nil {
		log.Println("uncompress failed: ", err)
		return nil, err
	}

	log.Println("Uncompress rEAD: ", vLen)
	return out, nil

}

func (s *iSocketHolder) Write(b []byte) error {
	defer utils.Trace("iSocketHolder.Write")()

	out, err := s.I.Compress(b, s.Key)
	if err != nil {
		log.Println("zip content failed: ", err)
		return err
	}
	oLen := len(out)

	if oLen+4 > pageSize {
		log.Println("out of memory!!", oLen+4)
		return errors.New("out of memory")
	}

	vLenBuffer := utils.Int2Bytes(uint32(oLen))
	_, err = s.EChannel.Write(vLenBuffer)
	if err != nil {
		log.Println("write length of content failed: ", err)
		return err
	}
	_, err = s.EChannel.Write(out)
	if err != nil {
		log.Println("write content failed: ", err)
		return err
	}

	return nil

}
