package dns

import (
	"net"
	"strings"
	"testing"

	"github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

type fakeExchanger struct {
	dnsMap map[string]string
}

func (f fakeExchanger) Exchange(m *dns.Msg) (*dns.Msg, error) {
	domain := strings.TrimRight(m.Question[0].Name, ".")
	res := new(dns.Msg)
	res.Authoritative = true
	res.SetReply(m)
	if address, ok := f.dnsMap[domain]; ok {
		res.Answer = append(res.Answer, &dns.A{Hdr: dns.RR_Header{Name: domain + ".", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A: net.ParseIP(address)})

	}
	return res, nil
}

type fakeServiceEntryCreator struct {
	created []*v1alpha3.ServiceEntry
}

func (f *fakeServiceEntryCreator) Create(se *v1alpha3.ServiceEntry) {
	f.created = append(f.created, se)
}

func TestServer(t *testing.T) {
	sec := &fakeServiceEntryCreator{}
	proxy := NewServer(&fakeExchanger{dnsMap: map[string]string{"google.com": "1.2.3.4"}},
		sec,
		[]string{".cluster.local"})

	started := make(chan struct{})
	proxy.Start(":60653", "udp", func() {
		started <- struct{}{}
	})

	defer proxy.Stop()
	<-started
	c, m := dns.Client{Net: "udp"}, &dns.Msg{}
	m.SetQuestion("google.com.", dns.TypeA)

	r, _, err := c.Exchange(m, "127.0.0.1:60653")

	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "google.com.", r.Answer[0].Header().Name)
	require.Len(t, sec.created, 1)
}

func TestServerLocalAddress(t *testing.T) {
	sec := &fakeServiceEntryCreator{}
	proxy := NewServer(&fakeExchanger{dnsMap: map[string]string{"google.com": "1.2.3.4"}},
		sec,
		[]string{".cluster.local"})

	started := make(chan struct{})
	proxy.Start(":60654", "udp", func() {
		started <- struct{}{}
	})

	defer proxy.Stop()
	<-started
	c, m := dns.Client{Net: "udp"}, &dns.Msg{}
	m.SetQuestion("my.svc.cluster.local.", dns.TypeA)

	_, _, err := c.Exchange(m, "127.0.0.1:60654")

	if err != nil {
		t.Fatal(err)
	}

	require.Len(t, sec.created, 0)
}
