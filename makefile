CLANG ?= clang-14
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)
ARCH := $(shell uname --hardware-platform)

vmlinux:
	chmod a+x ./ebpf/headers/vmlinux.sh
	cd ./ebpf/headers && ./vmlinux.sh
headers: vmlinux
	chmod a+x ./ebpf/headers/update.sh
	cd ./ebpf/headers && ./update.sh

generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate: headers
	go generate ./...

mod:
	go mod tidy
build: mod generate
	@echo "Building..."
	rm -rfv ./bin
	mkdir -p ./bin
	GOARCH=amd64 GOOS=linux go build -o ./bin/ebpf-xdp main.go
	GOARCH=arm64 GOOS=linux go build -o ./bin/ebpf-xdp-arm64 main.go

docker-run: vmlinux
	docker run -ti --rm -v ```pwd```:/code -w /code quay.io/cilium/ebpf-builder:1648566014
