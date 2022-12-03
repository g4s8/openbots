Telegram Bot framework with declarative YAML specification.

[![CI](https://github.com/g4s8/openbots/actions/workflows/go.yml/badge.svg)](https://github.com/g4s8/openbots/actions/workflows/go.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/g4s8/openbots)](https://hub.docker.com/r/g4s8/openbots)

## Example

Write bot specification in `bot.yml` file:
```yaml
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
Create new Telegram bot and get its token: https://core.telegram.org/bots#6-botfather

Provide this token to docker image as `BOT_TOKEN` environment.

Run your bot:
```sh
docker run \
    -v $PWD/bot.yml:/w/config.yml \
    --env BOT_TOKEN="$BOT_TOKEN" \
    g4s8/openbots:latest
```

## Quick start

See [quick start](https://g4s8.github.io/openbots/quickstart/) guide to create Telegram bot in minutes.

## Documentation

The full documentation is available here: [g4s8.github.io/openbots/documentation](https://g4s8.github.io/openbots/documentation/).

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
 - [x] edit message
 - [ ] reply with images
 - [ ] handle image messages
 - [ ] delete messages
 - [ ] API:
   - [x] send message to particular user
 - [x] call webhook on update
 - [x] database storage
 
## Persistence

Bot can keep its state in two modes:
 - `memory` - store all data in memory
 - `database` - connect PostgreSQL database

For persistence configuration see [documentation](https://g4s8.github.io/openbots/persistence).

## Extending

This bot could be extended with custom handlers on Go.
