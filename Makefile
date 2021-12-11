SHELL = bash
MKDIR = mkdir

GOTEST := go test
ifneq ($(shell which gotestsum),)
	GOTEST := gotestsum --
endif

builddir:
	@mkdir -p build

cgexec: CGO_ENABLED=0
cgexec: GOOS=linux
cgexec: GOARCH=amd64
cgexec: BUILDFLAGS=-buildmode pie -tags 'osusergo netgo static_build'
cgexec: builddir cmd/cgexec/cgexec.go
	go build -o build/goexec cmd/cgexec/cgexec.go
.PHONY: cgexec

clean:
	$(RM) -r build
.PHONY: clean

test:
	$(GOTEST) -v -race ./...
.PHONY: test
