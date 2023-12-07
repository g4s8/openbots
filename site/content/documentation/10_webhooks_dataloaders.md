---
title: "Webhooks and Data Loaders"
date: 2023-12-07T22:14:38+04:00
weight: 100
menuTitle: "Webhooks, Loaders"
---

Enhance the capabilities of your bot by incorporating webhooks and data loaders.
These features allow your bot to interact with external services through HTTP calls
and fetch dynamic data to enrich user interactions.

## Webhooks

Configure webhooks in your bot to trigger HTTP calls to external REST APIs.
The `webhook` object allows you to specify the URL, method, headers, and data payload for the HTTP request.

### Webhook Object Elements

 * `url` (required): The URL to send the webhook request.
 * `method` (optional, default: 'GET'): The HTTP method for the webhook request (e.g., 'GET', 'POST', 'PUT').
 * `headers` (optional): Key-value pairs of strings representing HTTP request headers for the webhook.
 * `data` (optional): Key-value pairs of strings representing the JSON payload for the HTTP request body.

**Example:**

```yml
handlers:
  - on:
      message: test-webhook
    reply:
      - message:
          text: Ok
      - webhook:
          url: https://example.com/webhook
          method: POST
          headers:
            Authorization: "Bearer ${secret.apikey}"
          data:
            name: "${state.name}"
            user_id: "${user.id}"
            source: test-webhook
```

In this example, a POST request is sent to `https://example.com/webhook`
with specified headers and data payload.

## Data Loaders

Data loaders enable your bot to fetch external data via REST calls and use it within message templates.
Fetch external data during user interactions to provide dynamic and personalized responses.

### Data Loader Object Elements

 * `method` (optional, default: 'GET'): The HTTP request method for the data loader (e.g., 'GET', 'POST', 'PUT').
 * `url` (required): The external URL to fetch data from.
 * `headers` (optional): Key-value pairs of strings representing HTTP request headers for the data loader.

**Example:**

```yml
bot:
  handlers:
  - on:
      message:
        command: start
    data:
      fetch:
        method: GET
        url: "https://jsonplaceholder.typicode.com/todos/1"
    reply:
      - message:
          text: 'id: {{ index .Data "userId"}}, title: {{ .Data.title }}'
          template: go
```

In this example, the bot fetches JSON data from `https://jsonplaceholder.typicode.com/todos/1`
and uses Go templating to render it into the message text.

Webhooks and data loaders open up possibilities for integrating your bot with external services,
making interactions more dynamic.

