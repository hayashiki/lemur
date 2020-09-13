VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE ?= on

.PHONY: show-version test deploy

test:
	go test -v ./...

deploy:
	gcloud app deploy -q

.PHONY:
show-version: $(GOBIN)/gobump
	@gobump show -r .

$(GOBIN)/gobump:
	@cd && go get github.com/x-motemen/gobump/cmd/gobump

.PHONY: tag
tag:
	git tag -a "v$(VERSION)" -m "Release $(VERSION)"
	git push --tags

