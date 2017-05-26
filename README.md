# What is Locker?

Locker is a simple web application for claiming and releasing locks.
It contains a Concourse resource to help make locking various tasks
and jobs in Concourse pipelines easier to manage.

This project is similar to the [pool-resource](https://github.com/concourse/pool-resource), with a few key differences:

1) It uses a webserver with file-backed persistence for the locks. This isn't
   terribly scalable, but since it's handling locks for Concourse pipelines,
   it isn't anticipated to have an enormous traffic load.
2) Pools + lock-names do not need to be pre-defined in order to be used

# How do I use it?

1) Deploy `locker` along side your Concourse database node via the [locker-boshrelease](https://github.com/cloudfoundry-community/locker-boshrelease)
2) Use the `locker-resource` in your Concourse pipelines (see below for resource
   configuration details).

# locker-resource

## Source Configuration

* `locker_uri` **required** - specifies the base URI of the locker API
* `username` - specifies the username to use when authenticating to the locker API
* `password` - specifies the password to use when authenticating to the locker API
* `ca_cert` - specifies a PEM-encoded CA certificate for validating
  SSL while communicating with the locker API
* `skip_ssl_validation` - determines if ssl validation is ignored while
  communicating with the locker API
* `lock_name` **required** - specifies the name of the lock to be acquired

## `in` - Get status of a Pool

Retrieves the current status of the lock pool, creates a
`lock` file with the name of who owns the lock currently.
If no one owns the lock, the file will be empty.

### Configuration

None.

## `out` - Lock/Unlock a Pool

Allows you to lock or unlock a lock. If the pool is currently locked,
even by the item requesting the lock, it will loop indefinitely, until
the lock is released. Unlocking is only allowed when the `key` matches
the key that has currently locked the lock.

* `key` **required** - specifies the name of the item that is attempting to claim
  the lock.
* `lock_by` -  specifies the item requesting to lock/unlock the lock with the key
* `lock_op` **required** - specifies whether you want to `lock` or `unlock` the pool


### Example of exclusive locking between two jobs

```
resource:
- name: myLock
  source:
    locker_uri: http://10.10.10.10:8910
    username: test
    password: test
    lock_name: myLock

jobs:
- name: exclusive-job-1
  public: true
  serial: true
  plan:
  - put: myLock
    params:
      key: exclusive-job-1
      lock_op: lock
  - task: do_stuff
    .....
  - put: myLock
    params:
      key: exclusive-job-1
     lock_op: unlock
- name: exclusive-job-2
  public: true
  serial: true
  plan:
  - put: myLock
    params:
      key: exclusive-job-2
      lock_op: lock
  - task: do_stuff
    .....
  - put: myLock
    params:
      key: exclusive-job-1
     lock_op: unlock
```

### Example of shared locking between some jobs that cannot co-exist with another

```
resource:
- name: myLock
  source:
    locker_uri: http://10.10.10.10:8910
    username: test
    password: test
    lock_name: myLock

jobs:
# coexisting-1 can run if coexisting-2 is running, but not exclusive-job
- name: coexisting-1
  public: true
  serial: true
  plan:
  - put: myLock
    params:
      key: can-coexist
      locked_by: coexisting-1
      lock_op: lock
  - task: do_stuf
    ...
  - put: myLock
    params:
      key: can-coexist
      locked_by: coexisting-1
      lock_op: unlock
# coexisting-2 can run if coexisting-1 is running, but not exclusive-job
- name: coexisting-2
  public: true
  serial: true
  plan:
  - put: myLock
    params:
      key: can-coexist
      locked_by: coexisting-2
      lock_op: lock
  - task: do_stuf
    ...
  - put: myLock
    params:
      key: can-coexist
      locked_by: coexisting-2
      lock_op: unlock
# exclusive-job *cannot* run if coexisting-1 or coexisting-2 are running
- name: exclusive-job
  public: true
  serial: true
  plan:
  - put: myLock
    params:
      key: cannot-coexist
      locked_by: exclusive-job
      lock_op: lock
  - task: do_stuff
    ...
  - put: myLock
    params:
      key: cannot-coexist
      locked_by: exclusive-job
      lock_op: unlock
```
