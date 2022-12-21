---
title: "Payments"
date: 2022-12-21T22:37:45+03:00
weight: 8
---

Bot could be configured to accept payments.

Read first about payments API, register payment provider for your bot,
and obtain payment provider token: [Telegram bot payments](https://core.telegram.org/bots/payments#the-payments-api).

The payment flow consist of multiple steps, bot handles three of these steps:
 1. Bot sends an invoice to user with details.
 2. User opens the invoice, enter payment details, and clicks "Pay".
 3. Telegram back-end communicates with payment provider API, and sends pre-checkout query to bot.
 4. Telegram bot decides to either accept or decline purchase, and respond with decision.
 5. Telegram back-end sends bot's decision to payment provider, and if decision is "OK", and provider
 is good with it too, then telegram sends successful payment status to bot, and automatically notify user
 about success.
 6. Telegram bot handles success (post-checkout) request.

## Bot configuration

To configure payment providers for bot, add additional bot configuration.
It's a list of providers with `name` and `token` fields:
```yaml
bot:
  config:
    paymentProviders:
      - name: stripe
        token: "123456"
```

Then, you can use `stripe` as provider name in bot handlers.

## Invoice handler

Add invoice handler to bot declaration. All fields are required.
 - `provider` - name of provider from config.
 - `title` - invoice message title.
 - `description` - invoice message primary text.
 - `payload` - invoice identifier. By this payload you will be able to handle checkout
 queries later.
 - `currency` - invoice currency.
 - `prices` - breakdown of invoice prices with labels and amounts (amount is integer with cents, e.g. $10 = `1000`).

```yaml
    - on:
        message:
          test: Buy
      reply:
      - invoice:
          provider: stripe # provider name from configuration
          title: Invoice title
          description: This invoice is a test invoice
          payload: invoice-1 # by this payload you can identify checkout queries
          currency: USD
          prices:
            - label: For testing
              amount: 10000
            - label: Fees
              amount: 100
```

## Pre-checkout

When user enters card details and click pay for the invoice, telegram sends pre-checkout query to bot.
It can be handled by `preCheckout` trigger. It should be answered with `preCheckout` handler.
**Important: pre-checkout query doesn't have any chat association, it's not possible to reply
with some message or anything else expect pre-checkout response**.

Pre-checkout trigger should specify `invoicePayload` identifier.
The response should set `ok` field to either `true` or `false` (to accept or decline checkout).
If it's set to `false`, then `error` field should be provided to describe the reason of decline.

```yaml
    - on:
        preCheckout:
          invoicePayload: invoice-1 # the same string as for invoice payload
      reply:
        - preCheckout:
            ok: true
            # error: "Sorry, our products were ended."
```

## Post-checkout

Post-checkout queries could be handled by `postCheckout` trigger, it should specify
the same `invoicePayload` as for invoice or pre-checkout.

This request belongs to some chat, so it's possible to reply with messages to it:
```yaml
    - on:
        postCheckout:
          invoicePayload: invoice-1
      reply:
        - message:
            text: Success!
```

## Example

This is a full example of a payment flow for the bot which sends secret message for $1 payment:
```yaml
bot:
  debug: true
  config:
    paymentProviders:
      - name: stripe
        token: "1234"
  handlers:
    - on:
        message:
          command: start
      reply:
      - invoice:
          provider: stripe
          title: Secret message
          description: Buy secret message
          payload: secret-message
          currency: USD
          prices:
            - label: Message
              amount: 100
    - on:
        preCheckout:
          invoicePayload: secret-message
      reply:
        - preCheckout:
            ok: true
    - on:
        postCheckout:
          invoicePayload: secret-message
      reply:
        - message:
            text: "It is a secret message"
```
