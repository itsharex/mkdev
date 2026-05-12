package mdns

import (
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/hashicorp/mdns"
	"github.com/venkatkrishna07/mkdev/internal/store"
)

// Publisher manages a set of mDNS service registrations — one per enabled
// route whose TLD is ".local". Other TLDs are silently skipped.
type Publisher struct {
	mu      sync.Mutex
	ip      net.IP
	servers map[string]*mdns.Server // keyed by route domain
}

// New constructs a Publisher bound to the given LAN IPv4.
func New(ip net.IP) *Publisher {
	return &Publisher{ip: ip, servers: map[string]*mdns.Server{}}
}

// Set diffs the desired route set against the currently published set and
// adjusts: registers new .local enabled routes, deregisters removed ones.
func (p *Publisher) Set(routes []store.Route) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	desired := map[string]store.Route{}
	for _, r := range routes {
		if !r.Enabled || !r.Shared || !strings.HasSuffix(r.Domain, ".local") {
			continue
		}
		desired[r.Domain] = r
	}
	for dom, srv := range p.servers {
		if _, keep := desired[dom]; !keep {
			_ = srv.Shutdown() // best-effort; map is replaced regardless
			delete(p.servers, dom)
		}
	}
	for dom, r := range desired {
		if _, exists := p.servers[dom]; exists {
			continue
		}
		srv, err := registerOne(dom, r.Target, p.ip)
		if err != nil {
			return err
		}
		p.servers[dom] = srv
	}
	return nil
}

// Close deregisters everything.
func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	var firstErr error
	for dom, srv := range p.servers {
		if err := srv.Shutdown(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(p.servers, dom)
	}
	return firstErr
}

func registerOne(domain, target string, ip net.IP) (*mdns.Server, error) {
	if ip == nil {
		return nil, errors.New("mdns: nil LAN ip")
	}
	host := strings.TrimSuffix(domain, ".local") + ".local."
	service, err := mdns.NewMDNSService(
		strings.TrimSuffix(domain, ".local"),
		"_https._tcp",
		"",
		host,
		443,
		[]net.IP{ip},
		[]string{"target=" + target, "managed=mkdev"},
	)
	if err != nil {
		return nil, err
	}
	srv, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}
	return srv, nil
}
