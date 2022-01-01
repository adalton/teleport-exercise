# Teleport Interview Exercise

This repository contains my implementation of the Teleport programming
challenge exercise.

## Requirements
This subsection is copied from from
https://github.com/gravitational/careers/blob/main/challenges/systems/challenge.md

### Summary

Implement a prototype job worker service that provides an API to run arbitrary
Linux processes.

### Rationale

This exercise has two goals:

* It helps us to understand what to expect from you as a developer, how you
  write production code, how you reason about API design and how you
  communicate when trying to understand a problem before you solve it.
* It helps you get a feel for what it would be like to work at Teleport, as
  this exercise aims to simulate our day-as-usual and expose you to the type
  of work we're doing here.

We believe this technique is not only better, but also is more fun compared
to whiteboard/quiz interviews so common in the industry.  It's not without the
downsides - it could take longer than traditional interviews.

[Some of the best teams use coding challenges.](https://sockpuppet.org/blog/2015/03/06/the-hiring-post/)

We appreciate your time and are looking forward to hack on this project together.

### Requirements

There are 6 engineering levels at Teleport. It's possible to score on level 1-5
through coding challenge.

Level 6 is only for internal promotions. Check the
[engineering levels document](https://raw.githubusercontent.com/gravitational/careers/main/levels.pdf)
for more details.

Start with a very brief Google doc that covers the edge cases and design
approach and post it to the Slack channel. After the doc is approved, implement
interfaces and an example program using the library.

Add a couple of high quality tests that cover happy and unhappy scenarios.

Split the submission into 2-3 pull requests for us to review. We will review
every pull request and provide our feedback.

We are going to compile the program, test it and get back to you.

#### Level 5

##### Library

* Worker library with methods to start/stop/query status and get the output of a job.
* Library should be able to stream the output of a running job.
  * Output should be from start of process execution.
  * Multiple concurrent clients should be supported.
* Add resource control for CPU, Memory and Disk IO per job using cgroups.
* Add resource isolation for using PID, mount, and networking namespaces.

##### API

* [GRPC](https://grpc.io) API to start/stop/get status/stream output of a running process.
* Use mTLS authentication and verify client certificate. Set up strong set of
  cipher suites for TLS and good crypto setup for certificates. Do not use any
  other authentication protocols on top of mTLS.
* Use a simple authorization scheme.

##### Client

* CLI should be able to connect to worker service and start, stop, get status,
  and stream output of a job.

### Guidance

#### Interview process

The interview team will join the Slack channel. The team consists of the
engineers who will be working with you. Ask them about the engineering culture,
work and life balance, or anything else that you would like to learn about
Teleport.

Before writing the actual code, create a brief design document in Google Docs or
markdown in GitHub and share with the team.

This document should cover at least: design approach, trade-offs, scope,
proposed API, and security considerations. For security, please include
considerations for transport layer, authentication, and authorization.

Please avoid writing an overly detailed design document. Use this document to
make sure the team could provide feedback on your design and demonstrate that
you've investigated the problem space.

Split your code submission using pull requests and give the team an opportunity
to review the PRs. A good "rule of thumb" to follow is that the final PR
submission is a formality adding a small feature set - it means that the team
had an opportunity to contribute the feedback during multiple well defined
stages of your work.

Our team will do their best to provide a high quality review of the submitted
pull requests in a reasonable time frame. You are spending your time on this,
we are going to contribute our time too.

After the final submission, the interview team will assemble and vote using a
"+1, -2" anonymous voting system: +1 is submitted whenever a team member accepts
the submission, -2 otherwise.

In case of a positive result, we will connect you to our HR team who will
collect one-two references and will work out other details. You can start the
reference collection process in parallel if you would like to speed up the
process.

After reference collection, our ops team will send you an offer.

In case of a negative score result, hiring manager will contact you and share a
list of the key observations from the team that affected the result.

## Code and project ownership

This is a test challenge and we have no intent of using the code you've
submitted in production. This is your work, and you are free to do whatever
you feel is reasonable with it. In the scenario when you don't pass, you can
open source it with any license and use it as a portfolio project.

#### Areas of focus

Teleport focuses on networking, infrastructure and security.

These are the areas we will be evaluating in the submission:

* Use consistent coding style. We follow
  [Go Coding Style](https://github.com/golang/go/wiki/CodeReviewComments)
  for the Go language. If you are going to use a different language, please pick
  coding style guidelines and let us know what they are.
* At the minimum, create tests for authentication, networking, and unhappy scenario.
* Make sure builds are reproducible. Pick any vendoring/packaging system that
  will allow us to get consistent build results.
* Ensure error handling and error reporting is consistent. The system should
  report clear errors and not crash under non-critical conditions.
* Avoid concurrency and networking errors. Most of the issues we've seen in
  production are related to data races, networking error handling or goroutine
  leaks. We will be looking for those errors in your code.
* Security. Use strong authentication and simplest, but robust authorization.
  Set up the strongest transport encryption you can. Test it.

#### Trade-offs

Write as little code as possible, otherwise this task will consume too much
time and code quality will suffer.

Please cut corners, for example configuration tends to take a lot of time, and
is not important for this task.

Use hardcoded values as much as possible and simply add TODO items showing your
thinking, for example:

```
  // TODO: Add configuration system.
  // Consider using CLI library to support both
  // environment variables and reasonable default values,
  // for example https://github.com/alecthomas/kingpin
```

Comments like this one are really helpful to us. They save yourself a lot of
time and demonstrate that you've spent time thinking about this problem and
provide a clear path to a solution.

Consider making other reasonable trade-offs. Make sure you communicate them to
the interview team.

Here are some other trade-offs that will help you to spend less time on the task:

Do not implement a system that scales or is highly performing. Describe which
performance improvements you would add in the future.
High availability. It is OK if the system is not highly available. Write down
how you would make the system highly available and why your system is not.
Do not try to achieve full test coverage. This will take too long. Take two
key components, e.g. authentication/authorization layer and networking and
implement one or two test cases that demonstrate your approach to testing.

#### Pitfalls and Gotchas

To help you out, we've composed a list of things that previously resulted in a
no-pass from the interview team:

* Scope creep. Candidates have tried to implement too much and ran out of time
  and energy. To avoid this pitfall, use the simplest solution that will work.
  Avoid writing too much code. We've seen candidates' code introducing caching
  and making many mistakes in the caching layer validation logic. Not having
  caching would have solved this problem.
* Data races. We will scan the code with a race detector and do our best to
  find data races in the code. Avoid global state as much as possible; if using
  global state, write down a good description why it is necessary and protect
  it against data races.
* Deadlocks. When using mutexes, channels or any other synchronization
  primitives, make sure the system won't deadlock. We've seen candidates' code
  holding a mutex and making a network call without timeouts in place. Be extra
  careful with networking and sync primitives.
* Unstructured code. We've seen candidates leaving commented chunks of code,
  having one large file with all the code, not having code structure at all.
* Not communicating. Some candidates have submitted all their code to the master
  branch without raising pull requests, which does not give us the ability to
  provide feedback on the various implementation phases. We are a distributed
  team, so structured, asynchronous communication is critical to us.
* Implementing custom security algorithms/authentication schemes is always a
  bad idea unless you are a trained security researcher/engineer. It is
  definitely a bad idea for this task - try to stick to industry proven security
  methods as much as possible.

#### Questions

It is OK to ask the interview team questions. Some folks stay away from asking
questions to avoid appearing less experienced, so we provide examples of
questions to ask and questions we expect candidates to figure out on their own.

Here is a great question to ask:

> Is it OK to pre-generate secret data and put the secrets in the repository
> for a proof of concept? I will add a note that we will auto-generate secrets
> in the future.

It demonstrates that you thought about this problem domain, recognize the trade
off and are saving you and the team time by not implementing it.

This is the question we expect candidates to figure out on their own:

> What version of Go should I use? What dependency manager should I use?

Unless specified in the requirements, pick the solution that works best for you.

### Tools

This task should be implemented in Go or Rust and should work on 64-bit Linux machines.

### Timing

You can split coding over a couple of weekdays or weekends and find time to ask
questions and receive feedback.

Once you join the Slack channel, you have between 1 to 2 weeks complete the
challenge depending on the challenge you choose.

Within this timeframe, we don't give higher scores to challenges submitted more
quickly. We only evaluate the quality of the submission.

We only start the coding challenge if there are several open positions available.

## Key Observations During the Interview
The following are some things that I observed from the team during the
interview process:

* The team doesn't like mocking APIs in order to perform unit testing
  * https://github.com/adalton/teleport-exercise/pull/1#discussion\_r768012202
  * https://github.com/adalton/teleport-exercise/pull/2#discussion\_r771454460

* The team doesn't like unit testing gRPC servers at all; they prefer using
  only full-blown integration tests.
  * https://github.com/adalton/teleport-exercise/pull/2#discussion\_r771448935
  * https://github.com/adalton/teleport-exercise/pull/2#discussion\_r771472278

* In tests, the team want to assert everything that could possibly go wrong
  instead of an Arrange/Act/Assert-based approach that asserts only the things
  a test is designed to prove:
  * https://github.com/adalton/teleport-exercise/pull/2#discussion\_r771451513

## Post-Interview Feedback and Response

I was not offered a position.  The rational provided by the hiring manager was:

> The main reason was complexity of the implementation (had many layers,
> heavy OO influence, prone to panic) which made it difficult to review,
> reason about, and have confidence in accuracy of implementation.
> For example, the panel was still unsure if io.ByteStream.Close() and related
> logic was correct.
>
> Positives we would like to highlight: strong knowledge of cgroups and Linux
> process execution model, good communication skills, and overall
> professionalism during the process.

* I have a strong OO background, so I do tend to apply OO design principles
  (e.g., I create a collection of loosely-coupled components, each of which
  exposes a well-defined API and solves a small piece of the larger problem).
  I unit test those small components in isolation, creating interfaces and mocks
  where necessary to enable that isolation.

* I suspect that the "many layers" concern stems from my design approach and
  my desire to enable unit testing components in isolation.  The team did not
  like my approach unit testing (and I think my approach aligns pretty well
  with established best practices).  I'll also note that the final
  implementation did not deviate significantly from my initial design.  The
  next subsection includes a link to the design document with a full revision
  history.

* I'm not sure about the "prone to panic" comment.  One review did comment about
  a single instance where calling `concreteJob.Stop()` immediately after calling
  `concreteJob.Start()` triggered a panic, and that was something that I fixed
  fairly easily.  It's possible that the reviewers found other similar errors
  and didn't communicate them to me.

* The implementation of `io.ByteStream.Close()` is the most conceptually
  difficult portion of my implementation.  The `ByteStream` component is
  responsible for streaming process output for one client; there may be 0 or
  more instances for any job.  The component maintains a goroutine that monitors
  an `OutputBuffer` for changes, and delivers newly-written data to a channel
  so that it can be streamed to a reader.  If a reader is interrupted early
  (i.e., if the `ByteStream` should stop early), then the implementation calls
  `Close()` to avoid leaking the goroutine.

  The implementation of `Close()` (1) sets a closed flag to let the goroutine
  know that it should stop (if one is running), (2) determines if the goroutine
  was started and if not closes the channel and returns immediately, and
  (3) if the goroutine was started, waits until the goroutine terminates by
  monitoring the channel and draining anything written to it (which in reality,
  would be at most one write).  There may well be a more elegant solution that
  I did not think of, but ignoring blank lines and comments, my contains only
  14 lines of code.  I included detailed comments to try to help readers
  understand the behavior.

## Design Document
The design document for this project, along with its change history, can be
found at:

https://docs.google.com/document/d/1hwHLGNZ25PtIlB9fvSlNQ3XJN5Ce-L0qLHxjZlG4eyY

## Source Tree Organization

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


## Notes on Certificates
The `certs` directory contains some test certificates with which we can
experiment.  These are _not_ production certificates and must not be used
as such.

* ca.cert.pem 
  The root CA used to sign the valid certs

* badca.cert.pem 
  A root CA that was used to sign none of the certs

* server.cert.pem, server.key.pem 
  A valid server certificate/key pair

* administrator.cert.pem, administrator.key.pem 
  A valid administrator certificate/key pair; userID = client1

* client1.cert.pem, client1.key.pem 
  A valid client certificate/key pair, userID = client1

* client2.cert.pem, client2.key.pem 
  A valid client certificate/key pair, userID = client2

* client3.cert.pem, client3.key.pem 
  A valid client certificate/key pair, userID = client3

* weakclient.cert.pem, weakclient.key.pem 
  A client cert/key pair that is too weak

* weakserver.cert.pem, weakserver.key.pem 
  A server cert/key pair that is too weak

* badclient.cert.pem, badclient.key.pem 
  A client cert/key that was not signed by the included root CA

* badserver.cert.pem, badserver.key.pem 
  A server cert/key that was not signed by the included root CA

## Sample Execution
* Start the server
  ```
  $ sudo ./build/jobmanager
  ```

* List jobs. Expect none.
  ```
  $ ./jobctl list
  There are no jobs
  ```

* Start a short-lived job
  ```
  $ ./jobctl start -j date -c $(which date)
  +------+--------------------------------------+
  | NAME |                  ID                  |
  +------+--------------------------------------+
  | date | eea9ab73-b726-4b12-a1ff-5124bf4e53db |
  +------+--------------------------------------+
  ```

* Verify that the job completed successfully
  ```
  $ ./jobctl query eea9ab73-b726-4b12-a1ff-5124bf4e53db
  +------+--------------------------------------+---------+--------+-----------+--------+-------+
  | NAME |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL | ERROR |
  +------+--------------------------------------+---------+--------+-----------+--------+-------+
  | date | eea9ab73-b726-4b12-a1ff-5124bf4e53db | false   | 416111 |         0 |        |       |
  +------+--------------------------------------+---------+--------+-----------+--------+-------+
  ```

* Verify that the job shows up in the list
  ```
  $ ./jobctl list
  +------+--------------------------------------+---------+--------+-----------+--------+-------+
  | NAME |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL | ERROR |
  +------+--------------------------------------+---------+--------+-----------+--------+-------+
  | date | eea9ab73-b726-4b12-a1ff-5124bf4e53db | false   | 416111 |         0 |        |       |
  +------+--------------------------------------+---------+--------+-----------+--------+-------+
  ```

* View the output generated by the job to standard output
  ```
  $ ./jobctl stream eea9ab73-b726-4b12-a1ff-5124bf4e53db
  Mon Dec 20 22:27:06 EST 2021
  ```

* Try to start a job with the same name, expect failure
  ```
  $ ./jobctl start -j date -c $(which date)
  Error: rpc error: code = AlreadyExists desc = job exists
  Usage:
    jobctl start [flags]
  
  Examples:
  start -j myJob -c /usr/bin/find -- /dir -type f
  
  Flags:
    -c, --command string   The command for the job to run; must supply full path
    -h, --help             help for start
    -j, --jobName string   The name of the job to create; must be unique
  
  Global Flags:
        --hostPort string   The <hostName>:<portNumber> of the jobmanager server (default ":24482")
    -u, --userID string     The name of the user (selects the appropriate credential) (default "client1")
  ```

* Start a long-lived job
  ```
  $ ./jobctl start -j longrunning -c /bin/bash -- -c 'for ((i = 0; i < 60; ++i)); do echo $i; sleep 1; done'
  +-------------+--------------------------------------+
  |    NAME     |                  ID                  |
  +-------------+--------------------------------------+
  | longrunning | 91655432-b9ed-49d7-8b60-57ac8f5c7eff |
  +-------------+--------------------------------------+
  ```

* Verify the job is still running
  ```
  $ ./jobctl query 91655432-b9ed-49d7-8b60-57ac8f5c7eff
  +-------------+--------------------------------------+---------+--------+-----------+--------+-------+
  |    NAME     |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL | ERROR |
  +-------------+--------------------------------------+---------+--------+-----------+--------+-------+
  | longrunning | 91655432-b9ed-49d7-8b60-57ac8f5c7eff | true    | 416259 |           |        |       |
  +-------------+--------------------------------------+---------+--------+-----------+--------+-------+
  ```

* Verify that all jobs shows up in the list
  ```
  $ ./jobctl list
  +-------------+--------------------------------------+---------+--------+-----------+--------+-------+
  |    NAME     |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL | ERROR |
  +-------------+--------------------------------------+---------+--------+-----------+--------+-------+
  | date        | eea9ab73-b726-4b12-a1ff-5124bf4e53db | false   | 416111 |         0 |        |       |
  | longrunning | 91655432-b9ed-49d7-8b60-57ac8f5c7eff | true    | 416259 |           |        |       |
  +-------------+--------------------------------------+---------+--------+-----------+--------+-------+
  ```

* Watch the command generate output, and the stream terminate when the command
  finishes
  ```
  $ ./jobctl stream 91655432-b9ed-49d7-8b60-57ac8f5c7eff
  0
  1
  2
  3
  ...
  57
  58
  59
  ```

* Start another long-running job
  ```
  $ ./jobctl start -j tobekilled -c $(which sleep) 1h
  +------------+--------------------------------------+
  |    NAME    |                  ID                  |
  +------------+--------------------------------------+
  | tobekilled | 6c3986f7-9628-40c2-bec8-b9552c36e415 |
  +------------+--------------------------------------+
  ```

* Get the job's pid
  ```
  $ ./jobctl query 6c3986f7-9628-40c2-bec8-b9552c36e415
  +------------+--------------------------------------+---------+--------+-----------+--------+-------+
  | NAME       |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL | ERROR |
  +------------+--------------------------------------+---------+--------+-----------+--------+-------+
  | tobekilled | 6c3986f7-9628-40c2-bec8-b9552c36e415 | true    | 416955 |           |        |       |
  +------------+--------------------------------------+---------+--------+-----------+--------+-------+
  ```

* Stop the job
  ```
  $ ./jobctl stop 6c3986f7-9628-40c2-bec8-b9552c36e415
  ```

* Get the job's status.  Note that it now has a value in the `SIGNAL` column
  ```
  $ ./jobctl query 6c3986f7-9628-40c2-bec8-b9552c36e415
  +------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  | NAME       |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL |      ERROR       |
  +------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  | tobekilled | 6c3986f7-9628-40c2-bec8-b9552c36e415 | false   | 416955 |           | killed | [signal: killed] |
  +------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  ```

* Start a job as a different non-admin user.  Here I'll use the same name as
  the first user -- that's OK.
  ```
  $ ./jobctl -u client2 start -j date -c $(which date)
  +------+--------------------------------------+
  | NAME |                  ID                  |
  +------+--------------------------------------+
  | date | b56b7ab5-bc70-47c3-9182-8c037f4a8892 |
  +------+--------------------------------------+
  ```

* Verify that `client1` cannot see `client2`'s job
  ```
  $ ./jobctl list
  +-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  |    NAME     |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL |      ERROR       |
  +-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  | date        | eea9ab73-b726-4b12-a1ff-5124bf4e53db | false   | 416111 |         0 |        |                  |
  | longrunning | 91655432-b9ed-49d7-8b60-57ac8f5c7eff | false   | 416259 |         0 |        |                  |
  | tobekilled  | 6c3986f7-9628-40c2-bec8-b9552c36e415 | false   | 416955 |           | killed | [signal: killed] |
  +-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  ```

* Verify that `administrator` can see all the jobs.  Note that the administrator's
  view includes the owner.
  ```
  $ ./jobctl -u administrator list
  +---------+-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  |  OWNER  |    NAME     |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL |      ERROR       |
  +---------+-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  | client1 | date        | eea9ab73-b726-4b12-a1ff-5124bf4e53db | false   | 416111 |         0 |        |                  |
  | client1 | longrunning | 91655432-b9ed-49d7-8b60-57ac8f5c7eff | false   | 416259 |         0 |        |                  |
  | client1 | tobekilled  | 6c3986f7-9628-40c2-bec8-b9552c36e415 | false   | 416955 |           | killed | [signal: killed] |
  | client2 | date        | b56b7ab5-bc70-47c3-9182-8c037f4a8892 | false   | 416667 |         0 |        |                  |
  +---------+-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  ```

* Start a long-running job as client2
  ```
  ./jobctl -u client2 start -j longrunning -c $(which sleep) 1h
  +-------------+--------------------------------------+
  |    NAME     |                  ID                  |
  +-------------+--------------------------------------+
  | longrunning | f81b22af-bb5d-4e71-b672-1536b09f7d23 |
  +-------------+--------------------------------------+
  ```

* Ensure `client1` cannot stop `client2`'s job
  ```
  $ ./jobctl stop f81b22af-bb5d-4e71-b672-1536b09f7d23
  Error: rpc error: code = NotFound desc = job not found
  Usage:
    jobctl stop [flags]
  
  Examples:
  jobctl stop 8de11b74-5cd9-4769-b40d-53de13faf77f
  
  Flags:
    -h, --help   help for stop
  
  Global Flags:
        --hostPort string   The <hostName>:<portNumber> of the jobmanager server (default ":24482")
    -u, --userID string     The name of the user (selects the appropriate credential) (default "client1")
  ```

* Ensure that `administrator` can stop `client2`'s job
  ```
  $ ./jobctl -u administrator stop f81b22af-bb5d-4e71-b672-1536b09f7d23
  $ ./jobctl -u client2 query f81b22af-bb5d-4e71-b672-1536b09f7d23
  +-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  |    NAME     |                  ID                  | RUNNING |  PID   | EXIT CODE | SIGNAL |      ERROR       |
  +-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  | longrunning | f81b22af-bb5d-4e71-b672-1536b09f7d23 | false   | 417450 |           | killed | [signal: killed] |
  +-------------+--------------------------------------+---------+--------+-----------+--------+------------------+
  ```

* Verify we can stream stdout
  ```
  ./jobctl start -j ip -c $(which ip)
  +------+--------------------------------------+
  | NAME |                  ID                  |
  +------+--------------------------------------+
  | ip   | 53c05283-3054-4cba-8665-813d4458fcee |
  +------+--------------------------------------+
  
  ./jobctl stream -s stderr 53c05283-3054-4cba-8665-813d4458fcee
  Usage: ip [ OPTIONS ] OBJECT { COMMAND | help }
         ip [ -force ] -batch filename
  where  OBJECT := { address | addrlabel | fou | help | ila | l2tp | link |
                     macsec | maddress | monitor | mptcp | mroute | mrule |
                     neighbor | neighbour | netconf | netns | nexthop | ntable |
                     ntbl | route | rule | sr | tap | tcpmetrics |
                     token | tunnel | tuntap | vrf | xfrm }
         OPTIONS := { -V[ersion] | -s[tatistics] | -d[etails] | -r[esolve] |
                      -h[uman-readable] | -iec | -j[son] | -p[retty] |
                      -f[amily] { inet | inet6 | mpls | bridge | link } |
                      -4 | -6 | -I | -D | -M | -B | -0 |
                      -l[oops] { maximum-addr-flush-attempts } | -br[ief] |
                      -o[neline] | -t[imestamp] | -ts[hort] | -b[atch] [filename] |
                      -rc[vbuf] [size] | -n[etns] name | -N[umeric] | -a[ll] |
                      -c[olor]}
  ```
