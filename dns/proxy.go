package dns

import "github.com/miekg/dns"

type Exchanger interface {
	Exchange(m *dns.Msg) (*dns.Msg, error)
}

type DnsProxy struct {
	client  dns.Client
	address string
}

func NewProxy(address string) *DnsProxy {
	return &DnsProxy{address: address, client: dns.Client{Net: "udp"}}
}

func (p DnsProxy) Exchange(m *dns.Msg) (*dns.Msg, error) {
	r, _, err := p.client.Exchange(m, p.address)
	return r, err
}
