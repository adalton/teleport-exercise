#
# Copyright 2021 Andy Dalton
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
SHELL = bash
MKDIR = mkdir -p
BUILDDIR = build
COVERAGEDIR=$(BUILDDIR)/coverage

GOTEST := go test
ifneq ($(shell which gotestsum),)
	GOTEST := gotestsum -- 
endif

all: proto $(BUILDDIR) $(BUILDDIR)/cgexec

$(BUILDDIR):
	$(MKDIR) $(BUILDDIR)

$(BUILDDIR)/cgexec: CGO_ENABLED=0
$(BUILDDIR)/cgexec: GOOS=linux
$(BUILDDIR)/cgexec: GOARCH=amd64
$(BUILDDIR)/cgexec: BUILDFLAGS=-buildmode pie -tags 'osusergo netgo static_build'
$(BUILDDIR)/cgexec: dep $(BUILDDIR) cmd/cgexec/cgexec.go
	go build -race -o $(BUILDDIR)/cgexec cmd/cgexec/cgexec.go

proto:
	@if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	protoc --proto_path=./service/v1/ --go_out=./service/v1 --go_opt=paths=source_relative --go-grpc_out=./service/v1 --go-grpc_opt=paths=source_relative ./service/v1/jobmanager.proto
.PHONY: proto

clean:
	$(RM) -r $(BUILDDIR) ./service/v1/jobmanager_grpc.pb.go  ./service/v1/jobmanager.pb.go
.PHONY: clean

$(COVERAGEDIR):
	$(MKDIR) -p $(COVERAGEDIR)

test: vet $(COVERAGEDIR)
	@$(GOTEST) -v -race -coverprofile=${COVERAGEDIR}/coverage.out -coverpkg=./... ./...
	@go tool cover -func=${COVERAGEDIR}/coverage.out -o ${COVERAGEDIR}/function-coverage.txt
	@go tool cover -html=${COVERAGEDIR}/coverage.out -o ${COVERAGEDIR}/coverage.html
.PHONY: test

# Not using $(GOTEST) here since root might not have it installed
inttest: vet $(BUILDDIR)/cgexec
	@cp $(BUILDDIR)/cgexec /tmp
	@sudo go test -v -race --tags=integration ./test/...
.PHONY: inttest

vet: dep
	@go vet -race ./...
.PHONY: vet

dep:
	@go mod download
	@go mod tidy
.PHONY: dep
