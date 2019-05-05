package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
	_ "github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned/typed/networking/v1alpha3"
	log "github.com/sirupsen/logrus"

	"os"
	"os/signal"
	"syscall"

	"github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned"
	"github.com/miekg/dns"
	"github.com/tufin/istio-discovery/proxy"
	"gopkg.in/yaml.v2"
	istio "istio.io/api/networking/v1alpha3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type flags struct {
	autoApply bool
	address   string
	zones     string
	forward   string
}

var (
	dnsProxy    *proxy.DnsProxy
	istioClient *versioned.Clientset

	serviceEntryCh = make(chan *v1alpha3.ServiceEntry)
	serEntList     = new(v1alpha3.ServiceEntryList)
	errFlag        = func(arg string) error { return errors.New(fmt.Sprintf("%s must be set", arg)) }
)

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGSTOP, syscall.SIGTERM)

	flags := parseFlags()
	handler := dns.NewServeMux()

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		log.Fatalf("error creating kubernetes config, %s", err)
	}

	istioClient, err = versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create istio client: %s", err)
	}

	dnsProxy = proxy.New(flags.forward)

	for _, zone := range strings.Split(flags.zones, ",") {
		log.Infof("skipping zone %s", zone)
		handler.HandleFunc(zone, forward)
	}

	handler.HandleFunc(".", createAndForward)

	stopDrain := make(chan struct{})
	go drain(stopDrain, flags.autoApply)

	log.Infof("listening on %s", flags.address)
	go dns.ListenAndServe(flags.address, "udp", handler)

	<-stop
	stopDrain <- struct{}{}

	asJson, _ := yaml.Marshal(serEntList)

	log.Info(string(asJson))

}

func drain(stop chan struct{}, apply bool) {

	for {
		select {
		case se := <-serviceEntryCh:
			if apply {
				istioClient.NetworkingV1alpha3().ServiceEntries(v1.NamespaceDefault).Create(se)
			} else {
				serEntList.Items = append(serEntList.Items, *se)
			}
		case <-stop:
			break
		}
	}

}

func forward(w dns.ResponseWriter, r *dns.Msg) {
	if res, err := dnsProxy.Exchange(r); err != nil {
		log.Info(err) //TODO: check this
	} else {
		w.WriteMsg(res)
	}

}

func createAndForward(w dns.ResponseWriter, r *dns.Msg) {

	go func() {
		var hosts []string
		for _, q := range r.Question {
			hosts = append(hosts, strings.TrimRight(q.Name, "."))
		}
		serviceEntryCh <- &v1alpha3.ServiceEntry{
			ObjectMeta: metav1.ObjectMeta{Name: hosts[0]},
			TypeMeta:   metav1.TypeMeta{Kind: "ServiceEntry", APIVersion: "networking.istio.io/v1alpha3"},
			Spec: v1alpha3.ServiceEntrySpec{
				istio.ServiceEntry{
					Hosts:      hosts,
					Ports:      []*istio.Port{{Number: 443, Protocol: "HTTPS", Name: "https"}, {Number: 80, Protocol: "HTTP", Name: "http"}},
					Location:   istio.ServiceEntry_MESH_EXTERNAL,
					Resolution: istio.ServiceEntry_NONE}}}
	}()

	forward(w, r)
}

func parseFlags() *flags {

	flags := flags{}

	flag.BoolVar(&flags.autoApply, "autoApply", false, "Automatically autoApply generated ServiceEntry to kubernetes")
	flag.StringVar(&flags.address, "address", "localhost:53", "Listen address (ip:port)")
	flag.StringVar(&flags.zones, "zones", "", "Kubernetes authoritative zones")
	flag.StringVar(&flags.forward, "forward", "", "DNS server address (required)")

	flag.Parse()

	if flags.zones == "" || flags.forward == "" {
		flag.Usage()
		os.Exit(1)
	}

	return &flags
}
