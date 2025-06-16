package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	apiURL      string
	token       string
	jwt         string
	username    string
	outputFormat string
	showCurl    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "srest",
	Short: "Modern CLI for Slurm REST API",
	Long: `srest is a modern command-line interface for the Slurm REST API.
It provides a familiar interface similar to traditional Slurm commands
while leveraging the power of the REST API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/slurm-rest-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "url", "", "Slurm REST API URL")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Authentication token")
	rootCmd.PersistentFlags().StringVar(&jwt, "jwt", "", "JWT for authentication")
	rootCmd.PersistentFlags().StringVar(&username, "user", "", "Username for authentication")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "table", "Output format (json, yaml, table)")
	rootCmd.PersistentFlags().BoolVar(&showCurl, "curl", false, "Show equivalent curl command instead of executing")

	// Bind flags to viper
	viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("auth.token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("auth.jwt", rootCmd.PersistentFlags().Lookup("jwt"))
	viper.BindPFlag("auth.username", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("preferences.output_format", rootCmd.PersistentFlags().Lookup("format"))

	// Environment variables
	viper.SetEnvPrefix("SLURM_REST")
	viper.BindEnv("api.url", "SLURM_REST_API_URL")
	viper.BindEnv("auth.token", "SLURM_REST_TOKEN")
	viper.BindEnv("auth.jwt", "SLURM_REST_JWT")
	viper.BindEnv("auth.username", "SLURM_REST_USER")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".slurm-rest-cli" (without extension).
		viper.AddConfigPath(home + "/.config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("slurm-rest-cli")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
