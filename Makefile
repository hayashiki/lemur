VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE ?= on

.PHONY: test
test:
	go test -v ./...

.PHONY: deploy
deploy:
	@gcloud app deploy -q

.PHONY: deploy-tasks
deploy-tasks:
	@gcloud app deploy queue.yaml

.PHONY: deploy-scheduler
deploy-scheduler:
	@gcloud scheduler jobs create app-engine crowl-articles \
		--project lemurapp \
		--schedule "* 2 * * *" \
		--description "crow article" \
		--service "default" \
		--relative-url "/cron/docbase" \
		--http-method "GET" \
		--headers "Authorization=secrets" \
		--time-zone "Asia/Tokyo" \
		--attempt-deadline "60s"


.PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r .

$(GOBIN)/gobump:
	@cd && go get github.com/x-motemen/gobump/cmd/gobump

.PHONY: tag
tag:
	git tag -a "v$(VERSION)" -m "Release $(VERSION)"
	git push --tags

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: vet
vet:
	go vet ./...
