//go:build android || ios || mobile_skel

package mobile

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	socks5 "github.com/armon/go-socks5"
)

// Локальный SOCKS5-сервер: 127.0.0.1:PORT (по умолчанию 1080).
// Служит целью для варианта A (TUN -> SOCKS).
//
// В проде сюда добавим кастомный Dial с protectFD(fd) для Android VpnService.

var (
	socksMu      sync.Mutex
	socksLn      net.Listener
	socksSrv     *socks5.Server
	socksAddr    = "127.0.0.1:1080"
	socksRunning bool
)

type countingConn struct{ net.Conn }

func (c *countingConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	if n > 0 {
		bytesIn.Add(uint64(n))
	}
	return n, err
}
func (c *countingConn) Write(p []byte) (int, error) {
	n, err := c.Conn.Write(p)
	if n > 0 {
		bytesOut.Add(uint64(n))
	}
	return n, err
}

// StartLocalSocks запускает локальный SOCKS5-сервер на host:port.
// Пустая строка = ок, иначе текст ошибки.
func StartLocalSocks(host string, port int) string {
	socksMu.Lock()
	defer socksMu.Unlock()

	if socksRunning {
		logI("SOCKS already running at " + socksAddr)
		return ""
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port <= 0 {
		port = 1080
	}
	socksAddr = net.JoinHostPort(host, strconv.Itoa(port))

	// кастомный Dial: меряем "rtt", считаем reconnects, заворачиваем в countingConn
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		d := &net.Dialer{}
		c, err := d.DialContext(ctx, network, addr)
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			logW("SOCKS dial fail: " + err.Error())
			return nil, err
		}

		// 🔒 добавляем защиту сокета
		if tcpConn, ok := c.(*net.TCPConn); ok {
			raw, _ := tcpConn.File()
			if raw != nil {
				protectFD(int(raw.Fd()))
				_ = raw.Close() // закрываем временную дубликат-дескриптор
			}
		}

		quicRttMs.Store(elapsed)
		reconnects.Add(1)
		return &countingConn{Conn: c}, nil
	}

	conf := &socks5.Config{Dial: dial}
	srv, err := socks5.New(conf)
	if err != nil {
		logE("SOCKS init failed: " + err.Error())
		return "socks init failed: " + err.Error()
	}

	ln, err := net.Listen("tcp", socksAddr)
	if err != nil {
		logE("SOCKS listen failed: " + err.Error())
		return "socks listen failed: " + err.Error()
	}
	socksSrv = srv
	socksLn = ln
	socksRunning = true

	logI(fmt.Sprintf("SOCKS listening at %s", socksAddr))
	emit("socks_started", fmt.Sprintf(`{"addr":%q}`, socksAddr))

	safeGo(func() {
		if err := srv.Serve(ln); err != nil {
			logI("SOCKS serve stopped: " + err.Error())
		}
	})
	return ""
}

// StopLocalSocks останавливает локальный SOCKS5-сервер.
func StopLocalSocks() {
	socksMu.Lock()
	defer socksMu.Unlock()

	if !socksRunning {
		return
	}
	_ = socksLn.Close()
	socksLn = nil
	socksSrv = nil
	socksRunning = false

	emit("socks_stopped", "{}")
	logI("SOCKS stopped")
}

func LocalSocksAddr() string {
	socksMu.Lock()
	defer socksMu.Unlock()
	return socksAddr
}

func StartLocalSocks1080() string { return StartLocalSocks("127.0.0.1", 1080) }

// (Опционально) заготовка под кастомный Dial с protectFD:
// func protectDial(ctx context.Context, network, addr string) (net.Conn, error) {
// 	d := &net.Dialer{}
// 	c, err := d.DialContext(ctx, network, addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// TODO: попытаться извлечь fd и вызвать protectFD(fd) на Android.
// 	return c, nil
// }
