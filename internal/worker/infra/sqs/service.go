package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clodoaldomarques/balances-api/config"
	"github.com/clodoaldomarques/balances-api/internal/shared/domain/events"
	"github.com/clodoaldomarques/core-sdk/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler struct {
	ctx         context.Context
	sqsClient   *sqs.Client
	queueUrl    string
	maxMessages int32
	waitTime    int32
}

func New(ctx context.Context, maxMessages, waitTime int32) *Handler {
	return &Handler{
		ctx:       ctx,
		sqsClient: NewSQSClient(ctx),
		queueUrl:  config.New().BalancesSQSQueue,
	}
}

func (c Handler) Retrieve(ctx context.Context) (map[string]events.Event, error) {
	result, err := c.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueUrl),
		MaxNumberOfMessages: c.maxMessages,
		WaitTimeSeconds:     c.waitTime,
	})
	if err != nil {
		logger.Error(ctx, "Couldn't get messages from queue.", logger.Fields{
			"queue": c.queueUrl,
			"error": err.Error(),
		})
		return map[string]events.Event{}, err
	}

	return messageToEvent(result.Messages)
}

func (c Handler) DeleteMessages(receiptHandles []string) error {
	entries := make([]types.DeleteMessageBatchRequestEntry, len(receiptHandles))
	for msgIndex, rh := range receiptHandles {
		entries[msgIndex].Id = aws.String(fmt.Sprintf("%v", msgIndex))
		entries[msgIndex].ReceiptHandle = &rh
	}
	_, err := c.sqsClient.DeleteMessageBatch(c.ctx, &sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(c.queueUrl),
	})
	if err != nil {
		logger.Error(c.ctx, "Couldn't delete messages from queue.", logger.Fields{
			"queue": c.queueUrl,
			"error": err.Error(),
		})
	}
	return err
}

func messageToEvent(messages []types.Message) (map[string]events.Event, error) {
	evts := map[string]events.Event{}
	for _, m := range messages {
		var objectMap map[string]string
		err := json.Unmarshal([]byte(*m.Body), &objectMap)
		if err != nil {
			return map[string]events.Event{}, err
		}

		var evt events.Event
		err = json.Unmarshal([]byte(objectMap["Message"]), &evt)
		if err != nil {
			return map[string]events.Event{}, err
		}
		evts[*m.ReceiptHandle] = evt
	}
	return evts, nil
}
