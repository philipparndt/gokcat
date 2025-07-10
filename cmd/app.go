package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/philipparndt/go-logger"
	"gokcat/config"
	"gokcat/internal/kafka"
	"gokcat/internal/kafka/schemaRegistry"
	"gokcat/message"
	"strconv"
)

func runCat(topic string, cfg config.Config, follow bool) {
	partition := int32(0)
	tlsConfig, err := kafka.NewTLSConfig(cfg.Certs.ClientCert, cfg.Certs.ClientKey, cfg.Certs.Ca, cfg.Certs.Insecure)
	if err != nil {
		logger.Panic("Failed to create TLS config", err)
	}

	kConfig := sarama.NewConfig()
	kConfig.Net.TLS.Enable = true
	kConfig.Net.TLS.Config = tlsConfig

	sr := schemaRegistry.New(cfg.SchemaRegistry.Url,
		cfg.SchemaRegistry.Username,
		cfg.SchemaRegistry.Password,
		cfg.SchemaRegistry.Insecure,
	)

	deserializer := sr.NewDeserializer()
	logger.Debug("Created deserializer successfully")

	client, err := sarama.NewClient([]string{cfg.Broker}, kConfig)
	if err != nil {
		logger.Panic("Failed to create client", err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		logger.Panic("Failed to create consumer from client", err)
	}
	defer consumer.Close()

	latestOffset, err := getOffset(client, partition)
	if err != nil {
		logger.Panic("Failed to get latest offset", err)
	}

	if latestOffset == 0 {
		logger.Info("No messages found in topic", topic)
		if !follow {
			fmt.Println("[]")
			return
		}
	}

	if follow {
		latestOffset = -1
		logger.Info("Following topic, press Ctrl+C to exit")
	} else {
		logger.Info("Consuming until offset", strconv.Itoa(int(latestOffset)))
	}

	pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		logger.Panic("Failed to consume partition", err)
	}
	defer pc.Close()

	fmt.Println("[")

	for msg := range pc.Messages() {
		// Check if the message is empty or too short to contain a schema ID
		if len(msg.Value) < 5 {
			logger.Error("Message too short to contain schema ID", "length", len(msg.Value))
			if msg.Offset >= latestOffset-1 {
				logger.Info("Reached end of topic. Exiting.")
				break
			}
			continue
		}

		var payloadData interface{} = nil
		var schema *schemaRegistry.Schema

		// Check if the message starts with magic byte (0x0 for Confluent wire format)
		if msg.Value[0] != 0x0 {
			payloadDataString := decodeBase64OrRaw(msg.Value)
			payloadData = decodeJSONOrRaw(payloadDataString)
		} else {
			schema, err = deserializer.LoadSchemaInfo(topic, msg)
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
		fmt.Print(string(jsonPayload))

		if msg.Offset >= latestOffset-1 && !follow {
			logger.Info("Reached end of topic. Exiting.")
			break
		}

		fmt.Println(",")
	}

	fmt.Println()
	fmt.Println("]")
}

// decodeBase64OrRaw tries to decode the input as base64, returns raw bytes if not base64
func decodeBase64OrRaw(data []byte) []byte {
	decoded, err := strconv.Unquote("\"" + string(data) + "\"")
	if err == nil {
		return []byte(decoded)
	}
	return data
}

// decodeJSONOrRaw tries to decode the input as JSON, returns raw bytes if not JSON
func decodeJSONOrRaw(data []byte) interface{} {
	var v interface{}
	if err := json.Unmarshal(data, &v); err == nil {
		return v
	}
	return data
}

func getOffset(client sarama.Client, partition int32) (int64, error) {
	// Get the latest offset (the "high watermark")
	latestOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return 0, err
	}

	if latestOffset == 0 {
		return 0, nil
	}

	return latestOffset - 1, nil
}
