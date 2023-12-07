---
title: "Message Editing and Deleting"
date: 2023-12-07T21:36:39+04:00
weight: 90
menuTitle: "Edit and Delete"
---

Allowing your bot to edit or delete messages based on user interactions enhances the interactive experience.
This feature is particularly useful when users engage with inlineMarkup keyboard buttons or other interactive elements.

## Message editing

Edit a message by using the `edit` element in your bot configuration.
This is often triggered by user interactions such as clicking a button:

```yml
handlers:
  - on:
      message:
        command: start
    reply:
      - message:
          text: "Click the button to see message edits"
          markup:
            inlineKeyboard:
            - - text: "Click"
                callback: inc-counter
  - on:
      callback: inc-counter
    reply:
      - callback:
          text: "Click"
      - edit:
          message:
            text: "You clicked this button ${state.counter} times."
            inlineKeyboard:
              - - text: "Button"
                  callback: inc-counter
      state:
        ops:
          - kind: add
            key: counter
            value: "1"
```

In this example, each button click increments a counter and updates the message text accordingly.

## Message Deletion

Delete a message by using the `delete: true` flag in your bot configuration.
This is often triggered by user interactions such as clicking a button:

```yml
handlers:
  - on:
      message:
        command: start
    reply:
      - message:
          text: "Click the button to delete this message"
          markup:
            inlineKeyboard:
              -
                - text: "Delete"
                  callback: delete
  - on:
      callback: delete
    reply:
      - callback:
          text: "Deleted"
        delete: true
```

This example deletes the message when the "Delete" button is clicked.

Enhance user engagement by leveraging message editing and deleting in response to specific user interactions.
