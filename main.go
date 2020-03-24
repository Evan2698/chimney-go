package main

import (
	"chimney-go/configure"
	"chimney-go/privacy"
	"chimney-go/socketcore"
	"chimney-go/socks5server"
	"chimney-go/utils"
	"flag"
	"fmt"
	"log"
	"net"

	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var s *bool

func main() {
	/*	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()*/
	var configpath string
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu * 4)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Print("can not combin config file path!")
		os.Exit(1)
	}
	configpath = dir + "/config.json"
	if (len(configpath)) == 0 {
		fmt.Println("config file path is incorrect!!", configpath)
		os.Exit(1)
	}

	config, err := configure.Parse(configpath)
	if err != nil {
		fmt.Println("load config file failed!", err)
		os.Exit(1)
	}
	//configure.DumpConfig(config)

	s = flag.Bool("s", false, "a bool")
	flag.Parse()

	user := privacy.BuildMacHash(privacy.MakeCompressKey(config.Password), "WhereRU")

	serverTCP := net.JoinHostPort(config.Server, strconv.Itoa(int(config.ServerPort)))
	serverQuic := net.JoinHostPort(config.Server, strconv.Itoa(int(config.QuicPort)))

	which := utils.FormatProtocol(config.Which)

	log.Println("protol ", which)

	if *s {
		// start quic server
		log.Println("I AM A SERVER!!")
		go func(){
			sconf := &socks5server.SConfig{
				ServerAddress: serverQuic,
				Network:       "quic",
				Tm:            config.Timeout,
				User:          user,
				Pass:          user,
				Key:           privacy.MakeCompressKey(config.Password),
				I:             privacy.NewMethodWithName(config.Method),
			}
			ss := socks5server.NewServer(sconf)
			ss.Serve()

		}():

		// start tcp server 
		sconf := &socks5server.SConfig{
			ServerAddress: serverTCP,
			Network:       "tcp",
			Tm:            config.Timeout,
			User:          user,
			Pass:          user,
			Key:           privacy.MakeCompressKey(config.Password),
			I:             privacy.NewMethodWithName(config.Method),
		}
		ss := socks5server.NewServer(sconf)
		ss.Serve()

	} else {
		log.Println("I AM A CLIENT!!")

		settings := &socketcore.ClientConfig{
			User:    user,
			Pass:    user,
			Key:     privacy.MakeCompressKey(config.Password),
			Proxy:   serverHost,
			Tm:      config.Timeout,
			Network: which,
		}
		sconf := &socks5server.SConfig{
			ServerAddress: net.JoinHostPort(config.Local, strconv.Itoa(int(config.LocalPort))),
			Network:       "tcp",
			CC:            settings,
			Tm:            config.Timeout,
		}
		ss := socks5server.NewServer(sconf)
		ss.Serve()

	}
}
