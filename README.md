# Teleport Exercise
This repository contains my implementation of the Teleport programming
challenge exercise.

## Tree Organization

The `cmd` package is contains the main executable programs.

The `pkg` package contains code that might be reused by some other component.
It contains:
* `adaptation` a collection of APIs that provide shims over standard library
  functions to enable unit testing of code that depends on those functions.
* `cgroup` is a collection of APIs to manage cgroups for jobs.  Currently
  this includes only v1.  The rationale for that is that on my development
  box, both cgroup v1 and v2 are available, but v1 is in use.
* `command` provides the implementation of commands from the `cmd` package.
* `config` contains the hard-coded configuration values.
* `io` contains a collection of components that implement i/o behavior
   (e.g., buffers, streams).
* `jobmanager` contains the JobManager components and components for
  managing individual jobs.

The `test` package includes a collection of test programs that enable us to
test functionality that isn't suitable for unit test (e.g., programs that
actually create and manage jobs, create and manage processes)

## Notes on running the tests

Currently the build depends on a go compiler in the user's path.

You can run the unit tests with `make test`.  You can run the integration
tests with `make inttest`.

The following integration tests are available:
* test/job/blkiolimit/blkiolimit\_test.go
  A test to illustrate that the blockio cgroup limit controls the job output.
  This could be extended to also check input limits.

* test/job/memorylimit/memorylimit\_test.go
  A test to illustrate that the memory cgroup limit controls the job output.

* test/job/cpulimit/cpulimit\_test.go
  A test to illustrate that the cpu cgroup limit controls the job output.

* test/job/pidnamespace/pidnamespace\_test.go
  A test to illustrate that the job is running in its own pid namespace

* test/job/networknamespace/networknamespace\_test.go
  A test to illustrate that the job is running in its own network namespace

* test/job/concurrentreads/concurrentreads\_test.go
  A test to illustrate that a single job can have multiple concurrent readers

You can build the `cgexec` binary using `make cgexec`.  The resulting binary
will be stored in `build/cgexec`.
