---
title: "User State Management"
date: 2023-12-07T19:20:55+04:00
weight: 40
menuTitle: "User State"
---

In the bot configuration, user states provide a way to maintain persistent information about each user.
Handlers can read, use, and modify the user's state, allowing for dynamic and personalized interactions.

## Modify User State

A basic example of modifying the user's state:

```yml
bot:
  handlers:
    - on:
        message: test
      state:
        # apply state changes
```

The `state` element in the handler can have the following sub-elements:

### Set

Assign values to state variables:

```yml
state:
  set:
    foo: "qwe"
    bar: "asd"
    baz: "zxc"
```

After this handler, the user's state will have three variables with values:
`{"foo": "qwe", "bar": "asd", "baz": "zxc"}`.

### Delete

Delete state variable(s):

```yml
state:
  delete: "bar"
```

After applying this handler to the previous state, the "bar" variable will be removed:
`{"foo": "qwe", "baz": "zxc"}`.

Delete value may have string or array of strings, all examples below are valid:

**Delete one variable:**

```yml
state:
  delete: key1
```

**Delete one variable via array syntax:**

```yml
state:
  delete: ["key1"]
```

**Delete multiple variables:**

```yml
state:
  delete: ["key1", "key2"]
```

### Ops

Apply more complex operations on the state. Operations consist of:

 * `kind` (required): Enum string, one of (set, delete, add, sub, mul, div).
 * `key` (required): String, the state variable name/key to apply the operation.
 * `value` (required for some ops): String, the argument for the operation.

**Example:**

```yml
state:
  ops:
    - kind: add
      key: "x"
      value: "4"
```

This operation is equivalent to `x += 4`.

### Using Interpolation in State Values

State values can use default interpolation strings. For example:

```yml
state:
  ops:
    - kind: add
      key: "x"
      value: "${state.y}"
```

This operation adds the value of the "y" variable to the existing "x" variable.

## Reading and Using State Variables

Most values in handlers can access state via interpolation. For instance, to reply using state:

```yml
- on:
    message:
      command: val
  reply:
    - message: "value: ${state.val}"
```

Or using Go templates:

```yml
- on:
    nessage: test
  reply:
    - message:
        template: go
        text: |
          {{$x := .State.x}}
          {{$y := .Update.Message.Text}}
          {{$x}} + {{$y}} = {{sum $x $y}}
```

In this example, the handler calculates the sum of the state variable "x" and the user's text and prints the result.

State variables can be utilized in webhooks, payment handlers, various reply handlers (text, images, documents, etc.),
edit replies, keyboards, and more.

Explore the power of user state management to create dynamic and personalized bot interactions.
