package data

import (
	"strings"
	"sync"
	"time"
)

type IpStore struct {
	Ip map[string]int64
	Mu sync.RWMutex
}

type Data struct {
	IpHistory    *sync.Map // map[string]string
	UserOnlineIp *sync.Map // map[string]*IpStore
	ipLimit      int
}

func New(limit int) *Data {
	return &Data{
		IpHistory:    &sync.Map{},
		UserOnlineIp: &sync.Map{},
		ipLimit:      limit,
	}
}

type UserIps struct {
	Uid int      `json:"Uid"`
	Ips []string `json:"Ips"`
}

func (p *Data) AddIpHistory(ipHistory *[]UserIps) {
	for i := range *ipHistory {
		if oldIp, ok := p.IpHistory.LoadOrStore((*ipHistory)[i].Uid, strings.Join((*ipHistory)[i].Ips, ",")); ok {
			for _, ip := range (*ipHistory)[i].Ips {
				if !strings.Contains(oldIp.(string), ip) {
					p.IpHistory.Store((*ipHistory)[i].Uid, oldIp.(string)+","+ip)
				}
			}
		}
	}
}

func (p *Data) ClearAllUsers() {
	p.IpHistory = &sync.Map{}
}

func (p *Data) SyncUserOnlineIP(userIps *[]UserIps) *[]UserIps {
	var ips *IpStore
	for i := range *userIps {
		if tmp, ok := p.UserOnlineIp.Load((*userIps)[i].Uid); ok {
			ips = tmp.(*IpStore)
			for _, ip := range (*userIps)[i].Ips {
				ips.Mu.Lock()
				if len(ips.Ip) >= p.ipLimit {
					ips.Mu.Unlock()
					break
				}
				ips.Ip[ip] = time.Now().Unix()
				ips.Mu.Unlock()
			}
		} else {
			ips = &IpStore{Ip: make(map[string]int64), Mu: sync.RWMutex{}}
			count := 0
			for _, ip := range (*userIps)[i].Ips {
				if count >= p.ipLimit {
					break
				}
				ips.Ip[ip] = time.Now().Unix()
				count++
			}
			p.UserOnlineIp.Store((*userIps)[i].Uid, ips)
		}
	}
	userIps = &[]UserIps{}
	p.UserOnlineIp.Range(func(user, ipList interface{}) bool {
		ips = ipList.(*IpStore)
		var ip []string
		ips.Mu.RLock()
		for i := range ips.Ip {
			ip = append(ip, i)
		}
		ips.Mu.RUnlock()
		*userIps = append(*userIps, UserIps{Uid: user.(int), Ips: ip})
		return true
	})
	p.ClearUserOnlineIP()
	return userIps
}

func (p *Data) ClearUserOnlineIP() {
	p.UserOnlineIp.Range(func(k, v interface{}) bool {
		ip := v.(*IpStore)
		ip.Mu.Lock()
		for i := range ip.Ip {
			if ip.Ip[i] <= time.Now().Unix()-120 {
				delete(ip.Ip, i)
			}
		}
		if len(ip.Ip) == 0 {
			p.UserOnlineIp.Delete(k)
		}
		ip.Mu.Unlock()
		return true
	})
}
