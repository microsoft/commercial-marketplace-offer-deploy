package notification

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestReadNotificationPump(t *testing.T) {
	

}

func TestStageNotificationPump_receiver_function_called(t *testing.T) {
	mockDB := getDatabase()

	var receiverCalled bool
	receiver := func(notification *model.StageNotification) error {
		receiverCalled = true
		return nil
	}

	// Create a test record with Done set to false
	record := &model.StageNotification{
		Done: false,
		OperationId: uuid.New(),
		CorrelationId: uuid.New(),
		ResourceGroupName: "test",
	}
	mockDB.Save(record)

	pump := NewStageNotificationPump(mockDB, 10 * time.Second, receiver)
	
	pump.Start()

	time.Sleep(15 * time.Second)

	// Assert that the receiver function was called
	assert.True(t, receiverCalled)
}

func TestStageNotificationPump_read_unsent(t *testing.T) {
	mockDB := getDatabase()

	// Create a test instance of StageNotificationPump
	pump := &StageNotificationPump{
		db: mockDB,
	}

	// Create a test record with Done set to false
	expectedRecord := &model.StageNotification{
		Done: false,
		// Set other required fields if any
	}
	mockDB.Create(expectedRecord)

	// Call the read function
	record, found := pump.read()

	// Assert that the record is found
	assert.True(t, found)
	assert.NotNil(t, record)

	// Assert that the returned record matches the expected record
	assert.Equal(t, expectedRecord.ID, record.ID)
	assert.Equal(t, expectedRecord.Done, record.Done)

	// Assert that Done field is false
	assert.False(t, record.Done)
}

func TestStageNotificationPump_read_sent(t *testing.T) {
	mockDB := getDatabase()

	// Create a test instance of StageNotificationPump
	pump := &StageNotificationPump{
		db: mockDB,
	}

	// Create a test record with Done set to true
	record := &model.StageNotification{
		Done: true,
		// Set other required fields if any
	}
	mockDB.Create(record)

	// Call the read function
	result, found := pump.read()

	// Assert that the record is not found
	assert.False(t, found)
	assert.Nil(t, result)
}


func getDatabase() *gorm.DB {
	database := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	return database
}