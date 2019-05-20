# istio-discovery

[![CircleCI](https://circleci.com/gh/Tufin/istio-discovery.svg?style=svg)](https://circleci.com/gh/Tufin/istio-discovery)

## In a Nutshell 

istio-discovery automatically detects egress traffic in a Kubernetes cluster and creates a matching [Service Entry](https://istio.io/docs/reference/config/networking/v1alpha3/service-entry/) object for each host. This improves visibility to the cluster traffic and allows cluster operators to enforce a tight network security policy.

## The Problem 

When deploying Istio onto a Kubernetes cluster, the user needs to explicitly define each external end-point that the application relies on – this process which usually relies on trial and error makes it harder to deploy Istio.  

For example, if the application relies on storage.googleapis.com, the user must instruct Istio to allow access to this end-point by adding a [Service Entry](https://istio.io/docs/reference/config/networking/v1alpha3/service-entry/) definition. 

Up until Istio 1.0, the default behavior was to block access to external end-points - this created a connectivity issue and applications were breaking until the operator could discover all the end-points and configure them manually. 

Istio 1.1 changed the default to allow access to all external end-points - this resolves the connectivity problem and creates a security problem instead: now any malicious traffic can leave the cluster. 

## The Solution 

istio-discovery is a container that is deployed into the Kubernetes cluster as a proxy in front of the Kubernetes DNS service. 

The proxy sees all attempts to connect to external end-points by monitoring DNS lookups and automatically configures Istio to allow them by adding an Istio [Service Entry](https://istio.io/docs/reference/config/networking/v1alpha3/service-entry/) for each hostname. 

The proxy can be easily deployed to learn the egress communication patterns automatically. When the learning phase is done, the service entries can be deployed into a production environment to enforce a tight security policy. 

### Installation
#### Mac
```sh
$ make install
```
#### Linux
```sh
$ ./deploy.sh
```

### Building
```sh
$ make build && make docker
```
### Kubernetes Network Policies
If you are restricting access to your DNS service with a Kubernetes network policy, please note that this will change the DNS pod to listen on UDP port 54 instead of the standard port 53 and you should update the policy. (the DNS service continues to listen on the standard port UDP 53).
