package migateway

import (
	"time"
)

var (
	DefaultConf = &Configure{
		WhoIsTimeOut:         3,
		WhoIsRetry:           5,
		DevListTimeOut:       3,
		DevListRetry:         5,
		ReadTimeout:          3,
		ReadRetry:            1,
		ReportForwardTimeout: 1,
		AESKey:               "",
	}
)

type Configure struct {
	WhoIsTimeOut         int
	WhoIsRetry           int
	DevListTimeOut       int
	DevListRetry         int
	ReadTimeout          int
	ReadRetry            int
	ReportForwardTimeout int

	AESKey string
}

func NewConfig() *Configure {
	return &Configure{
		WhoIsTimeOut:         3,
		WhoIsRetry:           5,
		DevListTimeOut:       3,
		DevListRetry:         5,
		ReadTimeout:          3,
		ReadRetry:            1,
		ReportForwardTimeout: 1,
		AESKey:               "",
	}
}

func (c *Configure) getRetryAndTimeout(req *Request) (int, time.Duration) {
	if req.Cmd == CMD_WHOIS {
		return c.WhoIsRetry, time.Duration(c.WhoIsTimeOut) * time.Second
	} else if req.Cmd == CMD_DEVLIST {
		return c.DevListRetry, time.Duration(c.DevListTimeOut) * time.Second
	} else {
		return c.ReadRetry, time.Duration(c.ReadTimeout) * time.Second
	}
}

func (c *Configure) SetAESKey(key string) {
	c.AESKey = key
}
