---
title: "Webhooks"
date: 2022-09-25T19:37:45+03:00
weight: 5
---

Webhook can be sent as HTTP request to configured URL when handler is called.
For example if you want to count some specific message for statistics
or run some code on your side when user clicks the button.

Webhook object has mandatory `url` field and optional `method` and `body`
fields:
 - `url` (required) - the URL to call.
 - `method` (optional, default `GET`) - HTTP method.
 - `body` (optional) - request body.

```yaml
on:
  callback: button-terms-accepted
reply:
  - callback:
      text: Thank you!
  - webhook:
      url: "https://you-site/terms"
      method: POST
      body: '{"accepted": true}'
```
