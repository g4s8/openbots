---
title: "Dynamic bot"
date: 2022-10-15T21:04:36+03:00
weight: 5
---

Here is simple dynamic bot which will ask your name and remember it.

```yaml
bot:
  handlers:
    - on:
        message:
          command: start
      reply:
        - message:
            text: "What is your name?"
      context:
        set: ask-name
    - on:
        context: ask-name
      state:
        set:
          name: "${message.text}"
      context:
        delete: ask-name
      reply:
        - message:
            text: "Now I remember your name."
    - on:
        message:
          text: "Hello"
      reply:
        - message:
            text: "Hello ${state.name}!"
```

Let's check each handler one by one:
```yaml
    - on:
        message:
          command: start
      reply:
        - message:
            text: "What is your name?"
      context:
        set: ask-name
```

This handler is triggered by `/start` command. It's automatically sent when you start telegram bot,
or you can call it manually (just type `/start` to bot).

On receiving this command, bot reply with "What is your name?" question, and changing current context
to `ask-name`. This context will be used to filter out other messages.

```yaml
    - on:
        context: ask-name
      state:
        set:
          name: "${message.text}"
      context:
        delete: ask-name
      reply:
        - message:
            text: "Now I remember your name."
```
This handler performed if current context is `ask-name`, it doesn't matter what message or command
will be sent to the bot. After previous handler the context will be switched to `ask-name`, so this handler
will be triggered on any message after replying to "What is your name?" question.

On message, it saving message text from user (`"${message.text}"`) into state with key `name`.

And reply with a message "Now I remember your name."

Also, it reset current context `ask-name`, so after this handler context will be empty and it'll not be triggered
again, until another handler set context to `ask-name` back.

```yaml
    - on:
        message:
          text: "Hello"
      reply:
        - message:
            text: "Hello ${state.name}!"
```

And last handler in this example answers on "Hello" messages. Now it knows user's name and can respond with
greeting: `Hello ${state.name}`, where `${state.name}` will be replaced with state value of `name` key, which was set
on previous step.
