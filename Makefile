# TODO: add mac-os go installation scripts.

GO_VERSION = 1.20
IMAGE = snippet_bin
SHA = $(shell git rev-parse --short HEAD)

setup: install-go setup-go-env upgrade-go build containerize

install-go:
	@echo "[**** installing go binary ****]"
	# download go binary files for linux os.
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	# decompress and copy go binary files to [/usr/local] directory.
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	# delete the compressed go binary files.
	rm go$(GO_VERSION).linux-amd64.tar.gz

setup-go-env:
	@echo "[**** setting up go environment ****]"
	# update the os [PATH] to include the path to the go binary files.
	@echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.bashrc
	@echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.bashrc
	@echo "[**** successfully updated go binary path ****]"

upgrade-go:
	# delete the compressed go binary files.
	@echo "[**** deleting go binary directory for update ****]"
	rm -rf /usr/local/go
	# download go binary files for linux os.
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	# decompress and copy go binary files to [/usr/local] directory.
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	# delete the compressed go binary files.
	rm go$(GO_VERSION).linux-amd64.tar.gz

build:
	@echo "[**** building app ****]"
	go build -a -ldflags "-extldflags '-static' -w -s" -o app ./cmd/web


containerize:
	@echo "[**** containerizing app ****]"
	docker build --file ./Dockerfile --tag ${IMAGE}:${SHA} . 2>&1 | sed -e "s/^/ | /g"

test-with-coverage: test cover report

test:
	go test -v  -coverprofile=coverage.out ./...

cover:
	go tool cover -func coverage.out | grep "total:" | awk '{print ((int($$3) > 80) != 1)}'

report:
	go tool cover -html=coverage.out -o cover.html


