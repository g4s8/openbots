---
title: "User Context Management"
date: 2023-12-07T20:23:25+04:00
weight: 60
menuTitle: "Context"
---

User context is a powerful mechanism in your bot configuration to dynamically
adjust behavior based on the ongoing chat context. Different handlers can be triggered
depending on the current context, enhancing the flexibility of your bot.

## Setting and Clearing Context

To set and clear the user `context` within a handler, the context element is utilized.
This element can have two values, and only one can be used at a time:

 * `set`: Sets the current context to the specified value.
 * `delete`: Clears the current context, leaving it empty.

```yml
handlers:
  - on:
      command: 'delete'
    reply:
      text: 'Are you sure?'
    context:
      set: 'delete-question'
  - on:
      message: 'Yes'
      context: 'delete-question'
    reply:
      message:
        text: 'Deleted'
    context:
      delete: 'delete-question'
  - on:
      message: 'No'
    reply:
      message: 'Canceled'
    context:
      delete: 'delete-question'
```

In this example, when the command 'delete' is received, the context is set to 'delete-question'.
Subsequent messages like 'Yes' or 'No' within the 'delete-question' context trigger specific responses.

## Using Context in Triggers

To use the context in a trigger, simply add the `context` element to the `on` object:

```yml
  - on:
      message: 'Yes'
      context: 'delete-question'
```

This ensures that the handler is only triggered when the user sends 'Yes' within the
'delete-question' context.

Leverage user context to create dynamic and context-aware interactions in your bot.
