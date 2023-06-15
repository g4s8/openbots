---
title: "Webhooks"
date: 2022-09-25T19:37:45+03:00
weight: 5
---

Webhook can be sent as HTTP request to configured URL when handler is called.
For example if you want to count some specific message for statistics
or run some code on your side when user clicks the button. Webhook payloads
are created as JSON object based on cnofiguration.

Webhooks may include state or message fields.

Webhook object has mandatory `url` field and optional `method` and `body`
fields:
 - `url` (required) - the URL to call.
 - `method` (optional, default `GET`) - HTTP method.
 - `body` (optional) - request body as YAML mapping

```yaml
on:
  callback: button-terms-accepted
reply:
  - callback:
      text: Thank you!
webhook:
  url: "https://you-site/terms"
  method: POST
  data:
    test: "Hello"
    message: "${message.text}"
    state: "${state.foo}"
```
This webhook send JSON with static `test` field, and dynamic
`message` and `state` fields: `message` field contains message text value,
and `state` field loads `foo` key from state:
```json
{
  "data": {
    "test": "Hello",
    "message": "User input",
    "state": "one"
  },
  "meta": {
    chat_id: 1234,
    timestamp: "01-01-2023T01:02:03.000Z"
  }
}
```
