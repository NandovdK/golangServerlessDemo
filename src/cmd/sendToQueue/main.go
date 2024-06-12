package main

import (
	"context"
	"encoding/json"
	"errors"

	"ac-it/src/constants"
	"ac-it/src/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
)

func Handler(ctx context.Context, sqsEvent events.SQSEvent, sqsClient *services.SqsClient, log *logrus.Logger) error {
	message := sqsEvent.Records[0].Body

	var messageBody map[string]interface{}
	if err := json.Unmarshal([]byte(message), &messageBody); err != nil {
		log.Error("Error unmarshalling message: ", err)
		return err
	}

	if messageBody["prompt"] == nil {
		err := errors.New("prompt not found in message")
		log.Error(err.Error())
		return err
	}

	err := sqsClient.PostToQueue(messageBody, "answerQueue")
	if err != nil {
		log.Error("Error posting to queue: ", err)
		return err
	}
	log.Info("Message processed successfully")
	return nil
}

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	}))

	sqsClient := services.NewSqsClient(sess)
	log := services.GetLogger()

	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		return Handler(ctx, sqsEvent, sqsClient, log)
	})
}
