.PHONY: docker

BINARY=istio-discovery
DOCKER_REPO=tufin
IMAGE=$(DOCKER_REPO)/istio-discovery

build:
	GOOS=linux GOARCH=amd64 go build -o .dist/$(BINARY)

clean:
	rm $(BINARY)

docker:
	docker build --build-arg=binary=$(BINARY) -t $(IMAGE) -f docker/Dockerfile .dist

deploy:
	docker push $(IMAGE)

