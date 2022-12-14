---
title: "Triggers"
date: 2022-09-25T19:37:45+03:00
weight: 3
---

Trigger is defined as `on` element in handler.

The trigger should include at least one condition:
 - Message `message` - check received message data.
 - Callback `callback` - check callback data received from Telegram.
 - Context `context` - additional option for message or callback condition,
 will be discussed later.

# Message trigger

Message trigger can include `text` or `command` conditions:

```yaml
on:
  message:
    text: "Ping"
```
filter messages by exact match of the text to be equal "Ping".

```yaml
on:
  message:
    command: "start"
```
filter messages by its command (e.g. `/start`, `/info`, `/help`). Message command should not include leading slash `/` character, because it's parsed into word by Telegram.

Text handler could be simplified by using string instead of `text` element:
```yaml
on:
  message:
    text: Hello
```
is equal to
```yaml
on:
  message: Hello
```

## Callback trigger

Callback is a special type of message from Telegram which is sent when user
clicks some inline button with callback data. You will learn about these buttons
on next page, now we assume that we can attach an arbitrary string data to each
button and process it with callback handler:
```yaml
on:
  callback: callback-data
  reply:
    - message:
        text: Thanks for clicking!
```

## Context

Context is a kind of local user state which can be changed by handler.
The trigger can include context condition to perform some action only when
context matches:
```yaml
handlers:
  - on:
      command: 'delete'
    reply:
      text: 'Are you sure?'
    context:
      set: delete-question
  - on:
      message: Yes
      context: delete-question
    reply:
      message:
        text: Deleted
    context:
      delete: delete-question
  - on:
      message: No
    reply:
      message: Canceled
    context:
      delete: delete-question
```

These "Yes" and "No" handlers will be triggered only when context is `delete-question`,
which is set to this value only after user command `/delete`.
