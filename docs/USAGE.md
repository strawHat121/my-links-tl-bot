# Usage

This document explains how to use the bot from Telegram.

---

## Saving a link (implicit)

Send any message containing a URL:

Check this out https://example.com #golang

The bot will:
- Extract the title
- Save the link
- Store tags
- Reply with `✅ Saved`

---

## Saving explicitly with `/save`

/save video https://youtube.com/watch?v=abc123 #art


Supported types:
- article
- video
- book
