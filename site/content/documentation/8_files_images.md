---
title: "Files, Images, and Documents"
date: 2023-12-07T21:18:42+04:00
weight: 80
menuTitle: "File, Images"
---

Handling files, images, and documents in your bot allows you to share multimedia content with users.

## Asset Management

There are slight differences between the cloud and self-hosted versions in managing assets.

### Cloud Version

For cloud users, assets can be managed on the main bot page.
Upload assets and use the provided names as keys for files or images in your bot configuration.

### Self-Hosted Version

Self-hosted users need to manually configure the assets provider.
Key formats may vary depending on the chosen provider.
See advanced configuration documentation for self-hosted setup.

## Replying with Images

To reply with an image, use the `image` element in your bot configuration:

```yml
handlers:
  - on:
      message:
        command: start
    reply:
      - image:
          name: Test image
          key: test-image.png
```

In this example, name represents the image name displayed in the chat, and key is the asset key.

## Replying as a Document

To reply with a document (file), use the `document` element in your bot configuration:

```yml
handlers:
  - on:
      message:
        command: start
    reply:
      - document:
          name: Test document
          key: test-document.txt
```

In this example, `name` represents the document name displayed in the chat,
and `key` is the asset key.

**Note:** when replying with an image or document, the asset key should point to the corresponding file.

Explore the possibilities of multimedia interactions to enhance user engagement with your bot.
