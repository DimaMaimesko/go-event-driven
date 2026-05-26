package main

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageHeader struct {
	ID         string `json:"id"`
	EventName  string `json:"event_name"`
	OccurredAt string `json:"occurred_at"`
}

type ProductOutOfStock struct {
	Header    MessageHeader `json:"header"`
	ProductID string        `json:"product_id"`
}

type ProductBackInStock struct {
	Header    MessageHeader `json:"header"`
	ProductID string        `json:"product_id"`
	Quantity  int           `json:"quantity"`
}

type Publisher struct {
	pub message.Publisher
}

func NewPublisher(pub message.Publisher) Publisher {
	return Publisher{
		pub: pub,
	}
}

func (p Publisher) PublishProductOutOfStock(productID string) error {
	event_name := "ProductOutOfStock"
	event := ProductOutOfStock{
		Header:    NewMessageHeader(event_name),
		ProductID: productID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}

func (p Publisher) PublishProductBackInStock(productID string, quantity int) error {
	event_name := "ProductBackInStock"
	event := ProductBackInStock{
		Header:    NewMessageHeader(event_name),
		ProductID: productID,
		Quantity:  quantity,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}

func NewMessageHeader(eventName string) MessageHeader {
	return MessageHeader{
		ID:         watermill.NewUUID(),
		EventName:  eventName,
		OccurredAt: time.Now().Format(time.RFC3339),
	}
}
