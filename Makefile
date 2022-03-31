# This option causes make to display a warning whenever an undefined variable is expanded.
MAKEFLAGS += --warn-undefined-variables

# Disable any builtin pattern rules, then speedup a bit.
MAKEFLAGS += --no-builtin-rules

# If this variable is not set, the program /bin/sh is used as the shell.
SHELL := /bin/bash

# The arguments passed to the shell are taken from the variable .SHELLFLAGS.
#
# The -e flag causes bash with qualifications to exit immediately if a command it executes fails.
# The -u flag causes bash to exit with an error message if a variable is accessed without being defined.
# The -o pipefail option causes bash to exit if any of the commands in a pipeline fail.
# The -c flag is in the default value of .SHELLFLAGS and we must preserve it.
# Because it is how make passes the script to be executed to bash.
.SHELLFLAGS := -eu -o pipefail -c

# Disable any builtin suffix rules, then speedup a bit.
.SUFFIXES:

# Sets the default goal to be used if no targets were specified on the command line.
.DEFAULT_GOAL := help

#
# Variables for the file and directory path
#
ROOT_DIR ?= $(shell $(GIT) rev-parse --show-toplevel)
MARKDOWN_FILES ?= $(shell find . -name '*.md')
YAML_FILES ?= $(shell find . -name '*.y*ml')
SHELL_FILES ?= $(shell find . -name '*.sh')

#
# Variables to be used by Git and GitHub CLI
#
GIT ?= $(shell \command -v git 2>/dev/null)
GH ?= $(shell \command -v gh 2>/dev/null)
GIT_EXCLUSIVES ?= ':!*.md' ':!Makefile' ':!VERSION' ':!.*' ':!.github/*'

#
# Variables to be used by Docker
#
DOCKER ?= $(shell \command -v docker 2>/dev/null)
DOCKER_PULL ?= $(DOCKER) pull
DOCKER_WORK_DIR ?= /work
DOCKER_RUN_OPTIONS ?=
DOCKER_RUN_OPTIONS += -it
DOCKER_RUN_OPTIONS += --rm
DOCKER_RUN_OPTIONS += -v $(ROOT_DIR):$(DOCKER_WORK_DIR)
DOCKER_RUN_OPTIONS += -w $(DOCKER_WORK_DIR)
DOCKER_RUN_SECURE_OPTIONS ?=
DOCKER_RUN_SECURE_OPTIONS += --user 1111:1111
DOCKER_RUN_SECURE_OPTIONS += --read-only
DOCKER_RUN_SECURE_OPTIONS += --security-opt no-new-privileges
DOCKER_RUN_SECURE_OPTIONS += --cap-drop all
DOCKER_RUN_SECURE_OPTIONS += --network none
DOCKER_RUN ?= $(DOCKER) run $(DOCKER_RUN_OPTIONS)
SECURE_DOCKER_RUN ?= $(DOCKER_RUN) $(DOCKER_RUN_SECURE_OPTIONS)

#
# Variables for the image name
#
REGISTRY ?= ghcr.io/tmknom/dockerfiles
PRETTIER ?= $(REGISTRY)/prettier:latest
MARKDOWNLINT ?= $(REGISTRY)/markdownlint:latest
YAMLLINT ?= $(REGISTRY)/yamllint:latest
ACTIONLINT ?= rhysd/actionlint:latest
SHELLCHECK ?= koalaman/shellcheck:stable
SHFMT ?= mvdan/shfmt:latest

#
# Variables for the version
#
VERSION ?= $(shell \cat VERSION)
SEMVER ?= "v$(VERSION)"
MAJOR_VERSION ?= $(shell version=$(SEMVER) && echo "$${version%%.*}")

#
# Development
#
.PHONY: all
all: build lint test run ## all

.PHONY: mod
mod: ## manage modules
	go mod tidy
	go mod verify

.PHONY: build
build: mod ## build executable binary
	go build -o bin/actdocs ./cmd/actdocs

.PHONY: run
run: build ## run command
	@printf "\n"
	@printf "fixtures/valid-workflow.yml: \033[32m\n"
	@bin/actdocs workflow fixtures/valid-workflow.yml
	@printf "\033[0m\n"
	@printf "fixtures/valid-action.yml: \033[35m\n"
	@bin/actdocs action fixtures/valid-action.yml
	@printf "\033[0m"

.PHONY: test
test: ## test all
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

#
# Lint
#
.PHONY: lint-all
lint-all: lint lint-markdown lint-yaml lint-action lint-shell ## lint all

.PHONY: lint-markdown
lint-markdown: ## lint markdown by markdownlint and prettier
	$(SECURE_DOCKER_RUN) $(MARKDOWNLINT) --dot --config .markdownlint.yml $(MARKDOWN_FILES)
	$(SECURE_DOCKER_RUN) $(PRETTIER) --check --parser=markdown $(MARKDOWN_FILES)

.PHONY: lint-yaml
lint-yaml: ## lint yaml by yamllint and prettier
	$(SECURE_DOCKER_RUN) $(YAMLLINT) --strict --config-file .yamllint.yml .
	$(SECURE_DOCKER_RUN) $(PRETTIER) --check --parser=yaml $(YAML_FILES)

.PHONY: lint-action
lint-action: ## lint action by actionlint
	$(SECURE_DOCKER_RUN) $(ACTIONLINT) -color -ignore '"permissions" section should not be empty.'

.PHONY: lint-shell
lint-shell: ## lint shell by shellcheck and shfmt
ifneq ($(SHELL_FILES),)
	$(SECURE_DOCKER_RUN) $(SHELLCHECK) $(SHELL_FILES)
endif
	$(SECURE_DOCKER_RUN) $(SHFMT) -i 2 -ci -bn -d .

#
# Format code
#
.PHONY: format
format: format-markdown format-yaml format-shell ## format all

.PHONY: format-markdown
format-markdown: ## format markdown by prettier
	$(SECURE_DOCKER_RUN) $(PRETTIER) --write --parser=markdown $(MARKDOWN_FILES)

.PHONY: format-yaml
format-yaml: ## format yaml by prettier
	$(SECURE_DOCKER_RUN) $(PRETTIER) --write --parser=yaml $(YAML_FILES)

.PHONY: format-shell
format-shell: ## format shell by shfmt
	$(SECURE_DOCKER_RUN) $(SHFMT) -i 2 -ci -bn -w .

#
# Release management
#
release: ## release
	$(GIT) tag --force --message "$(SEMVER)" "$(SEMVER)" && \
	$(GIT) push --force origin "$(SEMVER)"

bump: input-version commit create-pr ## bump version

input-version:
	@echo "Current version: $(VERSION)" && \
	read -rp "Input next version: " version && \
	echo "$${version}" > VERSION

commit:
	$(GIT) switch -c "bump-$(SEMVER)" && \
	$(GIT) add VERSION && \
	$(GIT) commit -m "Bump up to $(SEMVER)"

create-pr:
	$(GIT) push origin $$($(GIT) rev-parse --abbrev-ref HEAD) && \
	$(GH) pr create --title "Bump up to $(SEMVER)" --body "" --web

#
# Git shortcut
#
.PHONY: diff
diff: ## git diff only features
	@$(GIT) diff $(SEMVER)... -- $(GIT_EXCLUSIVES)

.PHONY: log
log: ## git log only features
	@$(GIT) log $(SEMVER)... -- $(GIT_EXCLUSIVES)

#
# General
#
.PHONY: install
install: ## install docker images
	$(DOCKER_PULL) $(PRETTIER)
	$(DOCKER_PULL) $(MARKDOWNLINT)
	$(DOCKER_PULL) $(YAMLLINT)
	$(DOCKER_PULL) $(ACTIONLINT)
	$(DOCKER_PULL) $(SHELLCHECK)
	$(DOCKER_PULL) $(SHFMT)

.PHONY: clean
clean: ## clean all
	echo "fix me"

.PHONY: help
help: ## show help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
