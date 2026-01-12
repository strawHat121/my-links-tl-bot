package db

import (
	"context"
	"fmt"
	"my-links-bot/models"
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

func SaveResource(r models.Resource) error {
	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	r.ResourceID = id
	r.CreatedAt = now
	r.SK = fmt.Sprintf("RES#%s#%s", now, id)

	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("USER#%d", r.UserID),
		},
		"SK": &types.AttributeValueMemberS{
			Value: r.SK,
		},
		"resource_id": &types.AttributeValueMemberS{Value: id},
		"type":        &types.AttributeValueMemberS{Value: r.Type},
		"title":       &types.AttributeValueMemberS{Value: r.Title},
		"url":         &types.AttributeValueMemberS{Value: r.URL},
		"status":      &types.AttributeValueMemberS{Value: r.Status},
		"notes":       &types.AttributeValueMemberS{Value: r.Notes},
		"created_at":  &types.AttributeValueMemberS{Value: now},
	}

	if len(r.Tags) > 0 {
		item["tags"] = &types.AttributeValueMemberSS{
			Value: r.Tags,
		}
	}

	_, err := Client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &[]string{TableName}[0],
		Item:      item,
	})

	return err
}

func ListResources(userID int64, limit int32) ([]models.Resource, error) {
	out, err := Client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              &[]string{TableName}[0],
		KeyConditionExpression: awsString("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("USER#%d", userID),
			},
			":sk": &types.AttributeValueMemberS{
				Value: "RES#",
			},
		},
		ScanIndexForward: awsBool(false),
		Limit:            &limit,
	})

	if err != nil {
		return nil, err
	}

	var res []models.Resource

	for _, item := range out.Items {
		r := models.Resource{
			UserID: userID,
			SK:     item["SK"].(*types.AttributeValueMemberS).Value,
			Type:   item["type"].(*types.AttributeValueMemberS).Value,
			Title:  item["title"].(*types.AttributeValueMemberS).Value,
			URL:    item["url"].(*types.AttributeValueMemberS).Value,
			Status: item["status"].(*types.AttributeValueMemberS).Value,
		}
		res = append(res, r)
	}

	return res, nil
}

func MarkDone(userID int64, sk string) error {
	_, err := Client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: &[]string{TableName}[0],
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("USER#%d", userID),
			},
			"SK": &types.AttributeValueMemberS{
				Value: sk,
			},
		},
		UpdateExpression: awsString("SET #s = :done"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status", // ✅ FIXED
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":done": &types.AttributeValueMemberS{
				Value: "completed",
			},
		},
	})

	return err
}

func awsString(s string) *string { return &s }
func awsBool(b bool) *bool       { return &b }
