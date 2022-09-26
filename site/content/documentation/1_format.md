---
title: "Format"
date: 2022-09-25T19:37:45+03:00
weight: 1
---

Bot config file should be valid [YAML](https://yaml.org/) file
with `bot` root element.

The `bot` mapping:
 - May include `token` value with Telegram bot token 
 `token` field.
 - May include `state` element to define bot global state.
 - Must include `handlers` list to define message handlers.

```yaml
bot:
  token: "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" # optional
  state: # optional
    admin_id: '1234'
    greeting: 'Hello'
  handlers: # required
    # ...
```

## Token

The `token` element may specify bot token value, it should be YAML string.
This `token` can be overridden by `BOT_TOKEN` environment variable if present.
It's recommended to use environment variable instead of `token` field,
but this field could be used for bot testing and prototyping.

## State

The state will be discussed later in this docs, for global state you need to know
that it's optional field, it's shared for all users,
and must be specified as key-value pairs, where value are YAML string.

## Handlers

Handlers will be discussed in the next topic.
