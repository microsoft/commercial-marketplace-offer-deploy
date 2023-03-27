package events

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

const contentTypeJson string = "application/json"

type MessageSender interface {
	Send(ctx context.Context, data any) error
}

type HttpSender struct {
	url    string
	apiKey string
}

func NewMessageSender(url string, apiKey string) MessageSender {
	sender := &HttpSender{url: url, apiKey: apiKey}
	return sender
}

func (sender *HttpSender) Send(ctx context.Context, data any) error {
	request, err := sender.createRequest(ctx, data)

	if err != nil {
		return err
	}

	err = retry.Do(func() error {
		client := http.Client{
			Timeout: 30 * time.Second,
		}

		response, err := client.Do(request)

		if err != nil {
			return err
		}

		defer response.Body.Close()
		var body []byte
		body, err = ioutil.ReadAll(response.Body)

		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("request failed with status [%d] '%s'", response.StatusCode, string(body))
		}

		return nil
	})

	return err
}

func (sender *HttpSender) createRequest(ctx context.Context, data any) (*http.Request, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(jsonData)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, sender.url, buffer)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentTypeJson)
	request.Header.Set("Authorization", sender.getAuthorizationHeaderValue())

	return request, nil
}

func (sender *HttpSender) getAuthorizationHeaderValue() string {
	encodedApiKey := base64.StdEncoding.EncodeToString([]byte(sender.apiKey))
	return "ApiKey " + encodedApiKey
}
