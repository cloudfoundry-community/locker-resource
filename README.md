# What is Locker?

Locker is a simple web application for claiming and releasing locks.
It contains a Concourse resource to help make locking various tasks
and jobs in Concourse pipelines easier to manage.

This project is similar to the [pool-resource], with a few key differences:

1) It uses a webserver with file-backed persistence for the locks. This isn't
   terribly scalable, but since it's handling locks for Concourse pipelines,
   it isn't anticipated to have an enormous traffic load.
2) Pools + lock-names do not need to be pre-defined in order to be used

# How do I use it?

1) Deploy `locker` along side your Concourse database node via the [locker-boshrelease]
2) Use the `locker-resource` in your Concourse pipelines (see below for an example pipeline
   + details on configuring the resource)

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

# locker API

## Supported Requests

* `GET /locks`

  Returns a JSON formatted list of locks + who owns them currently

* `PUT /lock/<pool-name>`

  Content: `{"lock":"item-requesting-the-lock"}`

  Issues a lock on `pool-name` to the value of the `lock` attribute in the JSON payload of the request.
  If the lock was already taken, it will immediately return a 423 error, and the client should back-off +
  re-try at a sane interval until the lock is obtained.

  Returns 200 on success, 423 on locking failure

  Example to lock `prod-deployments` with `prod-cloudfoundry`:

  ```
  curl -X PUT -d '{"lock":"prod-cloudfoundry"}' http://locker-ip:port/lock/prod-deployments
  ```

* `DELETE /lock/<pool-name>`

  Content: `{"lock":"item-requesting-unlock"}`

  Issues an unlock request on `pool-name` based on the value of the `lock` attribute in the JSON payload
  of the request. If the lock on `pool-name` is not currently held by `item-requesting-unlock`, the
  unlock is disallowed. If the lock is currently not held by anyone, returns 200.

  Returns 200 on success, 423 on failure.

  Example to unlock `prod-deployments` previously locked by `prod-cloudfoundry`:

  ```
  curl -X DELETE -d '{"lock":"prod-cloudfoundry"}' http://locker-ip:port/lock/prod-deployments
  ```

## Authentication

`locker` is protected with HTTP basic authentication. Credentials are passed in 

## Error handling

Errors will be reported as json objects, inside a top-level `error` key:

```
{"error":"This is what went wrong"}
```

## Running manually

If you want to run `locker` manually for testing/development:

```
go build
LOCKER_CONFIG=/tmp/locker-data.yml ./locker
```

`locker` can be configured using the following environment variables:

* `LOCKER_CONFIG` **required** - specifies the file that locks will be stored in
* `PORT` - Defaults to 3000, controls the port `locker` listens on
* `AUTH_USER` - If specified, requires `AUTH_PASS` and configures the username for
  HTTP basic auth
* `AUTH_PASS` - If specified, requires `AUTH_USER` and configures the password for
  HTTP basic auth
