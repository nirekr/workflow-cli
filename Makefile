ORGANIZATION = cpsd
PROJECT = workflow-cli
BINARYNAME = workflow-cli
GOOUT = ./bin
SHELL = /bin/bash

# variable definitions
PKGPATH = github.com/dellemc-symphony/workflow-cli
COMMITHASH = $(shell git describe --tags --always --dirty)
BUILDDATE = $(shell date -u)
GOVERSION = $(shell go version)
ifndef BUILD_ID
	RELEASEVERSION = v0.0.1-dev
else
	RELEASEVERSION = v0.0.1-${BUILD_ID}
endif

#Flags to pass to main.go
LDFLAGS = -ldflags "-X '${PKGPATH}/cmd.binaryName=${BINARYNAME}' \
	  -X '${PKGPATH}/cmd.buildDate=${BUILDDATE}' \
	  -X '${PKGPATH}/cmd.commitHash=${COMMITHASH}' \
	  -X '${PKGPATH}/cmd.goVersion=${GOVERSION}' \
	  -X '${PKGPATH}/cmd.releaseVersion=${RELEASEVERSION}' "


default: deps build test

creds:
	@$(eval CREDS = $(subst :, ,$(GIT_CREDS)))
	@$(eval GIT_USER = $(word 1, $(CREDS)))  # this variable is set by Jenkinsfile
	@$(eval GIT_PASS = $(word 2, ,$(CREDS)))
	@git config --global user.name $(GIT_USER)
	@echo "machine github.com\n login $(GIT_USER)\n password $(GIT_PASS)" > ~/.netrc
	@chmod 600 ~/.netrc

deps:
	go get github.com/Masterminds/glide
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	glide install
	env GOOS=windows go get -d ./...

integration-test: build
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bINTEGRATION\b"

unit-test: build
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bUNIT\b"

mock: build
	go run mock/main/main.go

test: build
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30

cover-cmd: test
	go tool cover -html=cmd/cmd.coverprofile

build: build-Linux build-Mac build-Windows

build-Linux:
	env GOOS=linux go build -o $(GOOUT)/linux/$(BINARYNAME) $(LDFLAGS)

build-Mac:
	env GOOS=darwin go build -o $(GOOUT)/darwin/$(BINARYNAME) $(LDFLAGS)

build-Windows:
	env GOOS=windows go build -o $(GOOUT)/windows/$(BINARYNAME).exe $(LDFLAGS)