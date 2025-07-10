package kafkaUtil

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/philipparndt/go-logger"
	"gokcat/internal/kafka/schemaRegistry"
	"gokcat/message"
	"strconv"
)

func Tail(topic string, amount int, follow bool, pc sarama.PartitionConsumer, deserializer schemaRegistry.Deserializer, latestOffset int64) {
	// For tail mode, use a simple sliding window approach
	// We'll consume from wherever we started and keep only the last N messages
	ctr := 0
	var messages []string

	for msg := range pc.Messages() {
		if ctr%1000 == 0 && ctr > 0 {
			logger.Debug("Processed " + strconv.Itoa(ctr) + " messages")
		}
		ctr++

		var payloadData interface{} = nil
		var schema *schemaRegistry.Schema

		// Check if the message starts with magic byte (0x0 for Confluent wire format)
		if len(msg.Value) < 5 || msg.Value[0] != 0x0 {
			payloadDataString := decodeBase64OrRaw(msg.Value)
			payloadData = decodeJSONOrRaw(payloadDataString)
		} else {
			schema, err := deserializer.LoadSchemaInfo(topic, msg)
			if err != nil {
				logger.Panic("Failed to load schema info", err)
			}

			payloadData = deserializer.Deserialize(schema, msg.Value[5:])
		}

		out := message.New(schema, payloadData, msg)

		jsonPayload, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			logger.Error("Failed to marshal payload data to JSON", "error", err)
			continue
		}

		// Keep only the last 'tail' messages in a sliding window
		messages = append(messages, string(jsonPayload))
		if len(messages) > amount {
			messages = messages[1:] // Remove the first element
		}

		if msg.Offset >= latestOffset-1 {
			break
		}
	}

	// Output the messages (which are already the last N)
	for i, msg := range messages {
		fmt.Print(msg)
		if i < len(messages)-1 {
			fmt.Println(",")
		}
	}
}
