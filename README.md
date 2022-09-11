Telegram Bot framework with declarative specification.

[![CI](https://github.com/g4s8/openbots/actions/workflows/go.yml/badge.svg)](https://github.com/g4s8/openbots/actions/workflows/go.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/g4s8/openbots)](https://hub.docker.com/r/g4s8/openbots)

## Example

Write declarative bot specification:
```yaml
# bot.yml
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
Run your bot:
```sh
docker run \
    -v $PWD/bot.yml:/w/config.yml \
    --env BOT_TOKEN="$BOT_TOKEN" \
    g4s8/openbots:latest
```

## Quick start

TBD: quick-start wiki and documentation

## About

This is a Telegram bot framework which allows you to write low-code bot project. You delcare the bot in YAML specification
file and start the bot using CLI or Docker image.

Full feature list:
 - [x] handle text messages
 - [x] handle bot commands
 - [x] handle inline queries callbacks (buttons)
 - [x] reply with text messages
 - [x] reply callbacks
 - [x] reply with inline buttons
 - [x] change keyboard layout (reply markup)
 - [x] reply with Markup, MarkupV2, HTML messages
 - [x] switch context, handle context-based updates
 - [x] keep state data and interpolate state in replies
 - [ ] edit message
 - [ ] reply with images
 - [ ] handle image messages
 - [ ] delete messages
 - [ ] send message to particular user
 - [ ] call webhook on update
 - [ ] ... more features
 
## Extending

This bot could be extended with custom handlers on Go.
