package cmd

import (
	"fmt"
	"gokcat/config"
	"gokcat/internal/kafka"
	"gokcat/internal/kafka/schemaRegistry"
	"gokcat/kafkaUtil"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/philipparndt/go-logger"
)

func runCat(topic string, cfg config.Config, follow bool, tail int) {
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

	var startOffset int64
	if tail > 0 {
		// For tail, try to start from a smart position to minimize processing
		startOffset = calculateTailStartOffset(client, partition, tail, latestOffset)
		logger.Info("Tailing last", strconv.Itoa(tail), "messages from offset", strconv.Itoa(int(startOffset)))
	} else if follow {
		startOffset = sarama.OffsetNewest
		latestOffset = -1
		logger.Info("Following topic, press Ctrl+C to exit")
	} else {
		startOffset = sarama.OffsetOldest
		logger.Info("Consuming until offset", strconv.Itoa(int(latestOffset)))
	}

	if latestOffset == 0 && !follow {
		logger.Info("No messages found in topic", topic)
		fmt.Println("[]")
		return
	}

	pc, err := consumer.ConsumePartition(topic, partition, startOffset)
	if err != nil {
		logger.Panic("Failed to consume partition", err)
	}
	defer pc.Close()

	fmt.Println("[")

	if tail > 0 {
		kafkaUtil.Tail(topic, tail, follow, pc, deserializer, startOffset, latestOffset)
	} else {
		kafkaUtil.Cat(topic, follow, pc, deserializer, latestOffset)
	}

	fmt.Println()
	fmt.Println("]")
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

func calculateTailStartOffset(client sarama.Client, partition int32, tail int, latestOffset int64) int64 {
	if latestOffset <= 0 {
		return sarama.OffsetOldest
	}

	// Get the oldest available offset
	oldestOffset, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
	if err != nil {
		logger.Debug("Failed to get oldest offset, starting from beginning")
		return sarama.OffsetOldest
	}

	totalOffsetRange := latestOffset - oldestOffset + 1

	// If the range is small, start from the beginning
	if totalOffsetRange <= int64(tail*2) {
		return oldestOffset
	}

	// Use a conservative heuristic: assume reasonable message density
	// Start from much closer to the end to minimize processing
	estimatedStart := latestOffset - int64(tail*5) // 5x buffer for gaps
	if estimatedStart < oldestOffset {
		estimatedStart = oldestOffset
	}

	return estimatedStart
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
