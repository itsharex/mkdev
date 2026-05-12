package proxy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeUpstream(t *testing.T) {
	cases := []struct {
		in     string
		host   string
		scheme string
		path   string
	}{
		{"localhost:3000", "localhost:3000", "http", ""},
		{"http://localhost:3000", "localhost:3000", "http", ""},
		{"https://gitlab.example.com", "gitlab.example.com", "https", ""},
		{"https://gitlab.example.com:8443", "gitlab.example.com:8443", "https", ""},
		{"https://gitlab.example.com/gitlab", "gitlab.example.com", "https", "/gitlab"},
		{"  localhost:3000  ", "localhost:3000", "http", ""},
	}
	for _, c := range cases {
		u, err := normalizeUpstream(c.in)
		require.NoError(t, err, c.in)
		require.Equal(t, c.host, u.Host, c.in)
		require.Equal(t, c.scheme, u.Scheme, c.in)
		require.Equal(t, c.path, u.Path, c.in)
	}
}

func TestNormalizeUpstreamRejectsEmpty(t *testing.T) {
	_, err := normalizeUpstream("")
	require.Error(t, err)

	_, err = normalizeUpstream("http://")
	require.Error(t, err)

	_, err = normalizeUpstream("https://")
	require.Error(t, err)
}
