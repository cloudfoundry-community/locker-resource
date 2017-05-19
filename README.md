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

FIXME: Coming soon


# locker API

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
