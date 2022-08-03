.PHONY: all build clean download get-build-deps vet lint test

# Load relative to the common.mk file
include $(dir $(lastword $(MAKEFILE_LIST)))/vars.mk

include ./vars.mk

all:
	@$(MAKE) get-build-deps
	@$(MAKE) download
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) build
	@$(MAKE) test

define binary
$(1)-$(2)-$(3)$(4):
	@GOOS=$(2) GOARCH=$(3) go build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$$@ ./cmd/$(1)
endef

define binaries
$(call binary,ssl-checker,$1,$2,$3)
$(call binary,smtp-checker,$1,$2,$3)
endef

$(eval $(call binaries,linux,amd64,))
$(eval $(call binaries,linux,arm64,))
$(eval $(call binaries,darwin,amd64,))
$(eval $(call binaries,windows,amd64,.exe))

build: ssl-checker-linux-amd64 smtp-checker-linux-amd64

clean:
	@rm -rf $(BUILD_DIR)

download:
	$(GO_MOD) download

get-build-deps:
	@echo "+ Downloading build dependencies"
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

vet:
	@echo "+ Vet"
	@go vet ./...

lint:
	@echo "+ Linting package"
	@staticcheck ./...
	$(call fmtcheck, .)

test:
	@echo "+ Testing package"
	$(GO_TEST) ./...
