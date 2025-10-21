//go:build android || ios || mobile_skel

package version

import (
	"fmt"

	"github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/mobile"
)

// Эти значения будут переопределяться через флаги линковки (ldflags) при сборке CI.
var (
	SdkVersion = "0.1.0-dev" // SemVer по умолчанию
	BuildTime  = "unknown"
	CommitHash = "unknown"
	EngineID   = "sing-tun" // можно менять на "skeleton"/"core" по контексту
)

// Version — возвращает человекочитаемую строку версии SDK.
func Version() string {
	return fmt.Sprintf("%s %s (%s)", mobile.SdkName, SdkVersion, EngineID)
}
