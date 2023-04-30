package hook

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
	message := &EventHookMessage{
		Id:        uuid.New(),
		EventType: "test.event",
		Body:      make(map[string]any),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var received = &EventHookMessage{}
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

func getPublisher(url string) Publisher {
	provider := newFakProvider(url)
	publisher := NewEventHookPublisher(provider)
	return publisher
}

//region fakes

type fakeProvider struct {
	subscriptions []*data.EventHook
}

func newFakProvider(url string) EventHooksProvider {
	provider := &fakeProvider{
		subscriptions: []*data.EventHook{
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

func (p *fakeProvider) Get() ([]*data.EventHook, error) {
	return p.subscriptions, nil
}

//endregion fakes
