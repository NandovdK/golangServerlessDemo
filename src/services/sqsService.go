package services

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SqsClient struct {
	client *sqs.SQS
}

func NewSqsClient(sess *session.Session) *SqsClient {
	return &SqsClient{
		client: sqs.New(sess),
	}
}

func (s *SqsClient) GetQueueUrl(queueName string) (string, error) {
	result, err := s.client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		fmt.Println("Error getting queue URL: ", err)
		return "", err
	}
	return *result.QueueUrl, nil
}

func (s *SqsClient) PostToQueue(message interface{}, queueName string) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message: ", err)
		return err
	}

	queueUrl, err := s.GetQueueUrl(queueName)
	if err != nil {
		fmt.Println("Error getting queue URL: ", err)
		return err
	}

	s.client.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBytes)),
		QueueUrl:    aws.String(queueUrl),
	})
	return nil
}
