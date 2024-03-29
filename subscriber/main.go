package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nats-go/model"
	"runtime"
	"time"

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

	t := time.Now()
	fmt.Println(t)
	createStream(js)
	createConsumer(js, t)
	createSub(js, t)

	runtime.Goexit()
}

// create stream
func createStream(js nats.JetStreamContext) {
	info, _ := js.StreamInfo(stream)
	if info != nil {
		js.DeleteStream(stream)
	}

	cfg := nats.StreamConfig{
		Name:      stream,
		Retention: nats.InterestPolicy,
		Replicas:  1,
		Subjects:  []string{subject},
	}
	// info, err := js.StreamInfo(stream)
	// if err != nil {
	_, err := js.AddStream(&cfg)
	// }
	if err != nil {
		log.Fatal("create stream error: ", err.Error())
	}

	log.Println("=========================== create stream success")
}

// create consumer
func createConsumer(js nats.JetStreamContext, now time.Time) {
	info, _ := js.ConsumerInfo(stream, durable)
	if info != nil {
		js.DeleteConsumer(stream, durable)
	}

	cfg := nats.ConsumerConfig{
		Durable:        durable,
		FilterSubject:  subject,
		MaxDeliver:     5,
		DeliverSubject: deliver,
		AckWait:        time.Second,
		DeliverPolicy:  nats.DeliverNewPolicy,
		DeliverGroup:   queue,
		AckPolicy:      nats.AckExplicitPolicy,
	}
	_, err := js.AddConsumer(stream, &cfg)
	if err != nil {
		log.Fatal("create consumer error: ", err.Error())
	}
	log.Println("=========================== create consumer success")

}

// register queue-subscribe
func createSub(js nats.JetStreamContext, now time.Time) {
	_, err := js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
		msg.Ack()

	}, nats.Bind(stream, durable),
		nats.ManualAck(),
		nats.MaxDeliver(5),
		nats.DeliverSubject(deliver),
		nats.AckWait(time.Second),
	)

	if err != nil {
		log.Fatal("create createSub error: ", err.Error())
	}

	log.Println("=========================== create createSub success")

}
