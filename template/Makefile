LDFLAGS += -X "main.buildTime=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.version=$(shell git rev-parse HEAD)"

.PHONY: init
init:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	@chmod +x ./hack/check.sh
	@chmod +x ./hooks/pre-commit

.PHONY: setup
setup: init
	git init
	@echo "Install pre-commit hook"
	@ln -s $(shell pwd)/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit || true
	dep init

.PHONY: check
check:
	@./hack/check.sh ${scope}

.PHONY: ci
ci: init
	@dep ensure
	@make check

.PHONY: build
build: check
	go build -ldflags '$(LDFLAGS)'

.PHONY: install
install: check
	go install -ldflags '$(LDFLAGS)'
