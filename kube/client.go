package kube

import (
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/rest"
)

type Client struct {
	restClient *rest.RESTClient
}

func New() (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}

	rc, err := rest.RESTClientFor(config)
	if err != nil {

		return nil, err
	}

	return &Client{restClient: rc}, nil
}

func (c Client) CreateCRDResource(namespace string, resource string, obj interface{}) error {

	return c.restClient.Post().
		Namespace(namespace).
		Resource(resource).
		Body(obj).
		Do().Error()
}
