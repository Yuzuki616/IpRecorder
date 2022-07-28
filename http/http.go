package http

import (
	"IpRecorder/conf"
	"IpRecorder/data"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Http struct {
	data          *data.Data
	addr          string
	token         string
	historyLimit  int
	onlineIpLimit int
	gin           *gin.Engine
}

func NewHttp(c *conf.Conf, data *data.Data) *Http {
	gin.SetMode(gin.ReleaseMode)
	return &Http{
		data:         data,
		addr:         c.Addr,
		token:        c.Token,
		historyLimit: c.HistoryIpLimit,
		gin:          gin.Default(),
	}
}

func (p *Http) initRoute() {
	p.gin.POST("/api/v1/SyncOnlineIp", func(context *gin.Context) {
		token := context.Query("token")
		if token == p.token {
			var userIp []data.UserIps
			if err := context.BindJSON(&userIp); err != nil {
				context.JSON(400, gin.H{
					"error": fmt.Sprintf("bind error: %v", err),
				})
				return
			}
			if p.historyLimit > 0 {
				p.data.AddIpHistory(&userIp)
			}
			if p.onlineIpLimit > 0 {
				ip := p.data.SyncUserOnlineIP(&userIp)
				context.JSON(200, ip)
			} else {
				context.String(200, "ok")
			}
		} else {
			context.Status(403)
			return
		}
	})
}

func (p *Http) Start() error {
	p.initRoute()
	err := p.gin.Run(p.addr)
	if err != nil {
		return fmt.Errorf("http start error: %v", err)
	}
	return nil
}
