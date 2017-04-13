package qframe_collector_gelf

import (
	"log"
  "net"
	"fmt"
	"encoding/json"

	"github.com/qnib/qframe-types"
	"github.com/zpatrick/go-config"
)

const (
	version = "0.0.0"
)

type Plugin struct {
	QChan qtypes.QChan
	Cfg config.Config
	Name string
}

func NewPlugin(qChan qtypes.QChan, cfg config.Config, name string) Plugin {
	return Plugin{
		QChan: qChan,
		Cfg: cfg,
		Name: name,
	}
}

func (p *Plugin) Run() {
	log.Printf("[II] Start GELF collector %s v%s", p.Name, version)
	port, _ := p.Cfg.StringOr(fmt.Sprintf("collector.%s.port", p.Name), "12201")
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("[EE] %v", err)
	}
	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		log.Printf("[EE] %v", err)
	}
	defer ServerConn.Close()
	log.Printf("[II] Start GELF server on '%s'", ServerAddr)
	buf := make([]byte, 1024)

	log.Printf("[II] Wait for incomming GELF message")
	for {
		n,addr,err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("[EE] %v", err)
			continue
		}
		qm := qtypes.NewQMsg("collector", p.Name)
		gmsg := GelfMsg{}
		json.Unmarshal(buf[0:n], &gmsg)
		gmsg.SourceAddr = addr.String()
		qm.Msg = gmsg.ShortMsg
		p.QChan.Data.Send(qm)
	}
}
