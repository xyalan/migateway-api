package migateway

import (
	"errors"
	"time"

	"github.com/bingbaba/util/logs"
)

var (
	LOGGER      = logs.GetBlogger()
	EOF    byte = 0
)

type MIHomeManager struct {
	reportChan chan *Device

	GateWay *GateWay
	Motions map[string]*Motion

	DiscoveryTime    int64
	FreshDevListTime int64
}

func NewMiHomeManager(c *Configure) (m *MIHomeManager, err error) {
	if c == nil {
		c = DefaultConf
	}
	conn := NewConn(c)

	//connection
	err = conn.initMultiCast()
	if err != nil {
		return
	}

	//MIHomeManager
	m = &MIHomeManager{
		Motions:       make(map[string]*Motion),
		DiscoveryTime: time.Now().Unix(),
	}

	//find gateway
	m.whois(conn)

	//show device list
	gw_ip := m.GateWay.IP
	err = conn.initGateWay(gw_ip)
	if err != nil {
		return
	}

	err = m.discovery()

	//report or heartbeat message
	go func() {
		for {
			m.putDevice(<-conn.devMsgs)
		}
	}()

	return
}

func (m *MIHomeManager) putDevice(dev *Device) (added bool) {
	LOGGER.Info("DEVICESYNC:: %s(%s): %s", dev.Model, dev.Sid, dev.Data)
	gateway := m.GateWay

	var saveDev *Device
	added = true
	switch dev.Model {
	case MODEL_GATEWAY:
		gateway.Set(dev)
		saveDev = gateway.Device
	case MODEL_MOTION:
		d, found := m.Motions[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.conn = gateway.conn
			m.Motions[dev.Sid] = NewMotion(dev)
		}
		saveDev = m.Motions[dev.Sid].Device
	default:
		added = false
		LOGGER.Warn("DEVICESYNC:: unknown model is %s", dev.Model)
	}

	LOGGER.Debug("save to report chan...")
	if saveDev != nil {
		saveDev.report(true)
	}
	LOGGER.Debug("save to report chan over!")

	return
}

func (m *MIHomeManager) whois(conn *GateWayConn) {
	//read msg
	iamResp := &IamResp{}
	conn.communicate(NewWhoisRequest(), iamResp)

	//gateway information
	dev := NewGateWay(iamResp.Device)
	dev.IP = iamResp.IP
	dev.Port = iamResp.Port
	dev.conn = conn

	m.GateWay = dev
}

func (m *MIHomeManager) discovery() (err error) {
	gateway := m.GateWay
	conn := gateway.conn

	//get device list response
	LOGGER.Info("start to discover the device...")
	devListResp := &DeviceListResp{}
	if !conn.communicate(NewDevListRequest(), devListResp) {
		return errors.New("show device list error")
	}
	//gateway.setToken(devListResp.Token)
	gateway.conn.token = devListResp.Token

	//every device
	for index, sid := range devListResp.getSidArray() {
		dev := conn.waitDevice(sid)
		dev.Token = devListResp.Token
		if m.putDevice(dev) {
			LOGGER.Warn("DISCOVERY[%d]: found the device %s(%s): %v", index, dev.Model, dev.Sid, dev.Data)
		} else {
			LOGGER.Warn("DISCOVERY[%d]: unknown model %s device: %v", index, dev.Model, dev)
		}
	}
	m.DiscoveryTime = time.Now().Unix()

	return
}
