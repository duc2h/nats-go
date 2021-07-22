package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"

	"nats-go/model"
)

const (
	streamName     = "student"
	streamSubjects = "student.*"
	subjectName    = "student.Created"
)

func main() {
	// Connect to NATS
	opt, err := nats.NkeyOptionFromSeed("../nkey-cert.yaml")
	if err != nil {
		log.Fatalln("error nkey: ", err)
	}
	nc, err := nats.Connect("nats://localhost:4223", opt)

	if err != nil {
		log.Fatal("connect err: ", err)
	}
	fmt.Println("nc--- ", nc.Status())
	// Creates JetStreamContext
	js, err := nc.JetStream()
	fmt.Println(err)
	checkErr(err)
	// Creates stream
	err = createStream(js)
	checkErr(err)
	// Create orders by publishing messages
	err = createOrder(js)
	checkErr(err)
}

// createOrder publishes stream of events
// with subject "ORDERS.created"
func createOrder(js nats.JetStreamContext) error {
	var order model.Order
	for i := 1; i <= 3; i++ {
		order = model.Order{
			OrderID:    i,
			CustomerID: "Cust-" + strconv.Itoa(i),
			Status:     "created2sdsd",
		}
		orderJSON, _ := json.Marshal(order)
		_, err := js.Publish(subjectName, orderJSON)
		if err != nil {
			return err
		}
		log.Printf("Order with OrderID:%d has been published\n", i)
	}
	return nil
}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext) error {
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	fmt.Println("stream: ", stream)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{streamSubjects},
			Retention: nats.InterestPolicy,
			MaxAge:    time.Minute * 10000,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("check error: ", err)
	}
}
