package socks5client

import (
	"sync"
)

type singleHolder struct {
	lock sync.Mutex
	//instance quic.Session
}

var ks = &singleHolder{}

func getQuicInstance(host string) error {
	//tlsConf := &tls.Config{
	//	InsecureSkipVerify: true,
	//	NextProtos:         []string{socketcore.QuicProtocolName},
	//}

	// var err error

	// if ks.instance == nil {
	// 	ks.lock.Lock()
	// 	if ks.instance == nil {
	// 		ks.instance, err = quic.DialAddr(host, tlsConf, &quic.Config{})
	// 	}
	// 	ks.lock.Unlock()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// return ks.instance, err
	return nil
}

func destoryQuicInstance() {
	// if ks.instance != nil {
	// 	ks.lock.Lock()
	// 	if ks.instance != nil {
	// 		ks.instance.CloseWithError(0x12, "ok")
	// 		ks.instance = nil
	// 	}
	// 	ks.lock.Unlock()
	// }
}
