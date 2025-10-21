//go:build mobile_skel

package socks

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/proxy"
)

func startEchoServer(t *testing.T) (net.Listener, string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("echo listen: %v", err)
	}
	addr := ln.Addr().String()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				r := bufio.NewReader(conn)
				line, _ := r.ReadString('\n')
				conn.Write([]byte(line)) // echo
			}(c)
		}
	}()
	return ln, addr
}

func TestLocalSocks_EchoRoundTrip_Metrics(t *testing.T) {
	resetState()
	ResetBytesStats()

	// echo
	echoLn, echoAddr := startEchoServer(t)
	defer echoLn.Close()

	// socks на random-порту
	if err := StartLocalSocks("127.0.0.1", 0); err != "" {
		t.Fatalf("StartLocalSocks: %v", err)
	}
	defer StopLocalSocks()

	socks := LocalSocksAddr()
	if !strings.HasPrefix(socks, "127.0.0.1:") {
		t.Fatalf("bad socks addr: %s", socks)
	}

	// клиент через SOCKS
	dialer, err := proxy.SOCKS5("tcp", socks, nil, proxy.Direct)
	if err != nil {
		t.Fatalf("SOCKS5 dialer: %v", err)
	}
	conn, err := dialer.Dial("tcp", echoAddr)
	if err != nil {
		t.Fatalf("dial via socks: %v", err)
	}
	defer conn.Close()

	msg := "ping-ping\n"
	if _, err := conn.Write([]byte(msg)); err != nil {
		t.Fatalf("write: %v", err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	buf := make([]byte, len(msg))
	if _, err := conn.Read(buf); err != nil {
		t.Fatalf("read: %v", err)
	}

	// проверяем HealthJSON
	var h Health
	if err := json.Unmarshal([]byte(HealthJSON()), &h); err != nil {
		t.Fatalf("HealthJSON invalid: %v", err)
	}
	if h.BytesOut < uint64(len(msg)) {
		t.Fatalf("BytesOut=%d, want >= %d", h.BytesOut, len(msg))
	}
	if h.BytesIn < uint64(len(msg)) {
		t.Fatalf("BytesIn=%d, want >= %d", h.BytesIn, len(msg))
	}
	if h.Reconnects < 1 {
		t.Fatalf("Reconnects=%d, want >= 1", h.Reconnects)
	}
	// RTT может быть 0 на очень быстром локале — не заваливаем, просто читаем
	_ = h.QuicRttMs
}
