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

Currently the build depends on a go compiler in the user's path.  I plan to
eventually move the build process into a docker container.

You can run the unit tests with `make test`.

You can build the `cgexec` binary using `make cgexec`.  The resulting binary
will be stored in `build/cgexec`  The tests current expect that binary to
be available in the `/tmp` directory; you'll have to copy it there manually.

I've provided the following tests to exercise functionality that is suitable
for unit testing.  I currently run this via `sudo`, but I plan to look into
setting up a cgroup configuration that would enable these to work as a
non-privileged user.

* test/job/blkiolimit/blkiolimit.go (build/test-blkiolimit)
  A test to illustrate that the blockio cgroup limit controls the job output.
  This could be extended to also check input limits.

* test/job/memorylimit/memorylimit.go (build/test-memorylimit)
  A test to illustrate that the memory cgroup limit controls the job output.

* test/job/cpulimit/cpulimit.go (build/test-cpulimit)
  A test to illustrate that the cpu cgroup limit controls the job output.

* test/job/pidnamespace/pidnamespace.go (build/test-pidnamespace)
  A test to illustrate that the job is running in its own pid namespace

* test/job/networknamespace/networknamespace.go (build/test-networknamespace)
  A test to illustrate that the job is running in its own network namespace

* test/job/concurrentreads/concurrentreads.go (build/test-concurrentreads)
  A test to illustrate that a single job can have multiple concurrent readers
