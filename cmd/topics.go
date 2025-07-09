package cmd

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/philipparndt/go-logger"
	"gokcat/config"
	"gokcat/internal/kafka"
	"sort"

	"github.com/spf13/cobra"
)

// topicsCmd represents the topics command
var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List all Kafka topics",
	Long:  `List all Kafka topics available on the configured Kafka cluster.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" && systemAlias == "" {
			return errors.New("you must specify a config file or system alias")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if systemAlias != "" {
			configFile = "~/.config/gokcat/" + systemAlias + "/config.json"
		}

		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			logger.Panic("Failed to load config", configFile, ",", err)
		}

		runTopics(cfg)
	},
}

func init() {
	rootCmd.AddCommand(topicsCmd)
	topicsCmd.Flags().StringVarP(&topic, "topic", "t", "", "Kafka topic to consume messages from")
	topicsCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the configuration file")
	topicsCmd.Flags().StringVarP(&systemAlias, "systemAlias", "s", "", "System alias")
}

func runTopics(cfg config.Config) {
	tlsConfig, err := kafka.NewTLSConfig(cfg.Certs.ClientCert, cfg.Certs.ClientKey, cfg.Certs.Ca, cfg.Certs.Insecure)
	if err != nil {
		logger.Panic("Failed to create TLS config", err)
	}

	kConfig := sarama.NewConfig()
	kConfig.Net.TLS.Enable = true
	kConfig.Net.TLS.Config = tlsConfig

	client, err := sarama.NewClient([]string{cfg.Broker}, kConfig)
	if err != nil {
		logger.Panic("Failed to create client", err)
	}
	defer client.Close()

	// Get list of topics
	topics, err := client.Topics()
	if err != nil {
		logger.Panic("Failed to get topics", err)
	}

	if len(topics) == 0 {
		fmt.Println("No topics found")
		return
	}

	// Sort topics alphabetically
	sort.Strings(topics)

	fmt.Printf("Found %d topics:\n", len(topics))
	for _, topic := range topics {
		fmt.Println(topic)
	}
}
