package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Commands for managing configuration settings for the Slurm REST CLI.`,
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value from the config file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := viper.Get(key)
		if value == nil {
			fmt.Printf("Configuration key '%s' not found\n", key)
			os.Exit(1)
		}
		fmt.Printf("%v\n", value)
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value in the config file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		// Load existing config
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configFile := filepath.Join(home, ".config", "slurm-rest-cli.yaml")
		configMap := make(map[string]interface{})

		// Read existing config if it exists
		if _, err := os.Stat(configFile); err == nil {
			configData, err := os.ReadFile(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
				os.Exit(1)
			}

			if err := yaml.Unmarshal(configData, &configMap); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
				os.Exit(1)
			}
		}

		// Set the value in the config map
		setNestedValue(configMap, key, value)

		// Write the updated config back to file
		configDir := filepath.Join(home, ".config")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		yamlData, err := yaml.Marshal(configMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling config: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(configFile, yamlData, 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration key '%s' set to '%s'\n", key, value)
	},
}

// configListCmd represents the config list command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Long:  `List all configuration values from the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		allSettings := viper.AllSettings()
		yamlData, err := yaml.Marshal(allSettings)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(yamlData))
	},
}

// Helper function to set a nested value in a map
func setNestedValue(configMap map[string]interface{}, key string, value string) {
	parts := strings.Split(key, ".")
	current := configMap

	// Navigate to the nested map
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}
		if nestedMap, ok := current[part].(map[string]interface{}); ok {
			current = nestedMap
		} else {
			// Convert to map if it's not already
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		}
	}

	// Set the value
	current[parts[len(parts)-1]] = value
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
}
