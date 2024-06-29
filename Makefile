.DEFAULT_GOAL := build

.PHONY:fmt
fmt:
	go fmt ./

.PHONY: vet
vet: fmt
	go vet ./

.PHONY: staticcheck
staticcheck: vet
	staticcheck ./

.PHONY: build
build: staticcheck
	go build -o /tmp/go_bin && /tmp/go_bin
