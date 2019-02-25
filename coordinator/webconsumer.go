// -----------------------------------------------------------------------------
// Coordinator package used for defining queue listener and event aggregator.
// -----------------------------------------------------------------------------
package coordinator

import (
	"bytes"
	"encoding/gob"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"
	"log"

	"github.com/streadway/amqp"
)

type WebappConsumer struct {
	eventRaiser EventRaiser
	connection  *amqp.Connection
	channel     *amqp.Channel
	sources     []string
}

func NewWebappConsumer(eventRaiser EventRaiser) *WebappConsumer {
	consumer := WebappConsumer{
		eventRaiser: eventRaiser,
	}

	consumer.connection, consumer.channel = common.GetChannel(common.URL_GUEST)
	common.GetQueue(common.PERSISTENCE_QUEUE, consumer.channel, false)

	go consumer.ListenForDiscoveryRequests()

	consumer.eventRaiser.Subscribe(common.SENSOR_DISCOVER_EVENT, func(eventData Any) {
		consumer.Subscribe(eventData.(string))
	})

	consumer.channel.ExchangeDeclare(
		common.WEBAPP_SOURCE_EXCHANGE, common.FANOUT, false, false, false, false, nil)

	consumer.channel.ExchangeDeclare(
		common.WEBAPP_READINGS_EXCHANGE, common.FANOUT, false, false, false, false, nil)

	return &consumer
}

func (consumer *WebappConsumer) ListenForDiscoveryRequests() {
	queue := common.GetQueue(common.WEBAPP_DISCOVERY_QUEUE, consumer.channel, false)
	msgs, _ := consumer.channel.Consume(queue.Name, "", true, false, false, false, nil)

	log.Print("Web Consumer: ListenForDiscoveryRequests")

	for range msgs {
		for _, src := range consumer.sources {
			consumer.SendMessageSource(src)
		}
	}
}

func (consumer *WebappConsumer) SendMessageSource(src string) {
	log.Printf("Web Consumer: Sending message: %s", src)
	consumer.channel.Publish(
		common.WEBAPP_SOURCE_EXCHANGE, "", false, false, amqp.Publishing{Body: []byte(src)})
}

func (consumer *WebappConsumer) Subscribe(eventName string) {
	for _, v := range consumer.sources {
		if v == eventName {
			return
		}
	}

	consumer.sources = append(consumer.sources, eventName)
	consumer.SendMessageSource(eventName)

	toEvent := common.NewEvent(common.MESSAGE_RECEIVED_EVENT, eventName)
	consumer.eventRaiser.Subscribe(toEvent, func(eventData Any) {
		data := eventData.(dto.EventData)
		readout := dto.Readout{
			Name:      data.Name,
			Value:     data.Value,
			Timestamp: data.Timestamp,
		}

		buffer := new(bytes.Buffer)
		encoder := gob.NewEncoder(buffer)
		encoder.Encode(readout)

		message := amqp.Publishing{
			Body: buffer.Bytes(),
		}

		log.Printf("Web Consumer: Sending readout from: %s", data.Name)
		consumer.channel.Publish(common.WEBAPP_READINGS_EXCHANGE, "", false, false, message)
	})
}
