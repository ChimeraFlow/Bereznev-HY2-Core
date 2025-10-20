//go:build android || ios || mobile_skel

package mobile

import "fmt"

// Эти значения будут переопределяться через флаги линковки (ldflags) при сборке CI.
var (
	SdkVersion = "0.1.0-dev" // SemVer по умолчанию
	BuildTime  = "unknown"
	CommitHash = "unknown"
	EngineID   = "sing-tun" // можно менять на "skeleton"/"core" по контексту
)

// Version — возвращает человекочитаемую строку версии SDK.
func Version() string {
	return fmt.Sprintf("%s %s (%s)", sdkName, SdkVersion, EngineID)
}
