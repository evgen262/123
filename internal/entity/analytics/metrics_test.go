package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_XCFCUserAgentHeader_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		header XCFCUserAgentHeader
		want   bool
	}{
		{
			name:   "valid header web",
			header: "DeviceID=1369757d-523b-49fe-a5d4-46de64148078;DeviceType=web",
			want:   true,
		},
		{
			name:   "valid header android",
			header: "DeviceID=1369757d-523b-49fe-a5d4-46de64148078;DeviceType=Android",
			want:   true,
		},
		{
			name:   "valid header ios",
			header: "DeviceID=1369757d-523b-49fe-a5d4-46de64148078;DeviceType=iOS",
			want:   true,
		},
		{
			name:   "invalid header - wrong format missing '='",
			header: "DeviceID-invalid-uuidDeviceType-web",
			want:   false,
		},
		{
			name:   "invalid header - wrong format missing ';'",
			header: "DeviceID=invalid-uuidDeviceType=web",
			want:   false,
		},
		{
			name:   "invalid header - missing DeviceID",
			header: "DeviceType=web",
			want:   false,
		},
		{
			name:   "invalid header - missing DeviceType",
			header: "DeviceID=1369757d-523b-49fe-a5d4-46de64148078",
			want:   false,
		},
		{
			name:   "invalid header - wrong DeviceType",
			header: "DeviceID=1369757d-523b-49fe-a5d4-46de64148078;DeviceType=mobile",
			want:   false,
		},
		{
			name:   "invalid header - empty DeviceID",
			header: "DeviceID=;DeviceType=web",
			want:   false,
		},
		{
			name:   "invalid header - empty DeviceType",
			header: "DeviceID=;DeviceType=",
			want:   false,
		},
		{
			name:   "invalid header - empty header",
			header: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.header.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}
