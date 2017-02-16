package mocks

import "net"

// IPDetectorMock mock-impl of IPDetector-interface
type IPDetectorMock struct {
	DetectCall struct {
		Returns net.IP
	}
}

// Detect returns DetectCall.Returns on every call
func (ipd IPDetectorMock) Detect() net.IP {
	return ipd.DetectCall.Returns
}
