package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"

	"nats-go/model"
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

	// Create durable consumer monitor
	js.Subscribe("student.Created", func(msg *nats.Msg) {

		fmt.Println("Reply: ", msg.Reply)

		// err = msg.Ack()
		fmt.Println(err)
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
	}, nats.Durable("durable-queue-push1"), nats.ManualAck(), nats.MaxDeliver(5), nats.AckWait(time.Second))

	runtime.Goexit()

}
