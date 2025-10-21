//go:build android || ios || mobile_skel

package transport

import "context"

type TransportStatus struct {
	RTTms   int64
	Remote  string
	ALPN    string
	SNI     string
	LastErr string
}

type Transport interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status() TransportStatus
}
