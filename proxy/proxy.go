package proxy

import "github.com/miekg/dns"

type DnsProxy struct {
	proxy   dns.Client
	address string
}

func New(address string) *DnsProxy {
	return &DnsProxy{address: address, proxy: dns.Client{Net: "udp"}}
}

func (p DnsProxy) Exchange(m *dns.Msg) (*dns.Msg, error) {
	r, _, err := p.proxy.Exchange(m, p.address)
	return r, err
}
