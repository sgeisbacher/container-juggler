package generation

import (
	"log"
	"net"
)

// UplinkIPDetector defines a basic struct
type UplinkIPDetector struct{}

// Detect detects local-outbound ip-address
func (ipd UplinkIPDetector) Detect() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

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
