package dlqslack

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/sirupsen/logrus"
)

// SendNotifyArguments is parameters for dlqslack
type SendNotifyArguments struct {
	IncomingURL string
	Event       interface{}
}

type deadLetterQueue struct {
	ErrorCode    string
	ErrorMessage string
	RequsetID    string
	Source       string
	Events       []string
}

// SendNotify parses event and send information to Slack
func SendNotify(args SendNotifyArguments) error {
	Logger.WithField("args", args).Debug("Start handler")

	var dlq []*deadLetterQueue
	var err error

	switch args.Event.(type) {
	case events.SQSEvent:
		ev := args.Event.(events.SQSEvent)
		dlq, err = parseSQSEvent(ev)
		if err != nil {
			return err
		}

	case events.SNSEvent:
		ev := args.Event.(events.SQSEvent)
		dlq, err = parseSQSEvent(ev)
		if err != nil {
			return err
		}

	default:
		Logger.WithField("event", args.Event).Error("Unsupported event format")
		return fmt.Errorf("Unsupported event foramt: %v", args.Event)
	}

	if err := notifyToSlack(args.IncomingURL, dlq); err != nil {
		return err
	}

	return nil
}

func init() {
	Logger.SetLevel(logrus.ErrorLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
