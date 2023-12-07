---
title: "Reply Message Markups"
date: 2023-12-07T20:59:24+04:00
weight: 70
menuTitle: "Markups and Keyboard"
---

Telegram provides two types of reply markups that allow you to enhance user interactions:
"keyboard" and "inlineKeyboard." These markups enable you to customize the keyboard layout
or attach buttons to messages with callback data or external URLs.

## Chat Keyboard Markup

Chat keyboard markup is represented as an array of arrays of strings,
defining the layout of buttons. The first array represents rows,
and the nested arrays represent buttons in each row.

**Example:**
```yml
- on:
    message:
      command: start
  reply:
    - message:
        text: Hello
        markup:
          keyboard:
            - ["One", "Two"]
            - ["Three"]
            - ["Four", "Five", "Six"]
```

In this example, the bot replies with a "Hello" message and sets a chat keyboard
with the following layout:

```txt
[   One   ] [   Two   ]
[       Three         ]
[ Four ][ Five ][ Six ]
```

**Note:** chat keyboard buttons send the exact text displayed on the button when clicked.

## Inline Keyboard Markup

Inline keyboard markup attaches buttons to the current message.
Each button can include callback data or an external URL.

**Example:**

```yml
- on:
    message: test
  reply:
    - message:
        text: Test
        markup:
          inlineKeyboard:
            - - callback: callback-one
                text: One
            - - text: Open link
                url: https://example.com
```

In this example, the bot replies with a "test" message and attaches
an inline keyboard with two buttons.

**Note:** inline buttons send callback data to the bot when clicked.

### Handling Inline Button Clicks

When an inline button with callback data is clicked,
the bot can handle the click using the callback event.

```yml
- on:
    callback: callback-one
  reply:
    - message:
        text: "Button 'One' clicked"
```

It is recommended to reply with a callback reply to inform the user that the callback was handled.
This ensures a smooth user experience.

```yml
- on:
    callback: callback-one
  reply:
    - callback:
        text: "Button 'One' clicked"
        alert: true
```

In this example, the bot replies with a callback message, showing "Button 'One' clicked" to the user,
and an alert window is displayed.

Callback reply can be configured with `alert: true` for an alert window or
`alert: false` (default) for a non-blocking popup on the current chat screen.

Explore the flexibility of reply message markups to create interactive and user-friendly bot interactions.
