package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conf struct {
	path           string
	Addr           string `json:"Addr"`
	Token          string `json:"Token"`
	IpDb           string `json:"IpDb"`
	MasterId       int64  `json:"MasterId"`
	BotToken       string `json:"BotToken"`
	HistoryIpLimit int    `json:"HistoryIpLimit"`
	OnlineIpLimit  int    `json:"OnlineIpLimit"`
}

func New(path string) (*Conf, error) {
	c := &Conf{
		path:           path,
		Addr:           "127.0.0.1:1211",
		Token:          "token",
		IpDb:           "./IP2LOCATION-LITE-DB3.BIN",
		MasterId:       123,
		BotToken:       "token",
		HistoryIpLimit: 3,
		OnlineIpLimit:  3,
	}
	return c, nil
}

func (c *Conf) LoadConfig() error {
	f, err := os.Open(c.path)
	if err != nil {
		return fmt.Errorf("open config file error: %v", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return fmt.Errorf("decode config error: %v", err)
	}
	return nil
}

func (c *Conf) SaveConfig() error {
	f, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("create config file error: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		return fmt.Errorf("encode config error: %v", err)
	}
	return nil
}
