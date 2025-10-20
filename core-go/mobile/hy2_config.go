//go:build android || ios || mobile_skel

package mobile

import (
	"encoding/json"
	"errors"
	"strings"
)

type HY2Config struct {
	Server       string   `json:"server"` // host:port
	Password     string   `json:"password"`
	SNI          string   `json:"sni,omitempty"`  // tls.server_name
	ALPN         []string `json:"alpn,omitempty"` // default ["h3"]
	UpMbps       int      `json:"up_mbps,omitempty"`
	DownMbps     int      `json:"down_mbps,omitempty"`
	IdleTimeoutS int      `json:"idle_timeout_s,omitempty"`
	Mode         string   `json:"mode,omitempty"` // "tun2socks" | "native-tun"
}

func (c *HY2Config) defaults() {
	if len(c.ALPN) == 0 {
		c.ALPN = []string{"h3"}
	}
	if c.Mode == "" {
		c.Mode = "tun2socks"
	}
}

func (c *HY2Config) validate() error {
	if c.Server == "" || !strings.Contains(c.Server, ":") {
		return errors.New("server must be host:port")
	}
	if c.Password == "" {
		return errors.New("password required")
	}
	return nil
}

// parseHY2Config читает cfgRaw (уже провалидированный расширенным JSON) и
// достает из него минимальный outbound:hysteria2 (Server/Password/SNI/ALPN/...).
// На первом шаге — простой unmarshal всей cfgRaw в HY2Config.
// Если у тебя сложная схема, позже можно вытащить поля из "outbounds".
func parseHY2Config() (HY2Config, error) {
	var hc HY2Config
	if len(cfgRaw) == 0 {
		return hc, errors.New("empty config")
	}
	// попытка простого JSON → struct; при необходимости заменишь на sjson.Node навигацию.
	if err := jsonUnmarshal(cfgRaw, &hc); err != nil {
		return hc, err
	}
	hc.defaults()
	return hc, hc.validate()
}

func jsonUnmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
