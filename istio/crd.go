package istio

import (
	_ "github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
	_ "github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned/typed/networking/v1alpha3"

	log "github.com/sirupsen/logrus"

	"github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
	"github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned"
	_ "github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type ServiceEntryCreator interface {
	Create(se *v1alpha3.ServiceEntry) error
}

type CrdCreator struct {
	istioClient *versioned.Clientset
}

func New() *CrdCreator {

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		log.Fatalf("error creating kubernetes config, %s", err)
	}

	istioClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create istio client: %s", err)
	}
	return &CrdCreator{istioClient: istioClient}
}

func (c CrdCreator) Create(se *v1alpha3.ServiceEntry) error {
	_, err := c.istioClient.NetworkingV1alpha3().ServiceEntries(v1.NamespaceDefault).Create(se)
	return err
}
