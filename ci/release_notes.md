# New Features

- `bosh_lock` can be used in place of `lock_name` now. It triggers a special mode to contact a BOSH
  director (URI specified by the value of `bosh_lock`), and use the director's name as the lock name.
