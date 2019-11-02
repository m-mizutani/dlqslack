package dlqslack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type secretValues struct {
	GithubToken string `json:"github_token"`
	SlackURL    string `json:"slack_url"`
}

type arguments struct {
	SecretArn        string
	GithubRepository string
	GithubEndpoint   string
	TeamMembers      []string
	IgnoreLabels     []string
}

type slackRequest struct {
	Text        string            `json:"text"`
	Attachments []slackAttachment `json:"attachments"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"bool"`
}

type slackAttachment struct {
	Color     string       `json:"color"`
	Title     string       `json:"title"`
	TitleLink string       `json:"title_link"`
	Text      string       `json:"text"`
	Fields    []slackField `json:"fields"`
	MrkdwnIn  []string     `json:"mrkdwn_in"`
}

var httpPost = http.Post

func notifyToSlack(url string, dlqSet []*deadLetterQueue) error {
	var attachments []slackAttachment

	for _, dlq := range dlqSet {
		attach := slackAttachment{
			MrkdwnIn: []string{"fields"},
			Color:    "#DEE735",
			Fields: []slackField{
				{
					Title: "ErrorCode",
					Value: dlq.ErrorCode,
					Short: true,
				},
				{
					Title: "RequestID",
					Value: dlq.RequsetID,
					Short: true,
				},
				{
					Title: "ErrorMessage",
					Value: dlq.ErrorMessage,
				},
				{
					Title: "Source",
					Value: dlq.Source,
				},
				{
					Title: "Events",
					Value: "```\n" + strings.Join(dlq.Events, "\n\n") + "\n```",
				},
			},
		}

		attachments = append(attachments, attach)
	}

	req := slackRequest{
		Attachments: attachments,
	}
	reqbuf, err := json.Marshal(req)
	if err != nil {
		return errors.Wrapf(err, "Fail to marshal slack requset: %v", req)
	}

	resp, err := httpPost(url, "application/json", bytes.NewReader(reqbuf))
	if err != nil {
		return errors.Wrap(err, "Fail to post message to Slack")
	}

	Logger.WithField("response", resp).Debug("Sent request to Slack")

	return nil
}
