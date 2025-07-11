package kafkaUtil

import (
	"encoding/json"
	"fmt"
	"gokcat/internal/kafka/schemaRegistry"
	"gokcat/message"
	"time"

	"github.com/IBM/sarama"
	"github.com/philipparndt/go-logger"
)

func Tail(topic string, amount int, follow bool, pc sarama.PartitionConsumer, deserializer schemaRegistry.Deserializer, startOffset int64, latestOffset int64) {
	// Check if the topic is empty
	if startOffset >= latestOffset {
		logger.Info("The topic is empty or compacted with no messages available.")
		return
	}

	ctr := 0
	var messages []string
	timeout := time.After(10 * time.Second) // Timeout for empty topics

	exitLoop := false
	for !exitLoop {
		select {
		case msg := <-pc.Messages():
			if msg == nil {
				continue
			}

			if ctr >= amount {
				exitLoop = true
				break // Exit the loop cleanly
			}

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

			messages = append(messages, string(jsonPayload))
			ctr++

			if follow {
				fmt.Print(string(jsonPayload))
				if ctr > 1 {
					fmt.Println(",")
				}
			}

		case <-timeout:
			logger.Warn("Timeout reached while waiting for messages")
			exitLoop = true
			break // Exit the loop cleanly
		}
	}

	// Output the collected messages (only relevant for non-follow mode)
	if !follow {
		for i, msg := range messages {
			fmt.Print(msg)
			if i < len(messages)-1 {
				fmt.Println(",")
			}
		}
	}
}
