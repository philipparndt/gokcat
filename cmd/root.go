package cmd

import (
	"errors"
	"gokcat/config"
	"os"

	"github.com/philipparndt/go-logger"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gokcat",
	Short: "Print messages from a Kafka topic",
	Long:  `Print messages from a Kafka topic.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" && systemAlias == "" {
			return errors.New("you must specify a config file or system alias")
		}
		if topic == "" {
			return errors.New("you must specify a topic to cat")
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

		runCat(topic, cfg, follow, tail)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var topic string
var configFile string
var systemAlias string
var follow bool
var tail int

func init() {
	rootCmd.Flags().StringVarP(&topic, "topic", "t", "", "Kafka topic to consume messages from")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the configuration file")
	rootCmd.Flags().StringVarP(&systemAlias, "systemAlias", "s", "", "System alias")
	rootCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow the topic (like tail -f)")
	rootCmd.Flags().IntVar(&tail, "tail", 0, "Read the last n messages from the topic")
}
