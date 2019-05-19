# istio-discovery

[![CircleCI](https://circleci.com/gh/Tufin/istio-discovery.svg?style=svg)](https://circleci.com/gh/Tufin/istio-discovery)

istio-discovery automatically detects egress traffic from Kubernetes cluster and assigns a matching ServiceEntry object for each host. It creates better visibility in the cluster traffic, allowing cluster operators impose better security boundaries in their security policy.

### Installation
```sh
$ make install
```
### Building
```sh
$ make build && make docker
```