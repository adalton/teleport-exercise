# Teleport Exercise
This repository contains my implementation of the Teleport programming
challenge exercise.  The specification for this project can be found at:

https://docs.google.com/document/d/1hwHLGNZ25PtIlB9fvSlNQ3XJN5Ce-L0qLHxjZlG4eyY

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


## Notes on Certificates
The `certs` directory contains some test certificates with which we can
experiment.

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
