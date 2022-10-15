+++
title = "Persistence"
date = 2022-09-25T19:13:39+03:00
weight = 3
chapter = true
pre = "<b>3. </b>"
+++

# Persistence

By default bot keeps all data in memory. But it's possible to connect Postgres database.

To configure persistence of bot, add `persistence` element to `configuration`:
```yaml
bot:
  config:
    persistence:
      type: memory
```

Persistence `type` may have two values:
 - `memory` (default) - keep all data in memory
 - `database` - use database as data storage

Database storage requires additional configuration.
