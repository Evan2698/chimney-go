package configure

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// AppConfig ..
type AppConfig struct {
	ServerPort uint16 `json:"server_port"`
	Server     string `json:"server"`
	QuicPort   uint16 `json:"quic_port"`
	UDPPort    uint16 `json:"udp_port"`
	LocalPort  uint16 `json:"local_port"`
	Local      string `json:"local"`
	Which      string `json:"which"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Timeout    uint32 `json:"timeout"`
}

// Parse ..
func Parse(path string) (config *AppConfig, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config = &AppConfig{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// DumpConfig ..
func DumpConfig(config *AppConfig) {
	log.Println("server :", config.Server)
	log.Println("server_port :", config.ServerPort)
	log.Println("QuicPort :", config.QuicPort)
	log.Println("udpport :", config.UDPPort)
	log.Println("local_port :", config.LocalPort)
	log.Println("local :", config.Local)
	log.Println("which :", config.Which)
	log.Println("method :", config.Method)
	log.Println("password :", config.Password)
	log.Println("timeout :", config.Timeout)
}
