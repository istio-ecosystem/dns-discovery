package dns

//
//import (
//	"testing"
//	"github.com/miekg/dns"
//	"strings"
//	"net"
//	"errors"
//	"github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
//	"github.com/stretchr/testify/require"
//	"time"
//	"fmt"
//)
//
//
//type fakeExchanger struct {
//	dnsMap map[string]string
//}
//
//func (f fakeExchanger) Exchange(m *dns.Msg) (*dns.Msg, error) {
//	domain := strings.TrimRight(m.Question[0].Name,".")
//	if address, ok := f.dnsMap[domain]; ok {
//		dnsRes := new(dns.Msg)
//		dnsRes.SetReply(m)
//		dnsRes.Answer = append(dnsRes.Answer, &dns.A{Hdr: dns.RR_Header{ Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60 },
//			A: net.ParseIP(address)})
//
//		return dnsRes, nil
//	}
//	return nil, errors.New("no such host")
//}
//
//type fakeServiceEntryCreator struct {
//	created []*v1alpha3.ServiceEntry
//}
//
//func (f *fakeServiceEntryCreator) Create(se *v1alpha3.ServiceEntry) {
//	f.created = append(f.created, se)
//}
//
//
//func TestServer(t *testing.T) {
//	proxy := NewServer(&fakeExchanger{dnsMap:map[string]string{"www.google.com": "1.2.3.4"}}, &fakeServiceEntryCreator{}, []string{".cluster.local"})
//	proxy.Start(":60653", "udp")
//	defer proxy.Stop()
//
//
//	fmt.Println("sleeping")
//	time.Sleep(200 * time.Second)
//	c, m := dns.Client{Net:"udp"}, &dns.Msg{}
//	m.SetQuestion("www.google.com.", dns.TypeA)
//
//	r, _, err := c.Exchange(m, ":6053")
//
//	if err != nil {
//		t.Fatal("failed to exchange", err)
//	}
//
//	require.Equal(t, "1.2.3.4", r.Answer[0].String())
//
//}
