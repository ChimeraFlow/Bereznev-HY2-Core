//go:build mobile_skel

package mobile

import (
	"encoding/json"
	"testing"
)

func TestErrCode_String(t *testing.T) {
	cases := []struct {
		code ErrCode
		want string
	}{
		{ErrOK, "ok"},
		{ErrAlreadyRunning, "already_running"},
		{ErrInvalidConfig, "invalid_config"},
		{ErrEngineInitFailed, "engine_init_failed"},
		{ErrNotRunning, "not_running"},
		{ErrCode(999), "unknown_error"},
	}
	for _, c := range cases {
		if got := c.code.String(); got != c.want {
			t.Fatalf("String(%v)=%q, want %q", c.code, got, c.want)
		}
	}
}

func TestErrCode_JSON(t *testing.T) {
	raw := ErrInvalidConfig.JSON("bad json")
	var me MobileError
	if err := json.Unmarshal([]byte(raw), &me); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if me.Code != ErrInvalidConfig || me.Name != "invalid_config" || me.Message != "bad json" {
		t.Fatalf("unexpected payload: %+v", me)
	}
}
