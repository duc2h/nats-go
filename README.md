# nats-go
This repository for project demo with nats jetstream

## Nats-server
Nats server based on "At most once". It's mean if any subscriber is down while forwarding message from publisher to subcriber systems. Then those messages will be missed in subscriber systems and thus there is no guarantee of delivery for the publisher systems.

## Nats-streaming
Nats streaming based on "at-least-once-delivery". It's mean when you publish a message, the published message persist into a customizable storage so that we can replay the messages. 
### Some concept in nats-streaming
* ACK
* Wildcard
* Sequence numbers 
* Queue groups
* Channels
* Rate limiting/matching

## Nats-jetstream
It have some concept like nats-streaming, but it's better with the capabilities of distributed security, multi-tenancy and horizontally scaling.
### What is distributed event streaming event?
* The streaming systems let you capture the streams of event from distributed system and microservices, software applications,... and persist these streams into persisten stores. So you can replay it for retrieval, process, and reactive to those event streams by using an event-driven architecture. Nats jetstream is a distributed streaming system with publish/subcribe as the messaging pattern.

### Some concept
* Streams: In jetstream, streams of events are stored as `streams`. Docs: https://docs.nats.io/jetstream/concepts/streams
* Consumers: In jetstream, consumers can either be push or pull mode. Docs: https://docs.nats.io/jetstream/concepts/consumers

### Pull based consumer and Push based consumer
* Jetstream provides two kind of consumer (subcriber) systems: Pull based consumer and Push based consumer. The Pull based consumer let jetstream pull the messages from consumer systems. Pull based consumer systems are like work queues. Because the jetstream provides a ACK(acknowledment) mechansim, you can easily scale Pull based consumer systems horizontally without the problem of duplication of messages. Pull based subcription is new to the Nats ecosystem. The Push based consumer let jetstream pushing the messages to consumer systems. Which can be a good choice for monitoring systems. Docs: https://docs.nats.io/jetstream/concepts


### How to tracing nats-jetstream
We use `nats-box` to trace nats-jetstream, run docker to start nats-box:
`docker run --rm --network host -it natsio/nats-box:latest`

Then, connecting to nats server:
`nats context save s1 --user=admin --password=admin --server=nats://127.0.0.1:4223 --select`