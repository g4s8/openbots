---
title: "API"
date: 2022-10-29T19:37:45+03:00
weight: 7
---

You can declare API handlers in bot spec and call them via HTTP.

E.g. add this code to spec:
```yaml
bot:
  api:
    handlers:
      - id: notify
        actions:
          - send-message:
              text: Hello ${data.text}
            chat-id: 1234 # optional, could be specified in payload request
```
And then call it:
```
POST /handlers/notify


{
  "data": {
    "text":"test"
  }
}
```

This HTTP request will send a message with "test" text to chat `1234`.


## The structure of API handler

Handler must have unique ID - it will be API URL path: `/handlers/$ID`.
Then handlers should have `actions` list, each action may have different behaviors:
 - `send-message` - sends message to chat

The `send-message` action should have `text` argument. It can be either static text:
```yaml
text:
  value: "Hello"
```
Or dynamic parameter parsed from HTTP call:
```yaml
text:
  param: message
```
In second case API caller must provide `message` parameter in payload.


## Calling API handler

API service accepts only `POST` requests at path `/handlers/$ID`, where `$ID` is
handler ID defined in spec.

The JSON body must contain chat ID as `chat_id` number field,
and optional `params` field, where all dynamic parameter should be provided as
strings:
```json
{
  "chat_id": 1234,
  "params": {
    "message": "Hello"
  }
}
```
