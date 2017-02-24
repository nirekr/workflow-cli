ORGANIZATION = cpsd
PROJECT = workflow-cli
BINARYNAME = workflow-cli
GOOUT = ./bin

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
	go get github.com/fatih/color
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
	env GOOS=linux go build -o $(GOOUT)/linux/$(BINARYNAME)

build-Mac:
	env GOOS=darwin go build -o $(GOOUT)/darwin/$(BINARYNAME)

build-Windows:
	env GOOS=windows go build -o $(GOOUT)/windows/$(BINARYNAME).exe
