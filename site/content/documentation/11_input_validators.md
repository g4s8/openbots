---
title: "User Input Validation"
date: 2023-12-07T22:34:26+04:00
weight: 110
menuTitle: "Input Validation"
---

Ensure that user inputs meet specific criteria by employing user input validation in your bot.
Handlers can include a validate object to check user inputs against predefined conditions.
If the validation fails, the bot sends an error message and prevents the execution of the subsequent handler actions.

## Validate Object Elements

The validate object includes the following elements:

 * `error_message` (required): A string representing the error message to be sent if the validation fails.
 * `checks` (required): An array of check names, specifying the validation criteria.

**Supported Validation Checks:**

 * `not_empty`: Validates that the input is not empty.
 * `is_int`: Validates that the input is an integer.
 * `is_float`: Validates that the input is a floating-point number.
 * `is_bool`: Validates that the input is a boolean value.

## Example

Consider the following example where the bot prompts the user to enter a number for the variable `X`.
The `validate` object checks that the input is a non-empty integer:

```yml
bot:
  handlers:
    - on:
        message:
          command: start
      reply:
        - message:
            text: "Enter X number"
        context:
          set: read-x
    - on:
        context: read-x
      validate:
        error_message: x should be a number
        checks: ["not_empty", "is_int"]
      state:
        set:
          x: "${message.text}"
      reply:
        - message:
            text: "Saved X variable"
        context:
          delete: read-x
```

In this example, if the user's input does not meet the validation criteria,
the bot sends the specified error message, and the state and reply actions are not executed.
