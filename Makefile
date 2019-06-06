.PHONY: docker

BINARY=istio-discovery
DOCKER_REPO=tufin
IMAGE=$(DOCKER_REPO)/istio-discovery

clean:
	rm $(BINARY)

build:
	GOOS=linux GOARCH=amd64 go build -o .dist/$(BINARY)

test:
	go test `go list ./...`

docker:
	docker build --build-arg=binary=$(BINARY) -t $(IMAGE) -f docker/Dockerfile .dist

deploy:
	@echo "$(DOCKER_PASS)" | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(IMAGE)

install:
	$(eval DNS_DEPLOYMENT=$(shell kubectl get deploy -n kube-system -l k8s-app=kube-dns -o=custom-columns=NAME:.metadata.name | tail -n1 2>/dev/null))
	@if [ -z $(DNS_DEPLOYMENT)  ]; then\
		echo "could not detect DNS deployment for K8s cluster. Can not install istio-discovery";\
		exit 1;\
	fi
	@echo patching "$(DNS_DEPLOYMENT)"
	@kubectl patch deploy -n kube-system $(DNS_DEPLOYMENT) -p "`<kubernetes/deploy_patch.yaml`"
	@kubectl patch svc -n kube-system kube-dns -p "`<kubernetes/service_patch.yaml`"
	@kubectl patch clusterrole system:$(DNS_DEPLOYMENT) -p "`sed "s/#DEPLOYMENT#/$(DNS_DEPLOYMENT)/g" kubernetes/clusterrole.yaml`"


