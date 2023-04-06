package events

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	. "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublisherPublish(t *testing.T) {
	message := &WebHookEventMessage{
		Id:        uuid.New(),
		EventType: "test.event",
		Body:      make(map[string]any),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var received = &WebHookEventMessage{}
		json.Unmarshal(body, &received)

		// assert that the message that was published was received by the server that was registered to the publisher
		assert.Equal(t, message.Id, received.Id)
	}))
	defer ts.Close()

	publisher := getPublisher(ts.URL)

	fmt.Printf("URL: %s/\n", ts.URL)

	err := publisher.Publish(message)
	require.NoError(t, err)
}

func getPublisher(url string) WebHookPublisher {
	provider := newFakeSubscriptionProvider(url)
	publisher := NewWebHookPublisher(provider)
	return publisher
}

// FAKE
type fakeSubscriptionsProvider struct {
	subscriptions []*data.EventSubscription
}

func newFakeSubscriptionProvider(url string) SubscriptionsProvider {
	provider := &fakeSubscriptionsProvider{
		subscriptions: []*data.EventSubscription{
			{Callback: url, Name: "test-subscription-1", ApiKey: "testapikey", BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{
				ID: uuid.New(),
			}},
			{Callback: url, Name: "test-subscription-2", ApiKey: "testapikey", BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{
				ID: uuid.New(),
			}},
			{Callback: url, Name: "test-subscription-3", ApiKey: "testapikey", BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{
				ID: uuid.New(),
			}},
		},
	}
	return provider
}

func (p *fakeSubscriptionsProvider) GetSubscriptions() ([]*data.EventSubscription, error) {
	return p.subscriptions, nil
}
