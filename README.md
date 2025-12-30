# my-links-bot

A small Go service that extracts URLs from incoming Telegram messages and stores them in DynamoDB for later reading/management.

**Features**
- Extracts URLs from Telegram updates and saves them as resources in DynamoDB.
- Minimal, dependency-light implementation using `github.com/go-telegram-bot-api/telegram-bot-api/v5` and AWS SDK v2.

**Project Structure**
- [bot/handler.go](bot/handler.go) : Parses Telegram update JSON, extracts URLs and forwards them to the DB layer.
- [db/dynamo.go](db/dynamo.go) : Initializes a DynamoDB client and provides `SaveResource` to persist items.
- [models/resouce.go](models/resouce.go) : `Resource` model used to represent stored links.

**Requirements**
- Go 1.24 (the module uses `go 1.24.4`).
- AWS credentials/config available to the runtime (for local testing this can be the AWS CLI-configured profile or environment variables).
- A DynamoDB table named `Resources` (see table schema below).

**DynamoDB table & item layout**
The code writes items with the following attributes (strings):
- `PK`: partition key
- `SK`: sort key
- `resource_id`, `type`, `url`, `status`, `notes`, `created_at`

Create a table with a string `PK` partition key and a string `SK` sort key. The code expects the table name to be `Resources` (constant in `db/dynamo.go`).
