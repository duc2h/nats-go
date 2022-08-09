package main

import (
	"encoding/json"
	"log"
	"nats-go/model"
	"strconv"

	"github.com/nats-io/nats.go"
)

const (
	stream  = "stream1"
	durable = "durable1"
	queue   = "queue1"
	subject = "subject1"
	deliver = "deliver1"
)

func main() {

	// Connect to NATS
	opt, err := nats.NkeyOptionFromSeed("../nkey-cert.yaml")
	if err != nil {
		log.Fatalln("error nkey: ", err)
	}
	nc, err := nats.Connect("nats://localhost:4223", opt)
	if err != nil {
		log.Fatalln(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalln(err)
	}

	err = createOrder(js)
	if err != nil {
		log.Fatalln(err)
	}
}

func createOrder(js nats.JetStreamContext) error {
	var order model.Order
	for i := 1; i <= 3; i++ {
		order = model.Order{
			OrderID:    i,
			CustomerID: "Cust-" + strconv.Itoa(i),
			Status:     "created2sdsd",
		}
		orderJSON, _ := json.Marshal(order)
		_, err := js.Publish(subject, orderJSON)
		if err != nil {
			return err
		}
		log.Printf("Order with OrderID:%d has been published\n", i)
	}
	return nil
}
