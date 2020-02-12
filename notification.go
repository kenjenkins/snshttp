package snshttp

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// Notification events are sent for messages that are published to the SNS
// topic.
type Notification struct {
	BaseMessage
	Subject        string
	UnsubscribeURL string

	// MessageAttributes contain any attributes added to the message when
	// publishing it to SNS. This is most commonly used when transmitting binary
	// data (using raw message delivery).
	//
	// Note that message attributes are not part of the content signed by AWS.
	MessageAttributes map[string]MessageAttribute
}

// Unsubscribe will notify Amazon to remove this subscription from the SNS
// topic. It will make a request to the UnsubscribeURL and error if the
// request times out or the response does not indicate success.
func (e *Notification) Unsubscribe(ctx context.Context) error {
	req, err := http.NewRequest("GET", e.UnsubscribeURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Server is expected to return 200 OK but we can treat any 200 level code as
	// success.
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return fmt.Errorf("server returned error status=%d", resp.StatusCode)
	}

	return nil
}

func (e *Notification) SigningString() string {
	fields := []string{
		"Message", e.Message,
		"MessageId", e.MessageID,
	}

	if e.Subject != "" {
		fields = append(fields, "Subject", e.Subject)
	}

	fields = append(fields,
		"Timestamp", e.Timestamp,
		"TopicArn", e.TopicARN,
		"Type", e.Type,
		"")
	return strings.Join(fields, "\n")
}

type MessageAttribute struct {
	Type  string
	Value string
}

func (m MessageAttribute) StringValue() string {
	return m.Value
}

func (m MessageAttribute) BinaryValue() ([]byte, error) {
	return base64.StdEncoding.DecodeString(m.Value)
}
