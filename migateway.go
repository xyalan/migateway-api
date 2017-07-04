package migateway

const (
	MODEL_GATEWAY     = "gateway"
	FIELD_GATEWAY_RGB = "rgb"
	FIELD_IP          = "ip"

	FLASHING_WEIGHT_WEAK   = "1"
	FLASHING_WEIGHT_NORMAL = "2"
	FLASHING_WEIGHT_STRONG = "3"
)

//GateWay Status
type GateWay struct {
	*Device
	IP       string
	Port     string
	lastRGB  uint32
	RGB      uint32
	callBack func(gw *GateWay) error
}

func NewGateWay(dev *Device) *GateWay {
	dev.ReportChan = make(chan interface{}, 1)
	g := &GateWay{Device: dev}
	g.Set(dev)
	return g
}

func (g *GateWay) Set(dev *Device) {
	if dev.hasFiled(FIELD_IP) {
		g.IP = dev.GetData(FIELD_IP)
	}
	if dev.hasFiled(FIELD_GATEWAY_RGB) {
		if g.RGB != 0 {
			LOGGER.Warn("Save Last RGB:%d", g.RGB)
			g.lastRGB = g.RGB
		}
		g.RGB = dev.GetDataAsUInt32(FIELD_GATEWAY_RGB)
	}
	if dev.Token != "" {
		g.setToken(dev.Token)
	}
	if dev.ShortID > 0 {
		g.ShortID = dev.ShortID
	}
	if g.callBack != nil {
		err := g.callBack(g)
		if err != nil {
			LOGGER.Error("exec callback error: %#v", err)
		}
	}
}

func (g *GateWay) setToken(token string) {
	g.Token = token
	g.conn.token = token
}
