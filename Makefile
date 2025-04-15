# Include: minimum
-include .makefiles/minimum/Makefile
.makefiles/minimum/Makefile:
	@git clone https://github.com/tmknom/makefiles.git .makefiles >/dev/null 2>&1

# Variables: Go
REPO_ORIGIN ?= $(shell \git config --get remote.origin.url)
REPO_NAME = $(shell \basename -s .git $(REPO_ORIGIN))
REPO_OWNER = $(shell \gh config get -h github.com user)
VERSION = $(shell \git tag --sort=-v:refname | head -1)
COMMIT = $(shell \git rev-parse HEAD)
DATE = $(shell \date +"%Y-%m-%d")
URL = https://github.com/$(REPO_OWNER)/$(REPO_NAME)/releases/tag/$(VERSION)
LDFLAGS ?= "-X main.name=$(REPO_NAME) -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.url=$(URL)"

# Targets: Go
.PHONY: all
all: mod build test run ## all

.PHONY: mod
mod: ## manage modules
	go mod tidy
	go mod verify

.PHONY: deps
deps:
	go mod download

.PHONY: build
build: deps ## build executable binary
	go build -ldflags=$(LDFLAGS) -o bin/$(REPO_NAME) ./cmd/$(REPO_NAME)

.PHONY: install
install: deps ## install
	go install -ldflags=$(LDFLAGS) ./cmd/$(REPO_NAME)

.PHONY: run
run: build ## run command
	@printf "\n"
	@printf "testdata/valid-action.yml: \033[33m\n"
	@bin/actdocs generate --debug --format=json testdata/valid-action.yml | jq .
	@printf "\033[0m\n"
	@printf "testdata/valid-workflow.yml: \033[32m\n"
	@bin/actdocs generate --debug --sort testdata/valid-workflow.yml
	@printf "\033[0m\n"
	@printf "inject: \033[35m\n"
	@bin/actdocs inject --debug --dry-run --file=testdata/output.md testdata/valid-action.yml
	@printf "\033[0m"
	@git checkout testdata/output.md

.PHONY: test
test: lint ## test
	go test ./...

.PHONY: lint
lint: goimports vet ## lint go

.PHONY: vet
vet: ## static analysis by vet
	go vet ./...

.PHONY: goimports
goimports: ## update import lines
	goimports -w .

.PHONY: install-tools
install-tools: ## install tools for development
	go install golang.org/x/tools/cmd/goimports@latest

# Targets: GitHub Actions
.PHONY: lint-gha
lint-gha: lint/workflow lint/yaml ## Lint workflow files and YAML files

.PHONY: fmt-gha
fmt-gha: fmt/yaml ## Format YAML files

# Targets: Release
.PHONY: release
release: release/run ## Start release process
