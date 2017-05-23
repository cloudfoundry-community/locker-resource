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

* `locker_uri` - specifies the base URI of the locker API
* `username` - specifies the username to use when authenticating to the locker API
* `password` - specifies the password to use when authenticating to the locker API
* `ca_cert` - specifies a PEM-encoded CA certificate for validating
  SSL while communicating with the locker API
* `skip_ssl_validation` - determines if ssl validation is ignored while
  communicating with the locker API
* `lock_pool` **required** - specifies the name of the lock pool to use

## `in` - Get status of a Pool

Retrieves the current status of the lock pool, creates a
`lock` file with the name of who owns the lock currently.
If no one owns the lock, the file will be empty.

### Configuration

None.

## `out` - Lock/Unlock a Pool

Allows you to lock or unlock a pool. If the pool is currently locked,
even by the item requesting the lock, it will loop indefinitely, until
the lock is released. Unlocking is only allowed when the `lock_with` matches
the item that had claimed the lock initially.

* `lock_with` **required** - specifies the name of the item that is attempting to claim
  the lock. 
* `lock_op` **required** - specifies whehter you want to `lock` or `unlock` the pool
