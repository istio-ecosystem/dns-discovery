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
	address  string
	zones    string
	forward  string
	loglevel string
}

func main() {

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGSTOP, syscall.SIGTERM)

	flags := parseFlags()

	var zones []string
	for _, zone := range strings.Split(flags.zones, ",") {
		log.Infof("skipping zone %s", zone)
		zones = append(zones, zone)
	}

	server := dns.NewServer(dns.NewProxy(flags.forward), istio.New(), zones)
	server.Start(flags.address, "udp", func() {
		log.Infof("listening on %s", flags.address)
	})

	defer server.Stop()
	<-sig
}

func parseFlags() *flags {

	flags := flags{}

	flag.StringVar(&flags.address, "address", "localhost:53", "Listen address (ip:port)")
	flag.StringVar(&flags.zones, "zones", "", "Kubernetes authoritative zones")
	flag.StringVar(&flags.forward, "forward", "", "DNS server address (required)")
	flag.StringVar(&flags.loglevel, "loglevel", "info", "set log level")

	flag.Parse()

	if flags.zones == "" || flags.forward == "" {
		flag.Usage()
		os.Exit(1)
	}

	setLogLevel(flags.loglevel)

	return &flags
}

func setLogLevel(levelname string) {
	level, err := log.ParseLevel(levelname)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
