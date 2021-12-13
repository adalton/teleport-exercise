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
EXECUTABLES  = $(BUILDDIR)/test-concurrentreads
EXECUTABLES += $(BUILDDIR)/test-blkiolimit
EXECUTABLES += $(BUILDDIR)/test-memorylimit
EXECUTABLES += $(BUILDDIR)/test-cpulimit
EXECUTABLES += $(BUILDDIR)/test-pidnamespace
EXECUTABLES += $(BUILDDIR)/test-networknamespace

GOTEST := go test
ifneq ($(shell which gotestsum),)
	GOTEST := gotestsum -- endif

all: $(BUILDDIR) $(BUILDDIR)/cgexec $(EXECUTABLES)

$(BUILDDIR):
	$(MKDIR) $(BUILDDIR)

$(BUILDDIR)/cgexec: CGO_ENABLED=0
$(BUILDDIR)/cgexec: GOOS=linux
$(BUILDDIR)/cgexec: GOARCH=amd64
$(BUILDDIR)/cgexec: BUILDFLAGS=-buildmode pie -tags 'osusergo netgo static_build'
$(BUILDDIR)/cgexec: dep $(BUILDDIR) cmd/cgexec/cgexec.go
	go build -race -o $(BUILDDIR)/cgexec cmd/cgexec/cgexec.go

$(BUILDDIR)/test-concurrentreads: dep $(BUILDDIR) test/job/concurrentreads/concurrentreads.go
	go build -race -o $(BUILDDIR)/test-concurrentreads test/job/concurrentreads/concurrentreads.go

$(BUILDDIR)/test-blkiolimit: dep $(BUILDDIR) test/job/blkiolimit/blkiolimit.go
	go build -race -o $(BUILDDIR)/test-blkiolimit test/job/blkiolimit/blkiolimit.go

$(BUILDDIR)/test-memorylimit: dep $(BUILDDIR) test/job/memorylimit/memorylimit.go
	go build -race -o $(BUILDDIR)/test-memorylimit test/job/memorylimit/memorylimit.go

$(BUILDDIR)/test-cpulimit: dep $(BUILDDIR) test/job/cpulimit/cpulimit.go
	go build -race -o $(BUILDDIR)/test-cpulimit test/job/cpulimit/cpulimit.go

$(BUILDDIR)/test-pidnamespace: dep $(BUILDDIR) test/job/pidnamespace/pidnamespace.go
	go build -race -o $(BUILDDIR)/test-pidnamespace test/job/pidnamespace/pidnamespace.go

$(BUILDDIR)/test-networknamespace: dep $(BUILDDIR) test/job/networknamespace/networknamespace.go
	go build -race -o $(BUILDDIR)/test-networknamespace test/job/networknamespace/networknamespace.go

clean:
	$(RM) -r $(BUILDDIR)
.PHONY: clean

$(COVERAGEDIR):
	$(MKDIR) -p $(COVERAGEDIR)

test: vet $(COVERAGEDIR)
	@$(GOTEST) -v -race -coverprofile=${COVERAGEDIR}/coverage.out -coverpkg=./... ./...
	@go tool cover -func=${COVERAGEDIR}/coverage.out -o ${COVERAGEDIR}/function-coverage.txt
	@go tool cover -html=${COVERAGEDIR}/coverage.out -o ${COVERAGEDIR}/coverage.html
.PHONY: test

vet: dep
	@go vet -race ./...
.PHONY: vet

dep:
	@go mod download
.PHONY: dep
