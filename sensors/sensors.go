// -----------------------------------------------------------------------------
// Package containing implementation of configurable sensors
// that will collect or generate data and deploy it to server
// for further processing. This package is also a main package
// since sensors will be able to execute as stand-alone application.
// -----------------------------------------------------------------------------
package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const millis = 1000

// Unique sensor name used as a routing key.
var name = flag.String("name", "SimulatedSensor", "Name of the sensor")

// Number of sensor readout samples per second.kada se zalepsi
var frequency = flag.Uint("freq", 5, "Update frequency in cycles/sec")

// Maximum value of sensor readout.
var max = flag.Float64("max", 5., "Maximum value for generated readouts")

// Minimum value of sensor readout.
var min = flag.Float64("min", 1., "Minimum value for generated readouts")

// Minimum value of sensor readout.
var deviation = flag.Float64("dev", 0.1, "Maximum allowed readout deviation")

// Random number generator used for generating readouts.
var random *rand.Rand

// Sensor value that contains readout.
var value float64

// Nominal value of the sensor readout.
var nominal float64

func init() {
	flag.Parse()
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	value = random.Float64()*(*max-*min) + *min
	nominal = (*max-*min)/2 + *min
}

func main() {
	work()
}

// -----------------------------------------------------------------------------
// work - Provides work routine of the given sensor.
// -----------------------------------------------------------------------------
func work() {

	// Create AMQP conenction, channel and queue used for transfering readouts.
	connection, channel := common.GetChannel(common.URL_GUEST)
	defer connection.Close()
	defer channel.Close()

	// Publish sensor name to the system using sensor advertisement queue.
	queue := common.GetQueue(*name, channel)
	common.Advertise(*name, common.GetSensorQueue(channel), channel)

	// Create reusable data buffer and surrounding encoder.
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	// Wait for the timer tick by blocking on unbuffered channel.
	for range tick() {
		sample()

		// Reset buffer content and encode readout.
		buffer.Reset()
		readout := dto.NewReadout(*name, value, time.Now())
		encoder.Encode(readout)

		common.Send(buffer.Bytes(), queue, channel)
		log.Printf("Sensor: %s sent value: %v\n", *name, readout.Value)
	}
}

// -----------------------------------------------------------------------------
// sample - Calculates sensor value that will become new readout sample.
// -----------------------------------------------------------------------------
func sample() {
	var minDeviation, maxDeviation float64

	// Depending on current readout value determin deviation bounderies.
	if value < nominal {
		minDeviation = *deviation * (value - *min) / (nominal - *min) * -1
		maxDeviation = *deviation
	} else {
		minDeviation = *deviation * -1
		maxDeviation = *deviation * (*max - value) / (*max - nominal)
	}

	// Sample is incremented by the random value within the deviation bounderies.
	value += random.ExpFloat64()*(maxDeviation-minDeviation) + minDeviation
}

// -----------------------------------------------------------------------------
// tick - Provides channel wiht periodic time tick used for sensor sampling.
// -----------------------------------------------------------------------------
func tick() (signal <-chan time.Time) {
	// Convert cycles/sec into milliseconds/cycle.
	duration, _ := time.ParseDuration(strconv.Itoa(millis/int(*frequency)) + "ms")

	// Create periodic sampling by getting a timer tick channel.
	return time.Tick(duration)
}
