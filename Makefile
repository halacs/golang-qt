DIST := dist
APPNAME := go-qt

all: clean build-qt build

clean:
	rm $(DIST)/* || true

env:
	go install github.com/mappu/miqt/cmd/miqt-uic@latest
	mkdir -p $(DIST)

build: build-linux build-windows

build-linux: generate-qt-ui
	GOARCH=amd64 GOOS=linux go build -v -v -o "$(DIST)/$(APPNAME)" .

MIQT_DOCKER_DIR := $(shell go list -m -f '{{.Dir}}' github.com/mappu/miqt)
DOCKER_IMAGE := miqt/win64-qt6-static:latest
build-windows: generate-qt-ui
	docker build -t $(DOCKER_IMAGE) -f "$(MIQT_DOCKER_DIR)/docker/win64-cross-go1.26-qt6.5-static.Dockerfile" "$(MIQT_DOCKER_DIR)"
	docker run --rm \
		-u $(shell id -u):$(shell id -g) \
		-v "$(PWD):/src" \
		-v "/usr/lib/go-1.27:/usr/lib/go-1.27:ro" \
        -v "/usr/share/go-1.27:/usr/share/go-1.27:ro" \
		-v "$(shell go env GOPATH)/pkg/mod:/go/pkg/mod:ro" \
		-v "$(shell go env GOCACHE):/cache/go-build" \
		-e GOROOT=/usr/lib/go-1.27 \
		-e GOCACHE=/cache/go-build \
		-e GOPATH=/go \
		-w /src \
		$(DOCKER_IMAGE) \
		/usr/lib/go-1.27/bin/go build -v --tags=windowsqtstatic -ldflags="-s -w -H=windowsgui" -o "$(DIST)/$(APPNAME).exe" .

build-qt:
	cd layouts && miqt-uic -Qt6 -InFile first.ui -OutFile ui_first.go -Package layouts

generate-qt-ui:
	go generate ./...