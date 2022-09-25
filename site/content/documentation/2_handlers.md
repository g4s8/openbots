---
title: "Handlers"
date: 2022-09-25T19:37:45+03:00
weight: 2
---

Handlers declare a bot message handler: it specify which conditions of
incoming message should be matched to perform some action.

The simple handler may look like this:
```yaml
bot:
  handlers:
    - on:
        message:
          text: "Hello"
      reply:
        - message:
            text: "Hi"
```

This handler reacts to user's message with text "Hello" by replying with answer "Hi".

Handler can declare these elements:
 - Trigger `on` (required) - an event matcher used as condition to run handler action.
 - Reply `reply` (optional) - the answers to incoming message.
 - Webhook `webhook` (optional) - call remote URL on message.
 - State `state` (optional) - change current state for the user.
 - Context `context` (optional) - change current context.

The trigger is always required and at least one action must be declared in handler.

And example of working handler:
```yaml
bot:
  handlers:
    on:
      command: start
    reply:
      - message:
        text: "Welcome to my bot!"
      - message:
        text: "You can use `/help` command for help"
    webhook:
      url: "https://stats.myserver.com/bot1"
      method: POST
      data: '{"event": "user_started"}'
    state:
      set:
        new_user: "true"
    context:
      set: onboarding
```
