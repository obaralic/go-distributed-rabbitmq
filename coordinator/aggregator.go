// -----------------------------------------------------------------------------
// Coordinator package used for defining queue listener and event aggregator.
// -----------------------------------------------------------------------------
package coordinator

import (
	"godistributed-rabbitmq/common/dto"
	"goutils/slices"
)

type Event string

// -----------------------------------------------------------------------------
// EventAggregator - Struct that contains the logic for
// -----------------------------------------------------------------------------
type EventAggregator struct {

	// Mapping event name into trigger function.
	triggers map[Event][]func(dto.EventData)
}

// -----------------------------------------------------------------------------
// NewAggregator - Creates new event aggregator.
// -----------------------------------------------------------------------------
func NewAggregator() *EventAggregator {
	return &EventAggregator{
		triggers: make(map[Event][]func(dto.EventData)),
	}
}

// -----------------------------------------------------------------------------
// Subscribe - Method used to add new trigger callback for the given event.
//
// event - Name of the event of interest.
// trigger - Callback function triggered for the given event.
// -----------------------------------------------------------------------------
func (aggregator *EventAggregator) Subscribe(event Event, trigger func(dto.EventData)) {
	aggregator.triggers[event] = append(aggregator.triggers[event], trigger)
}

// -----------------------------------------------------------------------------
// Unsubscribe - Method used to remove trigger callback for the given event.
//
// event - Name of the event of interest.
// trigger - Callback function triggered for the given event.
// -----------------------------------------------------------------------------
func (aggregator *EventAggregator) Unsubscribe(event Event, trigger func(dto.EventData)) {
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
func (aggregator *EventAggregator) Publish(event Event, data dto.EventData) {

	if triggers := aggregator.triggers[event]; triggers != nil {
		for _, trigger := range triggers {
			// We are sending data by value not refference so that it gets copied.
			trigger(data)
		}
	}
}
