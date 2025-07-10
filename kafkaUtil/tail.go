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
	reachedLatest := false

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

		// Check if we've reached the latest offset
		if msg.Offset >= latestOffset-1 {
			reachedLatest = true
		}

		if follow && reachedLatest {
			// In follow mode, once we've caught up, output messages immediately
			fmt.Print(string(jsonPayload))
			if ctr > 1 {
				fmt.Println(",")
			}
		} else {
			// Keep only the last 'amount' messages in a sliding window
			messages = append(messages, string(jsonPayload))
			if len(messages) > amount {
				messages = messages[1:] // Remove the first element
			}

			// If not following and we've reached the latest offset, break
			if !follow && reachedLatest {
				break
			}
		}
	}

	// Output the collected messages (only relevant for non-follow mode or initial batch in follow mode)
	if !follow || !reachedLatest {
		for i, msg := range messages {
			fmt.Print(msg)
			if i < len(messages)-1 {
				fmt.Println(",")
			}
		}
	}
}
