# Deployment

## Requirements

- AWS account
- Telegram bot token
- Go 1.22+

---

## Environment variables (Lambda)

TELEGRAM_BOT_TOKEN=your-token

---

## Infrastructure

- API Gateway (POST /webhook)
- Lambda (custom Go runtime)
- DynamoDB table: `Resources`

---

## Webhook setup

https://api.telegram.org/bot<TOKEN>/setWebhook