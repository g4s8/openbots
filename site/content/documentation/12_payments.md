---
title: "Payments"
date: 2023-12-07T23:02:40+04:00
weight: 120
menuTitle: "Payments"
---

Enable payment transactions in your bot by following the Telegram payment workflow,
which includes three key stages: Invoice, Pre Checkout, and Post Checkout.
Make sure to fulfill the prerequisites, such as registering with a supported payment provider
and obtaining the provider token.

## Prerequisites

 1. **Register with a Payment Provider:** Bot owners should register with any payment provider
 supported by Telegram and obtain the provider token. Refer to the
 [official documentation](https://core.telegram.org/bots/payments) for details.
 2. **Cloud Version Setup:** For the cloud version, navigate to the bot's settings tab and add
 the payment provider token under the Payments section. The token will be securely stored encrypted in the database.
 Remeber token name to use it for invoice provider name.
 3. **Self-Hosted Configuration:** For self-hosted versions, configure payment providers
 as outlined in the advanced configuration documentation.

## Payment Workflow

The Telegram payment workflow consists of three stages:

### Sending Invoice

To send an invoice, include the `invoice` object in your handler. For example:

```yml
bot:
  handlers:
    - on:
        message:
          command: start
      reply:
      - invoice:
          provider: test-stripe
          title: Test invoice
          description: This invoice is a test invoice
          payload: invoice-1
          currency: USD
          prices:
            - label: For testing
              amount: 10000
            - label: Fees
              amount: 100
```

Invoice Object Properties

 * `provider` (required): The name of the payment provider configured for the bot.
 * `title` (required): The title of the invoice.
 * `description`: Detailed description of the invoice.
 * `payload` (required): Payload ID to handle preCheckout and postCheckout events.
 * `currency` (required): Currency code for the invoice.
 * `prices` (required): Array of labels and amounts.

### Pre Checkout

Handle the pre-checkout request using the `preCheckout` object:

```yml
- on:
    preCheckout:
      invoicePayload: invoice-1
  reply:
    - preCheckout:
        ok: true
```

Reply with `ok: true` to continue processing or `ok: false` to stop the checkout.
This stage allows the bot to check if the purchase can be processed, e.g., if products are available.

### Post Checkout

Upon successful completion of the checkout, handle the post-checkout event:

```yml
- on:
    postCheckout:
      invoicePayload: invoice-1
  reply:
    - message:
        text: Success!
```

Here, the bot can process the purchase after confirming the completion of the checkout.
