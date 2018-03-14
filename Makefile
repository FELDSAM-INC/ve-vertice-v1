#Copyright (c) 2014 Megam Systems.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
###############################################################################
# Makefile to compile cib.
# lists all the dependencies for test, prod and we can run a go build aftermath.
###############################################################################

GOPATH  := $(GOPATH):$(shell pwd)/../../../../
BUILD_DATE := $(shell date +%Y-%m-%d_%H:%M:%S%Z)
COMMIT_HASH := $(shell git rev-parse HEAD)

define HG_ERROR

FATAL: you need mercurial (hg) to download vertice dependencies.
       Check README.md for details


endef

define GIT_ERROR

FATAL: you need git to download vertice dependencies.
       Check README.md for details
endef

define BZR_ERROR

FATAL: you need bazaar (bzr) to download vertice dependencies.
       Check README.md for details
endef

.PHONY: all check-path get hg git bzr get-code test dep

all: check-path get test

build: check-path get _go_test _vertice

# It does not support GOPATH with multiple paths.
check-path:
ifndef GOPATH
	@echo "FATAL: you must declare GOPATH environment variable, for more"
	@echo "       details, please check README.md file and/or"
	@echo "       http://golang.org/cmd/go/#GOPATH_environment_variable"
	@exit 1
endif
	@exit 0

get: hg git bzr dep

hg:
	$(if $(shell hg), , $(error $(HG_ERROR)))

git:
	$(if $(shell git), , $(error $(GIT_ERROR)))

bzr:
	$(if $(shell bzr), , $(error $(BZR_ERROR)))

get-code:
	# rm -rf ~/.go
	go get $(GO_EXTRAFLAGS) -u -d -t -insecure ./...

godep:
	go get $(GO_EXTRAFLAGS) github.com/tools/godep
	godep restore ./...

dep:
	dep ensure -v -vendor-only

_go_test:
	go clean $(GO_EXTRAFLAGS) ./...
	go test $(GO_EXTRAFLAGS) ./...

_vertice:
	rm -f vertice
	go build $(GO_EXTRAFLAGS) \
		-ldflags="-X main.date=$(BUILD_DATE) -X main.commit=$(COMMIT_HASH)_lib_$(shell cd vendor/github.com/virtengine/libgo && git rev-parse HEAD)" \
		-o vertice ./cmd/vertice

_verticer:
	./vertice -v start --config ./conf/vertice.conf

test: _go_test _vertice # _verticer

_install_deadcode: git
	go get $(GO_EXTRAFLAGS) github.com/remyoudompheng/go-misc/deadcode

deadcode: _install_deadcode
	@go list ./... | sed -e 's;github.com/virtengine/vertice/;;' | xargs deadcode
