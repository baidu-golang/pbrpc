BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOLDFLAGS="-X main.branch $(BRANCH) -X main.commit $(COMMIT)"

default: build

fmt:
	!(gofmt -l -s -d $(shell find . -name \*.go) | grep '[a-z]')

# go get honnef.co/go/tools/simple
gosimple:
	gosimple ./...

# go get honnef.co/go/tools/unused
unused:
	unused ./...


test:
	go test -timeout 20m -v -coverprofile cover.out -covermode atomic
	
race:
	go test -bench=. -benchtime=1s -run=^$

.PHONY: fmt test gosimple unused
