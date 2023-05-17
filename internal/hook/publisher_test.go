package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublisherPublish(t *testing.T) {
	message := &sdk.EventHookMessage{
		Id:   uuid.New(),
		Type: "test.event",
		Data: make(map[string]any),
	}

	done := make(chan struct{})
	defer close(done)

	// should be 3 because of the fake provider
	receiveCount := 0
	mutex := sync.Mutex{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		receiveCount++
		defer mutex.Unlock()

		body, _ := io.ReadAll(r.Body)
		var received = &sdk.EventHookMessage{}
		json.Unmarshal(body, &received)

		log.Printf("Received message: %v", string(body))
		// assert that the message that was published was received by the server that was registered to the publisher
		assert.Equal(t, message.Id, received.Id)

		if receiveCount == 3 {
			log.Print("All messages received")
			done <- struct{}{}
		}
	}))
	defer ts.Close()

	go func() {
		publisher := getPublisher(ts.URL)

		fmt.Printf("URL: %s/\n", ts.URL)

		err := publisher.Publish(message)
		require.NoError(t, err)
	}()
	<-done
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
