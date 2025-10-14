//go:build android || ios || mobile_skel

package mobile

import "encoding/json"

// ErrCode — тип кодов ошибок SDK (стабильный API для KMM/Swift).
type ErrCode int

const (
	// УСПЕХ / обще-системные
	ErrOK ErrCode = iota
	ErrAlreadyRunning
	ErrInvalidConfig
	ErrEngineInitFailed
	ErrNotRunning
)

// String — человеко-читаемая строка для логов/UI.
func (e ErrCode) String() string {
	switch e {
	case ErrOK:
		return "ok"
	case ErrAlreadyRunning:
		return "already_running"
	case ErrInvalidConfig:
		return "invalid_config"
	case ErrEngineInitFailed:
		return "engine_init_failed"
	case ErrNotRunning:
		return "not_running"
	default:
		return "unknown_error"
	}
}

// MobileError — JSON-форма ошибки, стабильная для клиентов.
type MobileError struct {
	Code    ErrCode `json:"code"`
	Name    string  `json:"name"`              // дублируем String() для удобства клиентов
	Message string  `json:"message,omitempty"` // опциональное описание
}

// JSON возвращает сериализованную ошибку (code + name + message).
func (e ErrCode) JSON(message string) string {
	b, _ := json.Marshal(MobileError{
		Code:    e,
		Name:    e.String(),
		Message: message,
	})
	return string(b)
}
