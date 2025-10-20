//go:build android || ios

package mobile

import (
	"context"
	"errors"
	"net"
	"syscall"
)

func protectedPacketConn(ctx context.Context) (net.PacketConn, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var ctrlErr error
			if err := c.Control(func(fd uintptr) {
				if !protectFD(int(fd)) {
					ctrlErr = errors.New("protect failed")
				}
			}); err != nil {
				return err
			}
			return ctrlErr
		},
	}
	return lc.ListenPacket(ctx, "udp", "0.0.0.0:0")
}

func protectedTCPDialer() *net.Dialer {
	return &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			var ctrlErr error
			if err := c.Control(func(fd uintptr) {
				if !protectFD(int(fd)) {
					ctrlErr = errors.New("protect failed")
				}
			}); err != nil {
				return err
			}
			return ctrlErr
		},
	}
}
