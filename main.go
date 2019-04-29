package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/istio/istio/pilot/pkg/config/kube/crd"
	"github.com/miekg/dns"
	"github.com/tufin/istio-discovery/kube"
	"github.com/tufin/istio-discovery/proxy"

	"os"
	"os/signal"
	"syscall"
	"time"
)

type flags struct {
	autoApply          bool
	kubeConfigFilePath string
	address            string
	zones              string
	forward            string
}

var (
	dnsProxy   *proxy.DnsProxy
	kubeClient *kube.Client

	serviceEntryCh = make(chan crd.ServiceEntry)
	errFlag        = func(arg string) error { return errors.New(fmt.Sprintf("%s must be set", arg)) }
)

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGSTOP, syscall.SIGTERM)

	flags := parseFlags()
	handler := dns.DefaultServeMux

	kc, err := kube.New()
	if err != nil {
		panic(err)
	}

	kubeClient = kc
	dnsProxy = proxy.New(flags.address)

	for _, zone := range strings.Split(flags.zones, " ") {
		handler.Handle(zone, dns.HandlerFunc(forward))
	}
	handler.Handle(".", dns.HandlerFunc(createAndForward))

	stopApply := make(chan bool)
	if flags.autoApply {
		go applyCrd(stopApply)
	}

	go dns.ListenAndServe(flags.address, "udp", handler)
	log.Infof("serving at %s", flags.address)

	<-stop
	log.Info("exiting")
	if stopApply != nil {
		stopApply <- true
	}

}

func applyCrd(stop chan bool) {

	for {
		select {
		case <-time.Tick(5 * time.Second):
			drainCRD()
		case <-stop:
			break
		}
	}

}

func drainCRD() {

	lenCRDCh := len(serviceEntryCh)

	for i := 0; i < lenCRDCh; i++ {
		//obj := <-serviceEntryCh
		//kubeClient.CreateCRDResource(obj)
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

	serviceEntryCh <- crd.ServiceEntry{}
	forward(w, r)
}

func parseFlags() *flags {

	flags := flags{}
	var zones string

	flag.BoolVar(&flags.autoApply, "autoApply", false, "Automatically autoApply generated ServiceEntry to kubernetes")
	flag.StringVar(&flags.kubeConfigFilePath, "kubeconfig.path", "", "Kubernetes cluster config file path (needed only when autoApply=true)")
	flag.StringVar(&flags.address, "address", "", "Listen address (ip:port)")
	flag.StringVar(&flags.zones, "zones", "", "Kubernetes authoritative zones")
	flag.StringVar(&flags.forward, "forward", "", "DNS server address")

	if zones == "" {
		panic(errFlag("zones"))
	}

	if flags.address == "" {
		panic(errFlag("address"))
	}

	flag.Parse()

	return &flags
}
