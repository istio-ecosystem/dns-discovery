package main

import (
	"flag"
)

type flags struct {
	apply              bool
	kubeConfigFilePath string
	zones              string
	forward            string
}

func main() {

	parseFlags()

}

func parseFlags() *flags {

	flags := flags{}

	flag.BoolVar(&flags.apply, "apply", false, "Automatically apply generated ServiceEntry to kubernetes")
	flag.StringVar(&flags.kubeConfigFilePath, "kubeconfig.path", "", "Kubernetes cluster config file path (needed only when apply=true)")
	flag.StringVar(&flags.zones, "zones", "", "Kubernetes authoritative zones")
	flag.StringVar(&flags.forward, "forward", "", "DNS server address")

	flag.Parse()
	return &flags
}
