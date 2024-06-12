package main

import (
	"ac-it/src/constants"
	"ac-it/src/services"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	}))

	sqsClient := services.NewSqsClient(sess)

	message := sqsEvent.Records[0].Body

	var messageBody map[string]interface{}
	if err := json.Unmarshal([]byte(message), &messageBody); err != nil {
		fmt.Println("Error unmarshalling message: ", err)
	}

	if messageBody["prompt"] == nil {
		return errors.New("prompt not found in message")
	}

	err := sqsClient.PostToQueue(messageBody, "answerQueue")
	if err != nil {
		fmt.Println("Error posting to queue: ", err)
		return err
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
