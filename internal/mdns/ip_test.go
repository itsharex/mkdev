package mdns_test

import (
	"net"
	"testing"

	"github.com/venkatkrishna07/mkdev/internal/mdns"
)

func TestPrimaryLANIPv4ReturnsSomething(t *testing.T) {
	ip, err := mdns.PrimaryLANIPv4()
	if err != nil {
		if err.Error() == "" {
			t.Fatal("empty error")
		}
		return
	}
	if ip == nil || ip.IsLoopback() || ip.To4() == nil {
		t.Fatalf("got bogus ip: %v", ip)
	}
	_ = net.IPv4len
}
