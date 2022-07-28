package main

import (
	"IpRecorder/bot"
	"IpRecorder/conf"
	"IpRecorder/cron"
	"IpRecorder/data"
	"IpRecorder/http"
	"flag"
	"log"
)

var path = flag.String("path", "./config.json", "config file path")

func main() {
	flag.Parse()
	config, err := conf.New(*path)
	if err != nil {
		log.Fatalln("Init config obj error: ", err)
	}
	dataObj := data.New(config.OnlineIpLimit)
	go func() {
		err := http.NewHttp(config, dataObj).Start()
		if err != nil {
			log.Fatalln("Start http service error: ", err)
		}
	}()
	botObj, err := bot.New(config)
	if err != nil {
		log.Fatalln("Init bot service error: ", err)
	}
	cronObj, err := cron.New(dataObj, botObj, config.HistoryIpLimit)
	if err != nil {
		log.Fatalln("Init cron error: ", err)
	}
	err = cronObj.Start()
	if err != nil {
		log.Fatalln("Start cron error: ", err)
	}
}
