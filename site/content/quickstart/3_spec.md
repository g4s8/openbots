---
title: "Bot spec"
date: 2023-12-06T18:46:12+04:00
weight: 4
---

Now let's write a bot specification.
The bot specification is a [YAML](https://yaml.org/) document
which describe bot handlers.


Before you start writing your bot specification, understand the deployment differences based on your setup.

 * **Self-Hosted Version:** If you're using the self-hosted version, 
 you'll need to create a YAML file for your bot specification.
 Follow the instructions below to deploy your bot.
 * **Cloud Version:** For the cloud version, simply open the bot page with the spec editor on the website.
 Changes made here will be automatically applied to your cloud-hosted bot.

## Writing Your First Bot Specification

Open spec file or bot page with spec editor. And write your first bot specification:

```yaml
bot:
  handlers:
  - on:
      message:
        command: start
    reply:
    - message:
        text: "Hello! Welcome to your new bot."
  - on:
      message: "How are you?"
    reply:
      - message:
          text: "I'm doing well, thank you!"
```

In this example, we have a simple bot with two handlers.
The first handler triggers when the user sends the command `/start`,
and the bot responds with a welcome message.
The second handler triggers when the user sends the message "How are you?"
and the bot responds with a positive message.

## Adding More Complexity

Let's enhance the bot specification by adding a handler that responds to a button click:

```yaml
bot:
  handlers:
  - on:
      message:
        command: start
    reply:
    - message:
        text: "Hello, how are you?"
        markup:
          inlineKeyboard:
          -
            - text: "Good"
              callback: callback-good
            - text: "Bad"
              callback: callback-bad
  - on:
      callback:
        data: "callback-good"
    reply:
    - callback:
        text: "Got it"
      message:
        text: "I'm glad to hear that you're doing well!"
  - on:
      callback:
        data: "callback-bad"
    reply:
    - callback:
        text: "Got it"
      message:
        text: "I'm sorry to hear that. Is there anything I can do to help?"
```

In this updated example:

 1. When the user sends the `/start` command, the bot responds with a question:
 "Hello, how are you?" along with two buttons for "Good" and "Bad."
 2. If the user clicks the "Good" button (`callback-good`), the bot responds with a positive message.
 3. If the user clicks the "Bad" button (`callback-bad`), the bot responds with a supportive message.

Feel free to test and modify this template based on your preferences and the desired bot behavior.

## Deploying Your Bot Specification

### Self-Hosted Version

If you're using the self-hosted version, you'll need to provide your bot specification file
when starting the bot service. Use the following CLI options (don't forget to specify Telegram bot token):

 * Using CLI option:

```bash
BOT_TOKEN="110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" ./openbots -config my_bot_spec.yml
```

 * Using Docker option:

```bash
docker run --rm --name openbots \
  --env BOT_TOKEN="110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" \
  -v $PWD/my_bot_spec.yml:/w/config.yml \
  g4s8/openbots:latest
```

### Cloud Version

For the cloud version, simply write your YAML spec on the website at the bot page in the spec editor.
The changes will be automatically applied to your cloud-hosted bot after deploying with "Save" button.

## Testing Your Bot Specification

To test your bot just start bot in Telegram app with "Start" button.

Congratulations! You've created a simple bot specification.
Feel free to explore more features and customize your bot according to your preferences.

Explore more features and customization options in the [full documentation](/openbots/documentation).
