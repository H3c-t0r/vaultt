package parseip

import (
	"testing"
)

func Test_TrimLeadingZeroes(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"127.0.0.1", "127.0.0.1"},
		{"010.010.20.5", "10.10.20.5"},
		{"1.1.1.010", "1.1.1.10"},
		{"64:ff9b::192.00.002.33", "64:ff9b::192.0.2.33"},
		{"2001:db8:122:344:c0:2:2100::", "2001:db8:122:344:c0:2:2100::"},
		{"2001:db8:122:344::192.0.2.033", "2001:db8:122:344::192.0.2.33"},
	}
	for _, tt := range tests {
		if got := TrimLeadingZeroesIP(tt.in); got != tt.want {
			t.Errorf("TrimLeadingZeroesIP() = %v, want %v", got, tt.want)
		}
	}

	for _, tt := range tests {
		want := tt.want + "/32"
		if got := TrimLeadingZeroesCIDR(tt.in + "/32"); got != want {
			t.Errorf("TrimLeadingZeroesCIDR() = %v, want %v", got, want)
		}
	}
}
