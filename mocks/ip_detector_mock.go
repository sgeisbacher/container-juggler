package mocks

import "net"

type IPDetectorMock struct {
	DetectCall struct {
		Returns net.IP
	}
}

func (ipd IPDetectorMock) Detect() net.IP {
	return ipd.DetectCall.Returns
}
