---
title: "Spec"
date: 2022-09-25T20:46:41+03:00
weight: 2
---

Create bot spec yaml file `bot.yml`:
```yml
bot:
  handlers:
    - on:
        message:
          command: start
      reply:
        - message:
            text: Hello
            markup:
              keyboard:
                - ["Hello!"]
                - ["How are you?"]
    - on:
        message: "Hello!"
      reply:
        - message:
            text: "Hi!"
    - on:
        message: "How are you?"
      reply:
        - message:
            text: "I'm fine, thank you"
```
This bot declaration specify three handlers:
 - When receive command `/start` (this message received by bot automatically
 when user starts the bot) - and bot respond to it with "Hello" message response
 and sets two keyboard buttons for user: "Hello" to answer hello and "How are you"
 to ask bot about his mood.
 - On receiving "Hello" message from user (if user press "Hello" button) the bot
 respond with message "Hi!"
 - On "How are you?" question the bot answers "I'm fine, thank you" message.

You can validate this declaration:
```bash
$ openbots-validator -config bot.yml 
Config is valid
```
