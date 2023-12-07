---
title: "Managin Secrets"
date: 2023-12-07T19:40:05+04:00
weight: 50
menuTitle: "Secrets"
---

Secrets provide a secure way to store sensitive information in your bot configuration.
These values can be accessed but not modified from the bot specification.
For the cloud version, secrets are managed through the settings tab on the bot's page.
In the self-hosted version, additional configuration is required for the secrets provider.

## Accessing Secrets

Secret variables can be accessed via interpolators in various parts of your bot configuration.
For instance, to include an authorization header in a webhook:

```yml
- on:
    message: test
  reply:
    - message:
        text: Ok
    webhook:
      url: https://example.com/webhook
      method: POST
      headers:
        Authorization: "Bearer ${secret.apikey}"
      data:
        value: "${state.value}"
```

In this example, the `Authorization` header is set with the value of the `apikey` secret.

## Using Go Templates

When working with Go templates, secrets can be accessed through the `.Secrets` object:

```yml
message:
  template: go
  text: "Secret value is {{.Secrets.Val}}"
```

In this case, the template retrieves and displays the value of the Val secret.

Note: Ensure that secrets containing sensitive information are handled securely and accessed only when necessary.

Explore the use of secrets for secure configuration in your bot.
