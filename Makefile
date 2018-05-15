GOOS       ?= linux
GOARCH     ?= amd64
DOCKER_IMG ?= lloydg/ecr-mop

all: compile docker

compile:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/ecr-mop cmd/ecr-mop/main.go

docker: docker-build docker-push

docker-build:
	docker build -t $(DOCKER_IMG) .

docker-push:
	docker push $(DOCKER_IMG)

.PHONY: all compile docker docker-build docker-push
