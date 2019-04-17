package main

import (
	"flag"
)

type flags struct {
	apply                 bool
	clusterConfigFilePath string
}

func main() {

	_ := parseFlags()

}

func parseFlags() *flags {

	flags := flags{}

	flag.BoolVar(&flags.apply, "apply", false, "Automatically apply generated ServiceEntry to kubeernetes")
	flag.StringVar(&flags.clusterConfigFilePath, "cluster.config.path", "", "Kubernetes cluster config file path (needed only when apply=true)")

	flag.Parse()
	return &flags
}
