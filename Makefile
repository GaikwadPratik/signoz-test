SHELL:= /bin/bash

ARCH_TYPE := $(shell uname -m)
SED:= $(shell command -v gsed || command -v sed)
ARCH_VARIANT:= $(shell echo "$(ARCH_TYPE)" | $(SED) -e 's/x86_64/amd64/g' -e 's/i686/i386/g' -e 's/aarch64/arm64/g')

BUILD_DATA:= build-data
SVC_NAME:=signoz-test

GO_FLAGS := CGO_ENABLED=0 GOARCH=$(ARCH_VARIANT) GOOS=linux

.PHONY: build
build:
	mkdir -p $(BUILD_DATA)
	$(GO_FLAGS) go build -o $(BUILD_DATA)/$(SVC_NAME) 

.PHONY: docker-build
docker-build:
	docker build --force-rm --no-cache -t signoz-test:latest .

.PHONY: docker-clean
docker-clean:
	docker image rm -f signoz-test:latest

.PHONY: docker-compose-up
docker-compose-up:
	docker compose up --build --remove-orphans  --force-recreate --wait

.PHONY: docker-compose-down
docker-compose-down:
	docker compose down --remove-orphans