package cmd

import (
	"encoding/json"
	"fmt"
	"gokcat/version"

	"github.com/spf13/cobra"
)

var jsonOutput bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information of gokcat, including build details.`,
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()

		if jsonOutput {
			output, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling version info: %v\n", err)
				return
			}
			fmt.Println(string(output))
		} else {
			fmt.Println(info.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version information in JSON format")
}
