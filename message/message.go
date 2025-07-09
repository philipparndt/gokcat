package message

import (
	"github.com/IBM/sarama"
	"gokcat/internal/kafka/schemaRegistry"
	"time"
)

type T struct {
	Headers struct {
		ContentType   string `json:"Content-Type"`
		Authorization string `json:"Authorization"`
		Field3        string `json:"1"`
		Field4        string `json:"2"`
	} `json:"headers"`
}

type Message struct {
	Schema struct {
		Id        int    `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	} `json:"schema,omitempty"`

	Metadata struct {
		Offset    int64             `json:"offset"`
		Timestamp string            `json:"timestamp"`
		Key       string            `json:"key"`
		Headers   map[string]string `json:"headers,omitempty"`
	} `json:"metadata,omitempty"`

	Payload interface{} `json:"payload"`
}

func New(schema *schemaRegistry.Schema, payloadData interface{}, msg *sarama.ConsumerMessage) Message {
	out := Message{}
	if schema != nil {
		out.Schema.Id = schema.ID
		out.Schema.Name = schema.Name
		out.Schema.Namespace = schema.Namespace
	}
	out.Metadata.Key = string(msg.Key)
	out.Metadata.Timestamp = msg.Timestamp.Format(time.RFC3339)
	out.Metadata.Offset = msg.Offset

	if msg.Headers != nil {
		out.Metadata.Headers = make(map[string]string, len(msg.Headers))
		for _, header := range msg.Headers {
			if header.Key == nil {
				continue
			}
			out.Metadata.Headers[string(header.Key)] = string(header.Value)
		}
	}

	out.Payload = payloadData

	return out
}
