# Architecture

The bot uses a fully serverless architecture.

Telegram
↓ 
Webhook API Gateway
↓
AWS Lambda (Go)
↓
DynamoDB

---

## Why this architecture?

- No servers to manage
- Extremely low cost
- Scales automatically
- Secure webhook delivery
