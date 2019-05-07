package main

import (
	"os"

	types "github.com/gogo/protobuf/types"
	v1alpha3 "github.com/michaelkipper/istio-client-go/pkg/apis/networking/v1alpha3"
	versionedclient "github.com/michaelkipper/istio-client-go/pkg/client/clientset/versioned"
	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	_ "github.com/golang/protobuf/jsonpb"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	log "github.com/sirupsen/logrus"
)

func createAndDeleteVirtualService(client *versionedclient.Clientset, namespace string) error {
	spec := v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-virtual-service",
		},
		Spec: v1alpha3.VirtualServiceSpec{
			VirtualService: istiov1alpha3.VirtualService{
				Hosts: []string{
					"*",
				},
				Gateways: []string{
					"test-gateway",
				},
				Http: []*istiov1alpha3.HTTPRoute{
					&istiov1alpha3.HTTPRoute{
						Match: []*istiov1alpha3.HTTPMatchRequest{
							&istiov1alpha3.HTTPMatchRequest{
								Uri: &istiov1alpha3.StringMatch{
									MatchType: &istiov1alpha3.StringMatch_Prefix{
										Prefix: "/",
									},
								},
							},
						},
						Timeout: &types.Duration{
							Seconds: int64(10),
						},
						Route: []*istiov1alpha3.HTTPRouteDestination{
							&istiov1alpha3.HTTPRouteDestination{
								Destination: &istiov1alpha3.Destination{
									Host: "test-service",
									Port: &istiov1alpha3.PortSelector{
										Port: &istiov1alpha3.PortSelector_Number{
											Number: uint32(8080),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	log.WithField("name", spec.GetName()).Info("Creating VirtualService")
	vs, err := client.NetworkingV1alpha3().VirtualServices(namespace).Create(&spec)
	if err != nil {
		log.WithField("error", err).Panic("Could not create virtual service")
	}

	err = client.NetworkingV1alpha3().VirtualServices(namespace).Delete(vs.GetName(), &metav1.DeleteOptions{})
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	namespace := os.Getenv("NAMESPACE")
	if len(kubeconfig) == 0 || len(namespace) == 0 {
		log.Fatalf("Environment variables KUBECONFIG and NAMESPACE need to be set")
	} else {
		log.WithFields(log.Fields{"kubeconfig": kubeconfig, "namespace": namespace}).Info("Building config")
	}
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}

	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}
	// Test VirtualServices
	vsList, err := ic.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get VirtualService in %s namespace: %s", namespace, err)
	}
	for i := range vsList.Items {
		vs := vsList.Items[i]
		log.Infof("Index: %d VirtualService Hosts: %+v\n", i, vs.Spec.VirtualService.GetHosts())
	}

	// Test DestinationRules
	drList, err := ic.NetworkingV1alpha3().DestinationRules(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get DestinationRule in %s namespace: %s", namespace, err)
	}
	for i := range drList.Items {
		dr := drList.Items[i]
		log.Printf("Index: %d DestinationRule Host: %+v\n", i, dr.Spec.GetHost())
	}

	// Test Policies
	pList, err := ic.AuthenticationV1alpha1().Policies(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Policy in %s namespace: %s", namespace, err)
	}
	for i := range pList.Items {
		p := pList.Items[i]
		log.Infof("Index: %d Policy Targets: %+v\n", i, p.Spec.GetTargets())
	}

	// Test MeshPolicies
	mpList, err := ic.AuthenticationV1alpha1().MeshPolicies().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list MeshPolicies: %+v", err)
	}
	for i := range mpList.Items {
		mp := mpList.Items[i]
		log.Infof("Index: %d MeshPolicy Name: %+v\n", i, mp.ObjectMeta.Name)
		_, err := ic.AuthenticationV1alpha1().MeshPolicies().Get(mp.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("Failed to get MeshPolicy named %s", mp.ObjectMeta.Name)
		}
	}

	// Test Gateways
	gList, err := ic.NetworkingV1alpha3().Gateways(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list Gateways: %+v", err)
	}
	for i, g := range gList.Items {
		log.Infof("Gateway %d: %s", i, g.ObjectMeta.GetName())
	}

	createAndDeleteVirtualService(ic, namespace)
}
