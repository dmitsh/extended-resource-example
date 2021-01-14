DOCKER_REPO=docker.io/dmitsh
DOCKER_IMAGE_VER=0.1
DOCKER_IMAGE=extres:${DOCKER_IMAGE_VER}

build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' ./cmd/extres

docker-build:
	docker build -t ${DOCKER_IMAGE} .

docker-push:
	docker tag ${DOCKER_IMAGE} ${DOCKER_REPO}/${DOCKER_IMAGE} && docker push ${DOCKER_REPO}/${DOCKER_IMAGE}

.PHONY: build docker-build docker-push
