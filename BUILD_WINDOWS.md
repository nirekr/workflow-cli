# Building the workflow-cli project on Windows
## Description

Dell EMC strongly recommends building the workflow-cli project on a Linux or Mac environment. However, it is possible to build on Windows. The instructions that follow presume that you are using the git-bash shell that is bundled with git on Windows.

## Before you begin

Make sure GOPATH is set (usually `c:\Users\<USER-NAME>\go`). 
This might have been done when you installed Go. To check using git-bash, run `echo $GOPATH`.
If it is not set, you have to set it. Using git-bash, this would be done by running `export GOPATH=/c/Users/<USER-NAME>/go`.

Make sure GOPATH\bin is added to your PATH variable. Using git-bash, this would be done by running `export PATH=$PATH:$GOPATH/bin`.

To install dependencies:
```
go get github.com/Masterminds/glide
go install github.com/Masterminds/glide
go get github.com/onsi/ginkgo/ginkgo
go install github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
glide install
GOOS=windows go get -d ./...
```
## Building

To build a Linux binary:
```
GOOS=linux go build -o ./bin/linux/workflow-cli
cp resources/endpoint_template.yaml ./bin/linux/endpoint.yaml
```

To build a Mac binary:

```
GOOS=darwin go build -o ./bin/darwin/workflow-cli
cp resources/endpoint_template.yaml ./bin/darwin/endpoint.yaml
```

To build a Windows binary:
```
GOOS=windows go build -o ./bin/windows/workflow-cli.exe
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
```

To run Integration Test:
```
GOOS=windows go build -o ./bin/windows/workflow-cli.exe
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bINTEGRATION\b"
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
```

To run Unit Tests:
```
GOOS=windows go build -o ./bin/windows/workflow-cli.exe
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30 --focus="\bUNIT\b"
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
```

To run all tests:
```
GOOS=windows go build -o ./bin/windows/workflow-cli.exe
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=30
cp resources/endpoint_template.yaml ./bin/windows/endpoint.yaml
```
