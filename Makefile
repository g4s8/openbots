.PHONY: build bin clean test test-race docker vet site

.DEFAULT_GOAL := build

# Call `make V=1` in order to print commands verbosely.
ifeq ($(V),1)
    Q =
else
    Q = @
endif

DOCKER_IMAGE := g4s8/openbots
DOCKER_TAG := local

BUILD_ENV := GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on
BIN_DIR := $(shell pwd)/bin
GO_PKG := $(shell go list ./...)
GO_CMDS := $(shell go list ./cmd/...)
GO_BUILD_TAGS := 'osusergo netgo static_build' 
LDFLAGS := -extldflags -static
GO_BUILD_ARGS := -ldflags "$(LDFLAGS)" -tags "$(GO_BUILD_TAGS)"
GO_TEST_FLAGS := -tags $(GO_BUILD_TAGS)
GO_VET_FLAGS :=
SHELL_ARGS :=
ifeq ($(V),1)
	GO_BUILD_ARGS += -v
	SHELL_ARGS += -v
	GO_TEST_FLAGS += -v
	GO_VET_FLAGS += -v
endif

build:
	${Q}${BUILD_ENV} go build $(GO_BUILD_ARGS) -o /dev/null $(GO_PKG)

bin: $(BIN_DIR)
	${Q}${BUILD_ENV} go build $(GO_BUILD_ARGS) -o $(BIN_DIR) $(GO_CMDS)

$(BIN_DIR):
	${Q}mkdir -p $(BIN_DIR)

test: build
	${Q}${BUILD_ENV} go test $(GO_TEST_FLAGS) $(GO_PKG)

test-race: test
	${Q}${BUILD_ENV} CGO_ENABLED=1 go test $(GO_TEST_FLAGS) -race $(GO_PKG)

vet:
	${Q}go vet $(GO_VET_FLAGS) $(GO_PKG)

docker:
	${Q}docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

clean:
	${Q}rm -rf $(BIN_DIR)

site:
	${Q}hugo --source ./site/
