TOOL=permissions-checker

all: lint build

include ../../scripts/go.mk

build:
	mkdir -p $(ROOT_PKG_DIR)/cmd/$(TOOL)/bin && cd $(ROOT_PKG_DIR)/cmd/$(TOOL)/bin && go build $(EXECUTABLE_FLAG) -ldflags="-X main.VERSION=$(VERSION)" -i ../...

test:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go test ./...

install:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go install ./...

generate:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go generate ./...

clean:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go clean ../... && (test ! -d bin || rm bin/ -r)

lint:
	golint
	go vet
