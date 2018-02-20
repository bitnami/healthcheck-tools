ROOT_PKG=github.com/bitnami-labs/healthcheck-tools
ROOT_PKG_DIR=${GOPATH}/src/$(ROOT_PKG)

# For creating the proper filename for later uploading with Travis
EXECUTABLE_FLAG=
VERSION=$(shell git describe --always --long --dirty)
ifdef GOOS
  ifdef GOARCH
	EXECUTABLE_FLAG:=-o $(TOOL)-v$(VERSION)-$(GOOS)-$(GOARCH)
	ifeq ($(GOOS),windows)
      EXECUTABLE_FLAG:=$(EXECUTABLE_FLAG).exe
    endif
  endif
endif


SELF_DIR:=$(dir $(lastword $(MAKEFILE_LIST)))

# since go1.8 people can use go without having to define a GOPATH env
# this is the default value the go tooling would assume.
GOPATH?=~/go

godep-save:
	cd $(ROOT_PKG_DIR) && godep save $$(scripts/gopkgs)

godep-restore:
	cd $(ROOT_PKG_DIR) && godep restore $$(scripts/gopkgs)
