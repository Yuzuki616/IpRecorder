package bot

import (
	"IpRecorder/conf"
	"fmt"
	teleBot "gopkg.in/telebot.v3"
	"time"
)

type Bot struct {
	client   *teleBot.Bot
	masterId int64
}

func New(c *conf.Conf) (*Bot, error) {
	b, newErr := teleBot.NewBot(teleBot.Settings{
		Token:  c.BotToken,
		Poller: &teleBot.LongPoller{Timeout: 10 * time.Second},
	})
	if newErr != nil {
		return nil, fmt.Errorf("init telegram teleBot error: %v", newErr)
	}
	return &Bot{
		client:   b,
		masterId: c.MasterId,
	}, nil
}

func (p *Bot) PushMsgToMaster(msg string) error {
	c, err := p.client.ChatByID(p.masterId)
	if err != nil {
		return fmt.Errorf("get chat error: %v", err)
	}
	_, err = p.client.Send(c, msg)
	if err != nil {
		return fmt.Errorf("send message error: %v", err)
	}
	return nil
}
