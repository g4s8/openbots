---
title: "Reply"
date: 2022-09-25T19:37:45+03:00
weight: 4
---

Reply declares one or more replies to message.
It's defined by `reply` element in `handler` and should
contain the list of reply actions.

Reply action could be either `message` or `callback`.

## Message

Message reply action declares answer message for user message.

It contains these elements:
 - Text to reply as `text` (required) - text message to reply
 - Parse mode `parseMode` (optional) - could be empty or one of:
 `Markdown`, `MarkdownV2` or `HTML`.
 - Reply markup as `markup` - specify Telegram markup changes on
 user's chat.

### Parse mode

Parse-mode could be used to render markdown or HTML text on user's device,
for instance:
```yaml
reply:
  - message:
      text: "Just plain text"
  - message:
      text: "Some *bold* text"
      parseMode: Markdown
  - message:
      text: "Telegram specific markdown, e.g. ||spoiler||"
      parseMode: MarkdownV2
  - message:
      text: |
        Code example:
        <pre><code class="language-go">func main() {
          fmt.Println("Hello world")
        }</code></pre>
      parseMode: HTML
```

### Markup

Now add some markup to our messages, it could be custom chat
keyboard or inline buttons for reply message.
Chat keyboard could be specified by
`keyboard` element of array of arrays of strings with button names (first array for
rows, second for columns). Inline buttons has more complex structure but it's also
specified as array of arrays.

Example:
```yaml
reply:
  - message:
      text: Where you go?
      markup:
        keyboard:
          - ["Up"]
          - ["Left", "Right"]
          - ["Down"]
```
changes user's keyboard in chat with 3-rows keyboard:
first row has one button "Up", second row two buttons "Left" and "Right",
and 3-d row one button "Down".

Inline keyboard:
```yaml
reply:
  - message:
      text: How to contact you?
      markup:
        inlineKeyboard:
          -
            - text: Email
              callback: contact-me-email
            - text: Phone
              callback: contact-me-phone
          -
            - text: Other
              url: https://contact-me/form
              
```
The first row has two buttons: first sends callback data `contact-me-email` on click,
the second sends `contact-me-phone`. Each callback could be processed by another handler.
The second row has one button which opens website URL `https://contact-me/form`.

## Callbacks

Now we know how to send callback messages from buttons and
on previous page we discussed how to handle these messages.

Telegram allows us to send popup messages to callback events.
It can be either popup at the top of chat window or alert with "OK" button.

Callback reply is declared with `callback` element of `reply` object it can be configured
by two parameters:
 - `text` string (required) with text to popup
 - `alert` bool (optional, default `false`) to show alert instead of popup

```yaml
handlers:
  - on:
      callback: contact-me-email
      reply:
        - callback:
            text: "We'll send you email soon"
        - message:
            text: "Please check your inbox in five minutes"
  - on:
      callback: contact-me-phone
      reply:
        - callback:
            text: "We'll call you in five minutes"
            alert: true
```

## Edit message

You can edit message with inline button by using `edit` element, it allows to
edit caption, text and inline keyboard of message.

*There are some restrictions on edit messages*:
 - Caption can't be used together with text or inlineKeyboard edits
 - Text edits and inlineKeynoard edits could be used in one edit reply
 - If you don't specify inlineKeyboard for edits, original buttons will be removed on edit

```yaml
- on:
    message:
      command: start
  reply:
    - message:
        text: "Click the button to see message edits"
        markup:
          inlineKeyboard:
            -
              - text: "Click"
                callback: callback-1
- on:
    callback: callback-1
  reply:
    - callback:
        text: "Click"
    - edit:
        message:
          text: "New text for message"
          inlineKeyboard:
            -
              - text: "Button"
                callback: callback-2
```

## Delete message

Callback handler can delete message of the callback button, just put `delete: true`
to reply handler:

```yaml
bot:
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

## Reply image

Bot can reply with image content. Image handler can be configured using these fields:
 - `name` (requied, string) - image name
 - `file` (requied, string) - image file path on file system

Example:
```yaml
- on:
    message:
      command: start
  reply:
    - image:
        name: Test image
        file: /tmp/assets/test-image.png
```
