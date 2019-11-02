package dlqslack

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
)

type awsEvent struct {
	Records []map[string]interface{} `json:"Records"`
}

func jsonToEvent(jdata string) ([]string, error) {
	var input []string

	var event awsEvent
	if err := json.Unmarshal([]byte(jdata), &event); err != nil {
		return nil, errors.Wrapf(err, "Fail to unmarshal base event: %v", jdata)
	}

	for _, record := range event.Records {
		var msg []byte
		switch record["EventSource"] {
		case "aws:sns":
			evData, err := json.Marshal(record)
			if err != nil {
				return nil, errors.Wrapf(err, "Fail to marshal inner SNS record: %v", record)
			}
			var snsRecord events.SNSEventRecord
			if err := json.Unmarshal(evData, &snsRecord); err != nil {
				return nil, errors.Wrapf(err, "Fail to unmarshal inner SNS record: %v", string(evData))
			}

			msg = []byte(snsRecord.SNS.Message)

		case "aws:sqs":
			evData, err := json.Marshal(record)
			if err != nil {
				return nil, errors.Wrapf(err, "Fail to marshal inner SQS record: %v", record)
			}
			var sqsRecord events.SQSMessage
			if err := json.Unmarshal(evData, &sqsRecord); err != nil {
				return nil, errors.Wrapf(err, "Fail to unmarshal inner SQS record: %v", string(evData))
			}

			msg = []byte(sqsRecord.Body)

		default:
			evData, err := json.Marshal(record)
			if err != nil {
				return nil, errors.Wrapf(err, "Fail to marshal inner record: %v", record)
			}

			msg = evData
		}

		var msgJSON bytes.Buffer
		if err := json.Indent(&msgJSON, msg, "", "  "); err != nil {
			return nil, errors.Wrapf(err, "Fail to unmarshal original event: %v ", string(msg))
		}

		input = append(input, msgJSON.String())
	}

	return input, nil
}

func parseSQSEvent(ev events.SQSEvent) ([]*deadLetterQueue, error) {
	var dlqList []*deadLetterQueue
	for _, record := range ev.Records {
		var dlq deadLetterQueue
		for attrKey, attrValue := range record.MessageAttributes {
			switch attrKey {
			case "ErrorCode":
				dlq.ErrorCode = aws.StringValue(attrValue.StringValue)
			case "ErrorMessage":
				dlq.ErrorMessage = aws.StringValue(attrValue.StringValue)
			case "RequestID":
				dlq.RequsetID = aws.StringValue(attrValue.StringValue)
			}
		}

		originEvents, err := jsonToEvent(record.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Fail to parse SQS body")
		}
		dlq.Events = originEvents
		dlq.Source = record.EventSourceARN

		dlqList = append(dlqList, &dlq)
	}

	return dlqList, nil
}

func parseSNSEvent(ev events.SQSEvent) ([]*deadLetterQueue, error) {
	return nil, nil
}
