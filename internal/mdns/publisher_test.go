package mdns_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/venkatkrishna07/mkdev/internal/mdns"
)

func TestPublisherEmptyClose(t *testing.T) {
	p := mdns.New(net.IPv4(192, 168, 1, 42))
	require.NoError(t, p.Close())
}
