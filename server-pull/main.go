package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"nats-go/model"
)

const (
	subSubjectName = "Student"
	pubSubjectName = "Student.approved"
	streamName     = "Student"
)

func main() {
	// Connect to NATS
	opt, err := nats.NkeyOptionFromSeed("../nkey-cert.yaml")
	if err != nil {
		log.Fatalln("error nkey: ", err)
	}
	nc, err := nats.Connect("nats://localhost:4223", opt)

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	err = js.DeleteStream("Student")
	if err == nil {
		fmt.Println("delete stream success")
		return
	}

	// err = createStream(js)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	sub, err := js.PullSubscribe("Student.*", "student-pull",
		nats.PullMaxWaiting(1),
		nats.MaxDeliver(5),
		nats.AckWait(time.Second*5))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	msgs, _ := sub.Fetch(5, nats.Context(ctx))
	for _, msg := range msgs {
		msg.Ack()
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("order-review service")
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
		// reviewOrder(js, order)
	}

	time.Sleep(time.Minute * 100)
}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext) error {
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}

	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, subSubjectName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{subSubjectName},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// reviewOrder reviews the order and publishes ORDERS.approved event
func reviewOrder(js nats.JetStreamContext, order model.Order) {
	// Changing the Order status
	order.Status = "approved"
	orderJSON, _ := json.Marshal(order)
	_, err := js.Publish(pubSubjectName, orderJSON)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Order with OrderID:%d has been %s\n", order.OrderID, order.Status)
}
