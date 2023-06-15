---
title: "State"
date: 2022-09-25T19:37:45+03:00
weight: 6
---

Each user may have a local state as a key-value map.

It could be changed with `state` handler. But first let's see
how you can access it and use in your bot.

Text message may include state references and bot interpolates
these references to actual values of the state. E.g. user have
key `name` with value `John` in state, then if you want to greet 
user you can create message reply which uses state reference:
```yaml
reply:
  - message:
      text: "Hello ${state.name}"
```
And this message will be interpolated to `Hello John`.

Now let's change the state with `state` handler's object:
```yaml
on:
  message: John
  context: set-name
reply:
  - message:
      text: "Welcome!"
state:
  set:
    name: John
```
Now `name` key in state was changed to `John` value.

To clear the state, use `delete`:
```yaml
state:
  delete: name
```

User's state can also be set to dynamic message values, e.g. message text:
```yaml
state:
  set:
    name: "${message.text}"
```
