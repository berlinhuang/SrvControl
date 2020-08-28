package test

import (
	errors "SrvControl/utils/error"
	"testing"
)

func TestErrCodeString(t *testing.T) {
	tests := []struct {
		err  errors.ErrCode
		want string
	}{
		{errors.ErrCodeProtocol, "PROTOCOL_ERROR"},
		//{0xd, "HTTP_1_1_REQUIRED"},
		{0xf, "unknown error code 0xf"},
	}
	for i, tt := range tests {
		got := tt.err.String()
		if got != tt.want {
			t.Errorf("%d. Error = %q; want %q", i, got, tt.want)
		}
	}
}
