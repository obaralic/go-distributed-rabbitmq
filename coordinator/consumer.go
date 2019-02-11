// -----------------------------------------------------------------------------
// Coordinator package used for defining queue listener and event aggregator.
// -----------------------------------------------------------------------------
package coordinator

import (
	"bytes"
	"encoding/gob"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"
	"time"

	"github.com/streadway/amqp"
)

const MAX_RATE = 5 * time.Second

// -----------------------------------------------------------------------------
// StorageConsumer - Struct that represents event storaging consumer.
// -----------------------------------------------------------------------------
type StorageConsumer struct {
	eventRaiser EventRaiser
	connection  *amqp.Connection
	channel     *amqp.Channel
	queue       *amqp.Queue
	sources     []string
}

// -----------------------------------------------------------------------------
// NewConsumer - Creates new event consumer.
// -----------------------------------------------------------------------------
func NewStorageConsumer(eventRaiser EventRaiser) (consumer *StorageConsumer) {
	consumer = &StorageConsumer{eventRaiser: eventRaiser}

	// Consider reusing connection between listener and consumer.
	consumer.connection, consumer.channel = common.GetChannel(common.URL_GUEST)
	consumer.queue = common.GetQueue(common.PERSISTENCE_QUEUE, consumer.channel, false)

	// Start listening for new sensors
	consumer.eventRaiser.Subscribe(common.SENSOR_DISCOVER_EVENT, func(eventData Any) {
		// Explicit type conversion from interface{} aka Any to string
		sensorName := eventData.(string)
		consumer.Subscribe(sensorName)
	})

	return
}

// -----------------------------------------------------------------------------
// Subscribe - Subscribes consumer to the given event.
// -----------------------------------------------------------------------------
func (consumer *StorageConsumer) Subscribe(eventName string) {
	// Check if we are already listening for the given event.
	for _, observed := range consumer.sources {
		if observed == eventName {
			return
		}
	}

	// If not register it to event raiser.
	toEvent := common.NewEvent(common.MESSAGE_RECEIVED_EVENT, eventName)
	consumer.eventRaiser.Subscribe(toEvent, func() func(Any) {

		prevTime := time.Unix(0, 0)
		buffer := new(bytes.Buffer)

		// returns closure function
		return func(triggerData Any) {

			// Explicit type conversion from interface{} aka Any to EventData
			eventData := triggerData.(dto.EventData)

			if time.Since(prevTime) > MAX_RATE {
				prevTime = time.Now()
				buffer.Reset()
				readout := dto.Convert(eventData)
				encoder := gob.NewEncoder(buffer)
				encoder.Encode(readout)

				// Send a message to the storage handling endpoint.
				common.Send(buffer.Bytes(), consumer.queue, consumer.channel)
			}
		}
	}())

}
