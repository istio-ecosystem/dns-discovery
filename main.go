package main

import (
	"flag"
	"strings"

	_ "github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned/typed/networking/v1alpha3"
	log "github.com/sirupsen/logrus"

	"os"
	"os/signal"
	"syscall"

	"github.com/tufin/istio-discovery/dns"
	"github.com/tufin/istio-discovery/istio"
)

type flags struct {
	address string
	zones   string
	forward string
}

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGSTOP, syscall.SIGTERM)

	flags := parseFlags()

	var zones []string
	for _, zone := range strings.Split(flags.zones, ",") {
		log.Infof("skipping zone %s", zone)
		zones = append(zones, zone)
	}

	server := dns.NewServer(dns.NewProxy(flags.forward), istio.New(), zones)
	server.Start(flags.address, "udp")
	log.Infof("listening on %s", flags.address)
	defer server.Stop()

	<-stop

}

func parseFlags() *flags {

	flags := flags{}

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
