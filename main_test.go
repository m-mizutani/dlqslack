package dlqslack_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"

	"github.com/m-mizutani/dlqslack"
)

func TestActualAction(t *testing.T) {
	url := os.Getenv("DLQ_SLACK_URL")
	fpath := os.Getenv("DLQ_SLACK_EVENT_FILE")

	if url == "" || fpath == "" {
		t.Skip("DLQ_SLACK_URL and DLQ_SLACK_EVENT_FILE are required")
	}

	raw, err := ioutil.ReadFile(fpath)
	require.NoError(t, err)
	var sqsEvent events.SQSEvent
	err = json.Unmarshal(raw, &sqsEvent)
	require.NoError(t, err)

	args := dlqslack.SendNotifyArguments{
		IncomingURL: url,
		Event:       sqsEvent,
	}

	err = dlqslack.SendNotify(args)
	require.NoError(t, err)
}
