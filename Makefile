# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: easynode all test clean

GORUN = env GO111MODULE=on go run

easynode:
	$(GORUN) build/ci.go install ./cmd/easynode
	@echo "Done building."

all:
	$(GORUN) build/ci.go install
	@echo "Done building."

test:
	$(GORUN) build/ci.go test ./test

test_blockchain:
	$(GORUN) build/ci.go test ./blockchain/chain/ether
	$(GORUN) build/ci.go test ./blockchain/chain/polygonpos
	$(GORUN) build/ci.go test ./blockchain/chain/tron

test_collect:
	$(GORUN) build/ci.go test ./collect/service/cmd/chain/ether
	$(GORUN) build/ci.go test ./collect/service/cmd/chain/polygonpos
	$(GORUN) build/ci.go test ./collect/service/cmd/chain/tron2

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace
