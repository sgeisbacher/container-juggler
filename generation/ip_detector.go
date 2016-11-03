package generation

import (
	"log"
	"net"
)

type UplinkIPDetector struct{}

func (ipd UplinkIPDetector) Detect() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	switch a := conn.LocalAddr().(type) {
	case *net.IPAddr:
		return a.IP
	case *net.TCPAddr:
		return a.IP
	case *net.UDPAddr:
		return a.IP
	case *net.IPNet:
		return a.IP
	}

	return nil
}
