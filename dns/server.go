package dns

import (
	"strings"

	"github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
	_ "github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned/typed/networking/v1alpha3"
	log "github.com/sirupsen/logrus"

	"github.com/miekg/dns"
	"github.com/tufin/istio-discovery/istio"
	istio_v1alpha3 "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Server struct {
	server         *dns.Server
	mux            *dns.ServeMux
	serviceEntryCh chan *v1alpha3.ServiceEntry
	serEntList     *v1alpha3.ServiceEntryList
	exchanger      Exchanger
	creator        istio.ServiceEntryCreator
}

func NewServer(exchanger Exchanger, creator istio.ServiceEntryCreator, kubeZones []string) *Server {

	server := &Server{exchanger: exchanger,
		creator:        creator,
		serviceEntryCh: make(chan *v1alpha3.ServiceEntry),
		serEntList:     new(v1alpha3.ServiceEntryList), mux: dns.NewServeMux()}

	stopDrain := make(chan struct{})
	go server.drain(stopDrain)

	for _, zone := range kubeZones {
		server.mux.HandleFunc(zone, server.forward)
	}

	server.mux.HandleFunc(".", server.createAndForward)

	return server
}

func (s *Server) Start(address, network string, onStart func()) {
	s.server = &dns.Server{Addr: address, Net: network, Handler: s.mux}
	s.server.NotifyStartedFunc = onStart
	go s.server.ListenAndServe()
}

func (s *Server) Stop() {
	s.server.Shutdown()
}

func (s *Server) drain(stop chan struct{}) {

	for {
		select {
		case se := <-s.serviceEntryCh:
			s.creator.Create(se)
		case <-stop:
			break
		}
	}

}

func (s *Server) forward(w dns.ResponseWriter, r *dns.Msg) {
	if res, err := s.exchanger.Exchange(r); err != nil {
		log.Error(err)
	} else {
		w.WriteMsg(res)

	}

}

func (s *Server) createAndForward(w dns.ResponseWriter, r *dns.Msg) {

	go func() {
		var hosts []string
		for _, q := range r.Question {
			hosts = append(hosts, strings.TrimRight(q.Name, "."))
		}
		s.serviceEntryCh <- &v1alpha3.ServiceEntry{
			ObjectMeta: metav1.ObjectMeta{Name: hosts[0]},
			TypeMeta:   metav1.TypeMeta{Kind: "ServiceEntry", APIVersion: "networking.istio.io/v1alpha3"},
			Spec: v1alpha3.ServiceEntrySpec{
				istio_v1alpha3.ServiceEntry{
					Hosts:      hosts,
					Ports:      []*istio_v1alpha3.Port{{Number: 443, Protocol: "HTTPS", Name: "https"}, {Number: 80, Protocol: "HTTP", Name: "http"}},
					Location:   istio_v1alpha3.ServiceEntry_MESH_EXTERNAL,
					Resolution: istio_v1alpha3.ServiceEntry_NONE}}}
	}()

	s.forward(w, r)
}
