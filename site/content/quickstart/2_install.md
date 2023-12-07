---
title: "Install"
date: 2022-09-25T19:37:45+03:00
weight: 3
---

Before you start creating a new bot,
you need to decide whether to use [the self-hosted](#self-hosted-version) version or
[the Cloud](#start-bot-at-chatbotcloudorg) version
available at [chatbotcloud.org](https://chatbotcloud.org).
This documentation provides instructions for both versions.


## Self hosted version

To run a self hosted version you can download a binary or run Docker image.

### Running binary

You can download the binary release from
[github.com/g4s8/openbots/releases](https://github.com/g4s8/openbots/releases):
```bash
# example for Linux x86_64
wget https://github.com/g4s8/openbots/releases/download/<version>/openbots_<version>_Linux_x86_64.tar.gz
tar -xvzf openbots_0.0.4_Linux_x86_64.tar.gz 
rm openbots_0.0.4_Linux_x86_64.tar.gz
```

Thenn you need to specify Telegram bot token (you can get it from [@BotFather](https://t.me/botfather)),
see [Telegram docs](https://core.telegram.org/bots/api), and spec path:
```bash
BOT_TOKEN="110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" ./openbots -config bot.yml
```

### Running Docker image

```bash
docker run --rm --name openbots \
  --env BOT_TOKEN="110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" \
  -v $PWD/bot.yml:/w/config.yml \
  g4s8/openbots:latest
```

Don’t forget to provide your actual bot token via `--env BOT_TOKEN` argument instead of placeholder,
and actual bot spec file path.


## Start bot at chatbotcloud.org

For the cloud version follow this steps:
 1. Go to https://chatbotcloud.org/ website.
 2. Click "SignUp" button if you don't have account, follow instructions to create a new account.
 Or if you have an account just click "SignIn" button.
 3. Navigate to the "Dashboard" page at https://chatbotcloud.org/dashboard
 4. Click "Add bot" button at the bottom of the "Dashboard" page, enter a new bot name and click "Create".
 5. Click "Save" button on bot spec text editor to save "Hello" bot examople.
 6. Get bot token from [@BotFather](https://t.me/botfather) official Telegram bot.
 7. Go to bot "Settings" tab, click "Set token" button, and enter bot token from "BotFather" into
 input, then click "Save".

Now your bot should be deployed to one of the cloud servers and running
to serve your Telegram bot requests. Just open your new created bot from "BotFather" and click "Start" button
in Telegram app.

