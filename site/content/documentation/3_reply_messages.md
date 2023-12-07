---
title: "Replying with Messages: Formatting and Templating"
date: 2023-12-07T01:10:25+04:00
weight: 30
menuTitle: "Reply Messages"
---

The `reply.message` element allows you to reply to a user event with a text message,
providing options for text formatting and templating.

## Basic Usage

You can use the reply.message element in two ways: as a simple string or as a YAML mapping structure.
Both representations achieve the same result:

**String Representation:**
```yml
reply:
- message: "Hello"
```

**YAML Mapping Structure:**
```yml
reply:
- message:
    text: "Hello"
```

## Options

The `reply.message` element can have the following options:

 * **text (required):** The string for the reply text.
 * **parseMode (optional, default: none):** An enum string, one of "Markdown", "MarkdownV2", "HTML".
 Specifies the Telegram parse-mode for parsing entities in the message text. "Markdown" is a legacy mode,
 "HTML" processes some HTML formatting in the message text, and "MarkdownV2" is an extended Telegram markdown parsing mode.
 * **markup (optional):** Reply markup settings.
 * **template (optional, default: "default"):** An enum string, one of "default", "go", "no".
 Represents the template rendering engine. Default template style interpolates local fields from `${state.key}` expressions,
 "go" style uses the Go template engine, and "no" ignores templates and returns text as is.

**Example:**

```yml
reply:
  - message:
      text: "**Hello,** ${user.first_name}!"
      parseMode: "MarkdownV2"
      template: "default"
```

In this example, the bot replies with a formatted message using the MarkdownV2 parse mode and default template style.

See for details about parse-mode: https://core.telegram.org/bots/api#formatting-options

## Templating

**Default Interpolator:**

The default interpolator uses ${key} syntax to render values into the string. Possible keys include:

 * `state.<key>`: Get the state value for the key.
 * `secret.<key>`: Get the secret value for the key.
 * `message.id`: Current message ID (the message that triggered this event).
 * `message.text`: Current message text.
 * `message.from.id`: Telegram ID of the message sender.
 * `chat.id`: Current Telegram chat ID.
 * `chat.type`: Type of chat ("private", "group", "supergroup", or "channel").


**Go Interpolator:**

The Go interpolator uses the GoLang template format. Possible template variables include:

 * `Update`: Current Telegram update object.
 * `State`: Key-value pairs of the user's state.
 * `Secrets`: Key-value pairs of bot secrets.
 * `Data`: JSON object loaded by the data-loader.


**No Template Engine:**

The "no" template engine renders text as is without any interpolation.

### Examples

**Default interpolator:**

```yml
reply:
  - message:
      text: "Hello, ${user.first_name}! You are ${state.age} years old."
      template: "default"
```

**Go template engine:**

```yml
reply:
  - message:
      text: "Hello, {{.Update.Message.From.FirstName}}! You are {{.State.age}} years old."
      template: "go"
```

**No Template Engine:**

```yml
reply:
  - message:
      text: "This is plain text with no template. {{Render}} ${this text as is}."
      template: "no"
```

Explore further to understand the full potential of text formatting and templating options in your bot replies.

### Possible interpolator variables

- `state.<key>` - get state value for key;
- `secret.<key>` - get secret value for key;
- `message.id` - get current message id;
- `message.text` for current message text;
- `message.from.id` - Telegram id of message sender;
- `chat.id` - current Telegram chat id;
- `chat.type` - Type of chat, can be either “private”, “group”, “supergroup” or “channel”;
- `chat.title` - Title for supergroups, channels and group chats;
- `chat.first_name` - FirstName of the other party in a private chat;
- `chat.last_name` - LastName of the other party in a private chat;
- `chat.username` - UserName for private chats, supergroups and channels if available;
- `user.id` - ID of user who sent this message;
- `user.is_bot` - true if sender is bot;
- `user.first_name` - FirstName user's or bot's last name;
- `user.last_name` - LastName user's or bot's last name;
- `user.username` - UserName user's or bot's username;
- `user.language_code` - LanguageCode IETF language tag of the user's language;

### Possible Go template engine variables

- `Update` - current telegram update object as https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5#Update
- `State` - key value pairs of user's state
- `Secrets` - key value pairs of bot secrets
- `Data` - json object loaded by data-loader
