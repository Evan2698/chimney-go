package configure

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Settings ..
type Settings struct {
	ServerPort uint16 `json:"server_port"`
	Server     string `json:"server"`

	LocalPort uint16 `json:"local_port"`
	Local     string `json:"local"`

	Password string `json:"password"`
	Timeout  uint32 `json:"timeout"`
}

// Parse ..
func Parse(path string) (config *Settings, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config = &Settings{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// DumpConfig ..
func DumpConfig(config *Settings) {
	log.Println("server :", config.Server)
	log.Println("server_port :", config.ServerPort)
	log.Println("local_port :", config.LocalPort)
	log.Println("local :", config.Local)
	log.Println("password :", config.Password)
	log.Println("timeout :", config.Timeout)
}
