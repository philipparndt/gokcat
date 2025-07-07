package message

import (
	"github.com/IBM/sarama"
	"gokcat/internal/kafka/schemaRegistry"
	"time"
)

type Message struct {
	Schema struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"schema"`

	Metadata struct {
		Offset    int64  `json:"offset"`
		Timestamp string `json:"timestamp"`
		Key       string `json:"key"`
	} `json:"metadata,omitempty"`

	Payload interface{} `json:"payload"`
}

func New(schema *schemaRegistry.Schema, payloadData map[string]interface{}, msg *sarama.ConsumerMessage) Message {
	out := Message{}
	out.Schema.Id = schema.ID
	out.Schema.Name = schema.Name
	out.Schema.Namespace = schema.Namespace
	out.Metadata.Key = string(msg.Key)
	out.Metadata.Timestamp = msg.Timestamp.Format(time.RFC3339)
	out.Metadata.Offset = msg.Offset

	out.Payload = payloadData

	return out
}
