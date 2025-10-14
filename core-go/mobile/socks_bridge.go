//go:build (android || ios) && !mobile_skel

package mobile

/*
#cgo android,linux LDFLAGS: -llog
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	core "github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

var (
	t2sRunning atomic.Bool
	tunFile    *os.File
)

// StartTun2Socks — реальный мост TUN→SOCKS.
func StartTun2Socks(tunFd int, socksHost string, socksPort int) string {
	if t2sRunning.Load() {
		logI("tun2socks already running")
		return ""
	}

	addr := fmt.Sprintf("%s:%d", socksHost, socksPort)
	logI("starting tun2socks → " + addr)

	// Открываем TUN fd
	f := os.NewFile(uintptr(tunFd), "tun")
	if f == nil {
		return "invalid tun fd"
	}
	tunFile = f

	// TCP/UDP обработчики (новый API)
	tcpHandler := socks.NewTCPHandler(socksHost, uint16(socksPort))
	udpHandler := socks.NewUDPHandler(socksHost, uint16(socksPort), 60*time.Second)
	core.RegisterTCPConnHandler(tcpHandler)
	core.RegisterUDPConnHandler(udpHandler)

	// Регистрируем выход (из TUN наружу)
	core.RegisterOutputFn(func(data []byte) (int, error) {
		atomic.AddUint64(&bytesOut, uint64(len(data)))
		return f.Write(data)
	})

	t2sRunning.Store(true)
	logI("tun2socks started successfully")
	emit("tun2socks_started", fmt.Sprintf(`{"socks":"%s"}`, addr))

	// Основной цикл чтения из TUN
	safeGo(func() {
		buf := make([]byte, 65535)
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF || err.Error() == "file already closed" {
					break
				}
				logE("TUN read error: " + err.Error())
				break
			}
			if n > 0 {
				atomic.AddUint64(&bytesIn, uint64(n))
				core.InputToTun(buf[:n]) // ✅ заменено с InputPacket
			}
		}
		t2sRunning.Store(false)
		logI("tun2socks stopped (read loop ended)")
	})

	return ""
}

// StopTun2Socks останавливает мост.
func StopTun2Socks() {
	if !t2sRunning.Load() {
		return
	}
	t2sRunning.Store(false)
	if tunFile != nil {
		_ = tunFile.Close()
	}
	emit("tun2socks_stopped", "{}")
	logI("tun2socks stopped")
}
