package main

import (
	"context"
	"my-links-bot/bot"
	"my-links-bot/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {

	bot.Handle([]byte(req.Body))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	db.Init()
	lambda.Start(handler)
}
