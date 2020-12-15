package main

import (
	"chimney-go/configure"
	"flag"
	"fmt"
	"log"

	"os"
	"path/filepath"
	"runtime"
)

var s *bool

func main() {
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

	if *s {
		// start quic server
		log.Println("I AM A SERVER!!")

	} else {
		log.Println("I AM A CLIENT!!")

	}
}
