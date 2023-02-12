GOCMD=$(or $(shell which go), $(error "Missing dependency - no go in PATH"))
GOBUILD=$(GOCMD) build
CGCTL_BINARY_NAME=ntokinctl
VERSION=$(shell git branch --no-color --no-column --show-current)

# This cannot detect whether untracked files have yet to be added.
# That is sort-of a git feature, but can be a limitation here.
DIRTY=$(shell git diff-index --quiet HEAD -- || echo '-SNAPSHOT')
GIT_REF=$(shell git rev-parse --short HEAD --)
REVISION=$(GIT_REF)$(DIRTY)
LD_FLAGS := "-X github.com/michaelraskansky/ntokin/pkg/version.Version=$(VERSION)-$(REVISION)"

export REVISION

define compile_bin
	$(GOBUILD) -o $(CGCTL_BINARY_NAME) -ldflags $(LD_FLAGS)
endef

.PHONY: go_version
go_version: 
	$(GOCMD) version

build: 
	$(compile_bin)

.PHONY: serve
serve:
	./$(CGCTL_BINARY_NAME) serve --username $(USERNAME) \
	--password $(PASSWORD) \
	--subscriptions darwin.pushport-v16 \
	--kinesis-stream-arn $(STREAM_ARN)