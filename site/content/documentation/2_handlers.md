---
title: "Event Handlers"
date: 2023-12-06T21:16:24+04:00
weight: 20
menuTitle: "Handlers"
---

Handlers in the bot configuration play a pivotal role in defining how your bot responds to different Telegram event updates.
A handler declares various elements that determine its behavior:

## Trigger

The `on` trigger specifies when the handler is triggered.
The trigger element is required for handler.

 * `message`: Triggers the handler on a text message or command.
 * `callback`: Triggers on a button callback (part of inline-buttons and callbacks feature).
 * `context`: Additional selector to trigger the handler only if the current user context is set to a specified value (context feature).
 * `preCheckout`: Triggers on pre-checkout events (payments feature).
 * `postCheckout`: Triggers on post-checkout events (payments feature).
 * `state`: Array of state conditions, an additional filter to run the handler only if the user's state matches these conditions (states feature).

Trigger should have at least one of `message`, `callback`, `context`, `preCheckout`, `postCheckout`,
`state` and `context` could be added to other elements. If trigger is a string, it will be treated as
message handler, there are two identical triggers below:
```yml
bot:
  handlers:
    # trigger message shorthand
  - on: hello
    reply:
    - message: Hello
```
```yml
bot:
  handlers:
    # full message trigger
  - on:
      message: hello
    reply:
    - message: hello
```
```yml
bot:
  handlers:
    # the fullest message trigger version
  - on:
      message:
        text: hello
    reply:
      message: hello
```

The wildcard trigger is handled as a default fallback trigger after trying all other triggers:
```yml
bot:
  handlers:
  - on: "*"
    reply:
      message: Fallback message
```

`message`, `callback`, `preCheckout`, `postCheckout` could not be mixed together.

Example:

```yml
bot:
  handlers:
    - on:
        message:
          command: start
      reply:
        - message:
            text: "Welcome to my bot!"
    - on:
        callback:
          data: "button_click"
      reply:
        - message:
            text: "You clicked a button!"
```

In this example, the first handler triggers on the `/start` command, replying with a welcome message.
The second handler triggers on a button click with the callback data "button\_click" and responds accordingly.

## Reply element

The `reply` element specifies how to reply to a user event. It can include:

 * **message:** Reply with a text message.
 * **callback:** Reply to a button callback with a popup text or alert message.
 * **edit:** Edit the message that triggers this event (e.g., if the message has a button).
 * **delete:** Delete the message that triggers this event.
 * **image:** Reply with an image.
 * **document:** Reply with a document.
 * **invoice:** Reply with an invoice for payment (discussed later as part of the payments feature).
 * **preCheckout:** Reply to a pre-checkout event (also part of the payments feature).

One handler may have multiple different reply items.

Example:

```yml
bot:
  handlers:
    - on:
        message:
          command: start
      reply:
        - message:
            text: "Welcome to my bot!"
        - image:
            name: "Test image"
            key: "test-image.png"
```

In this example, the handler replies with a welcome message and an image when triggered by the `/start` command.

Continue exploring to learn more about the webhook, state, context, data, and validate elements,
enabling you to create versatile and dynamic bot interactions.
