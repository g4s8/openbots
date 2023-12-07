---
title: "Bot Configuration Overview"
date: 2023-12-06T20:56:54+04:00
menuTitle: "Overview"
weight: 10
---

The `bot` element is the root configuration for your Telegram chat-bot.
It encapsulates various settings and specifications that define the behavior of your bot.
Below are the primary top-level elements within the bot configuration:

## Handlers

The `handlers` element is a crucial component of your bot configuration.
It represents an array of handlers that respond to different Telegram event updates.
Handlers define how your bot reacts to user messages, button clicks, and other events.
Each handler specifies conditions triggering its execution and the corresponding actions to take.

Example:

```yml
bot:
  handlers:
  - on:
      message:
        command: start
    reply:
    - message:
        text: "Hello! Welcome to your new bot."
```

In this example, a handler responds to the `/start` command with a welcome message.

## API

The `api` element allows you to configure an HTTP API for your bot.
This feature exposes a service API, enabling authorized individuals, such as the bot creator,
to call the API. By making API calls, you can render predefined templates with provided data
and send the resulting content to specified users. This provides a powerful mechanism for dynamic
content generation and external integrations.

Specification example:

```yml
bot:
  api:
    handlers:
    - id: notify
      actions:
      - send-message:
          text: Hello ${data.text}
        chat-id: 1234 # optional, could be specified in payload request
```

HTTP request to API example:

```txt
POST /handlers/notify
Content-Type: application/json

{
  "data": {
    "text": "test"
  }
}
```

This HTTP request will send a message with the text "test" to the chat with ID 1234.

Explore the `api` documentation further to understand the available actions,
handlers, and customization options for integrating dynamic content into your bot.
