package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var systemsCmd = &cobra.Command{
	Use:   "systems",
	Short: "List all configured systems",
	Long:  GetHelp("systems"),
	Run: func(cmd *cobra.Command, args []string) {
		runSystems()
	},
}

func init() {
	rootCmd.AddCommand(systemsCmd)
}

func runSystems() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error determining home directory:", err)
		os.Exit(1)
	}

	baseDir := filepath.Join(home, ".config", "gokcat")
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No systems configured.")
			return
		}
		fmt.Fprintln(os.Stderr, "Error reading config directory:", err)
		os.Exit(1)
	}

	var systems []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		cfgPath := filepath.Join(baseDir, entry.Name(), "config.json")
		if _, err := os.Stat(cfgPath); err == nil {
			systems = append(systems, entry.Name())
		}
	}

	if len(systems) == 0 {
		fmt.Println("No systems configured.")
		return
	}

	sort.Strings(systems)
	for _, s := range systems {
		fmt.Println(s)
	}
}
