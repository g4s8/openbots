---
title: "Database"
date: 2022-09-25T19:37:45+03:00
weight: 1
---

# Database configuration

When using `database` persistence configuration, then `db_config` must be specified
under `persistence`. It has these properties:
 - `user` (required, string) - database user name
 - `password` (required, string) - database user password
 - `host` (required, string)-  database host
 - `port` (required, int) - database port
 - `database` (required, string) - database name
 - `no_ssl` (optional, bool) - set `true` to disable SSL

For example, to use database `mybot` on `localhost` 5432 port as `john` as user and `qwerty` as password,
and disable SSL for connection, use:
```yaml
bot:
  config:
    persistence:
      type: database
      db_config:
        user: john
        password: qwerty
        host: localhost
        port: 5432
        database: mybot
        no_ssl: true
```
