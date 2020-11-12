DOCKER_IMAGE_VER=0.1

DOCKER_CONTAINER=extres:${DOCKER_IMAGE_VER}

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' ./cmd/extres

docker-build:
	docker build -t ${DOCKER_CONTAINER} .

docker-push:
	docker tag ${DOCKER_CONTAINER} docker.io/dmitsh/${DOCKER_CONTAINER} && docker push docker.io/dmitsh/${DOCKER_CONTAINER}

.PHONY: build docker-build docker-push