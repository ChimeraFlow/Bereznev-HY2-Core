# Bereznev HY2 Core
> Freedom-grade transport layer built on QUIC & TLS.

**HTTP/3 (QUIC/TLS 1.3) transport core for mobile VPNs.**  
Drop-in **Hysteria2 client SDK** for Android (KMM-friendly).  
Implements modern censorship-resistant transport based on QUIC and TLS.

---

## 🧠 Why

Traditional protocols like **WireGuard** and **OpenVPN** are easily detected and throttled by DPI systems in restrictive networks (Russia, Iran, China).  
**Hysteria2 (HY2)** uses real **HTTP/3 (QUIC over UDP + TLS 1.3)** connections with valid certificates, making it appear to DPI as regular browser traffic.

Result:  
- ✅ Looks like normal HTTPS (SNI + ALPN = `h3`)  
- ✅ Works through most censorship and throttling systems  
- ✅ Handles packet loss better than TCP (built-in congestion control)  
- ✅ Uses real TLS certificates (Let's Encrypt, Cloudflare, etc.)

---

## 🚧 Status

**Early Preview** — minimal working API + AAR build via `gomobile bind`.

Planned features:
- Obfuscation modes (`faketls`, `obfs`)
- Log sink & debug levels
- `protect(fd)` callback for Android VPNService
- Auto reconnect / backoff
- iOS xcframework support
- Benchmark suite & adaptive MTU
- Per-app routing include/exclude (Android)

---

## ⚙️ Quick Start (Android AAR)

### 1. Install Go & gomobile
```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

### 2. Build AAR
``` bash
cd core
go mod tidy
cd ..
bash android/build-aar.sh
# → ./android/build/outputs/hy2lib.aar
```

3. Add to your Android/KMM project

Copy hy2lib.aar to your Android module and load via Gradle or /libs.

Then call from Kotlin:

```java 
val ok = Hy2.start(configJson)
Hy2.startTun2Socks(tunFd, "127.0.0.1", 1080)
println(Hy2.status())
Hy2.stop()
```

Default architecture:
HY2 starts a local SOCKS proxy (127.0.0.1:1080)
→ Android VpnService creates a TUN
→ go-tun2socks forwards all traffic through that proxy.



---


🧩 Minimal Client Config Example
```json 
{
  "log": { "level": "info" },
  "inbounds": [
    {
      "type": "socks",
      "listen": "127.0.0.1",
      "listen_port": 1080,
      "udp_enable": true,
      "sniff": true
    }
  ],
  "outbounds": [
    {
      "type": "hysteria2",
      "server": "your-domain.com:443",
      "password": "REDACTED",
      "tls": {
        "enabled": true,
        "server_name": "your-domain.com",
        "alpn": ["h3"]
      },
      "down_mbps": 100,
      "up_mbps": 20,
      "tag": "hy2"
    },
    { "type": "direct", "tag": "direct" },
    { "type": "block",  "tag": "block" }
  ],
  "route": { "final": "hy2" }
}
```
---


🧭 Architecture
```text
┌──────────────┐     UDP + TLS 1.3     ┌──────────────┐
│ Android TUN  │ ──► go-tun2socks ──►  │ HY2 Client   │
│ (VpnService) │ ◄── SOCKS 127.0.0.1   │ QUIC / HTTP3 │
└──────────────┘                       └──────────────┘
```

---



🛠 Roadmap
	•	🔧 Java/Kotlin log sink
	•	🧠 protect(fd) callback
	•	🔄 Auto reconnect
	•	📱 Per-app routing
	•	🍎 iOS xcframework
	•	⚙️ Benchmarks
	•	🕵️ Obfuscation modes



---


📜 License

Licensed under the Apache License 2.0.
See LICENSE for details.



---
🧭 Repository Info

Module: github.com/ChimeraFlow/Bereznev-HY2-Core/core
Go version: 1.22+

Example go.mod:
```go
module github.com/ChimeraFlow/Bereznev-HY2-Core/core

go 1.22

require (
    golang.org/x/mobile v0.0.0-20240830-abcdef123456 // indirect
)
```

---
⚡️ Bereznev HY2 Core — a modern transport engine for the next generation of VPNs.  
Built for performance, stealth, and the freedom to connect anywhere.


---

