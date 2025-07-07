package kafka

import (
	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ConvertHeaders(headers []*sarama.RecordHeader) []kafka.Header {
	var result []kafka.Header
	for _, h := range headers {
		result = append(result, kafka.Header{
			Key:   string(h.Key),
			Value: h.Value,
		})
	}

	return result
}
