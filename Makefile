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
	echo "$(DOCKER_PASS)" >/dev/null | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(IMAGE)

