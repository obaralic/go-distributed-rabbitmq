// -----------------------------------------------------------------------------
// Coordinator package used for defining queue listener and event aggregator.
// -----------------------------------------------------------------------------
package coordinator

import (
	"godistributed-rabbitmq/common"
	"goutils/slices"
)

type Any interface{}

// -----------------------------------------------------------------------------
// EventSubscriber - Interface for providing trigger subscription for event.
// -----------------------------------------------------------------------------
type EventSubscriber interface {
	Subscribe(event common.Event, trigger func(Any))
}

// -----------------------------------------------------------------------------
// EventUnsubscriber - Interface for providing trigger unsubscription from event.
// -----------------------------------------------------------------------------
type EventUnsubscriber interface {
	Unsubscribe(event common.Event, trigger func(Any))
}

// -----------------------------------------------------------------------------
// EventRaiser - Interface that combines EventSubscriber and EventUnsubscriber
// -----------------------------------------------------------------------------
type EventRaiser interface {
	EventSubscriber
	EventUnsubscriber
}

// -----------------------------------------------------------------------------
// EventAggregator - Struct that contains the logic for
// -----------------------------------------------------------------------------
type EventAggregator struct {

	// Mapping event name into trigger function.
	triggers map[common.Event][]func(Any)
}

// -----------------------------------------------------------------------------
// NewAggregator - Creates new event aggregator.
// -----------------------------------------------------------------------------
func NewAggregator() *EventAggregator {
	return &EventAggregator{
		triggers: make(map[common.Event][]func(Any)),
	}
}

// -----------------------------------------------------------------------------
// Subscribe - Method used to add new trigger callback for the given event.
//
// event - Name of the event of interest.
// trigger - Callback function triggered for the given event.
// -----------------------------------------------------------------------------
func (aggregator *EventAggregator) Subscribe(event common.Event, trigger func(Any)) {
	aggregator.triggers[event] = append(aggregator.triggers[event], trigger)
}

// -----------------------------------------------------------------------------
// Unsubscribe - Method used to remove trigger callback for the given event.
//
// event - Name of the event of interest.
// trigger - Callback function triggered for the given event.
// -----------------------------------------------------------------------------
func (aggregator *EventAggregator) Unsubscribe(event common.Event, trigger func(Any)) {
	if triggers := aggregator.triggers[event]; triggers != nil {
		remove := slices.IndexOf(len(triggers), func(i int) bool {
			// Predicate for comparing function addresses and removing pased trigger function.
			return &triggers[i] == &trigger
		})

		triggers = append(triggers[:remove], triggers[remove+1:]...)
	}
}

// -----------------------------------------------------------------------------
// Publish - Method used to send event data to all registered consumers.
//
// event - Name of the event of interest.
// data - Event data that is being sent.
// -----------------------------------------------------------------------------
func (aggregator *EventAggregator) Publish(event common.Event, data Any) {

	if triggers := aggregator.triggers[event]; triggers != nil {
		for _, trigger := range triggers {
			// We are sending data by value not refference so that it gets copied.
			trigger(data)
		}
	}
}
