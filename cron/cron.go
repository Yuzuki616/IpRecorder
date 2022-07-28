package cron

import (
	"IpRecorder/bot"
	"IpRecorder/data"
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
	cron2 "github.com/robfig/cron/v3"
	"strconv"
	"strings"
)

type Cron struct {
	data         *data.Data
	bot          *bot.Bot
	historyLimit int
	ip           *ip2location.DB
}

func New(data *data.Data, bot *bot.Bot, limit int) (*Cron, error) {
	ip, err := ip2location.OpenDB("./ipdata.bin")
	if err != nil {
		return nil, err
	}
	return &Cron{
		historyLimit: limit,
		ip:           ip,
		data:         data,
		bot:          bot,
	}, nil
}

func (p *Cron) checkUserIpList() {
	p.data.IpHistory.Range(func(user, ip interface{}) bool {
		ips := strings.Split(ip.(string), ",")
		if len(ips) > p.historyLimit {
			ipAndRegions := make(map[string]string, len(ips))
			for _, ip := range ips {
				location, err := p.ip.Get_city(ip)
				if err != nil {
					continue
				}
				ipAndRegions[location.City] = ip
			}
			if len(ipAndRegions) > p.historyLimit {
				msg := "IP列表: \n"
				for region := range ipAndRegions {
					msg += "\n" + ipAndRegions[region] + " | " + region
				}
				err := p.bot.PushMsgToMaster("历史连接IP数超出限制通知\n\n用户: " + strconv.Itoa(user.(int)) + msg)
				if err != nil {
					fmt.Println("Push message error: ", err)
				}
			}
			err := p.bot.PushMsgToMaster("历史连接IP数超出限制通知\n\n用户: " + strconv.Itoa(user.(int)) +
				"\nIP: " + strings.Join(ips, " | "))
			if err != nil {
				fmt.Println("Push message error: ", err)
			}
		}
		return true
	})
}

func (p *Cron) Start() error {
	c := cron2.New()
	if p.historyLimit > 0 {
		_, err := c.AddFunc("0 0 * * *", func() {
			p.checkUserIpList()
			p.data.ClearAllUsers()
		})
		if err != nil {
			return fmt.Errorf("add write to file task error: %v", err)
		}
	}
	_, err := c.AddFunc("*/1 * * * *", func() {
		p.data.ClearUserOnlineIP()
	})
	if err != nil {
		return fmt.Errorf("add clear user online ip task error: %v", err)
	}
	c.Run()
	return nil
}
