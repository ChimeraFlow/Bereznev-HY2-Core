//go:build mobile_skel

package config

// В юнитах разрешаем пустой конфиг: подставляем безопасные дефолты.
// Это НЕ попадает в прод благодаря тегам.
func hy2TestFixup(c *HY2Config) {
	if c.Engine == "" {
		c.Engine = "sing"
	}
	if c.Server == "" {
		c.Server = "127.0.0.1:1" // валидный host:port, никуда не пойдём
	}
	if c.Password == "" {
		c.Password = "stub"
	}
	if len(c.ALPN) == 0 {
		c.ALPN = []string{"h3"}
	}
}
