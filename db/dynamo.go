package db

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var Client *dynamodb.Client

const TableName = "Resources"

func Init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	Client = dynamodb.NewFromConfig(cfg)
}

func SaveResource(userID int64, url, notes string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("USER#%d", userID),
		},
		"SK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("RES#%s#%s", now, id),
		},
		"resource_id": &types.AttributeValueMemberS{Value: id},
		"type":        &types.AttributeValueMemberS{Value: "article"},
		"url":         &types.AttributeValueMemberS{Value: url},
		"status":      &types.AttributeValueMemberS{Value: "to_read"},
		"notes":       &types.AttributeValueMemberS{Value: notes},
		"created_at":  &types.AttributeValueMemberS{Value: now},
	}

	_, err := Client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &[]string{TableName}[0],
		Item:      item,
	})

	return err

}
