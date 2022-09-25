---
title: "Start bot"
date: 2022-09-25T21:04:36+03:00
weight: 4
---

There are two simple ways to start the bot: using `openbots` binary
or from Docker container.

## Start the bot with openbots binary

If you downloaded and extracted `openbots` binary in current directory
and created bot config in the same directory too, then run:
```bash
env BOT_TOKEN="110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" ./openbots -config bot.yml
```
Use your bot token from previous step instead of placeholder for `BOT_TOKEN` environent
variable.

That's all - you can find your bot in telegram and click "Start" button,
the bot will respond you with "Hello".

## Using Docker container

Run this `docker` command:
```bash
docker run --rm --name openbots \
  --env BOT_TOKEN="110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" \
  -v $PWD/bot.yml:/w/config.yml \
  g4s8/openbots:latest
```

Don't forget to provide your actual bot token for `--env BOT_TOKEN` argument instead
of placeholder.

The bot is up and running.
