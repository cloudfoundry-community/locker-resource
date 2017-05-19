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
