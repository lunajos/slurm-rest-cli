package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var (
	keycloakServerURL string
	keycloakRealm     string
	keycloakClientID  string
	noStoreCreds      bool
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication related commands",
	Long:  `Commands for managing authentication with Slurm REST API.`,
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Keycloak and obtain JWT token",
	Long: `Login to Keycloak authentication server and obtain a JWT token
for use with the Slurm REST API. Credentials can be provided via
command line flags or will be prompted interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		var configDir string
		var configFile string
		var yamlData []byte
		var err error
		// Get username if not provided
		if username == "" {
			fmt.Print("Enter username: ")
			fmt.Scanln(&username)
		}

		// Get password (not storing in a variable visible in ps)
		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println() // Add newline after password input
		password := string(passwordBytes) // Convert to string for use with authentication

		// Fetch real token from Keycloak
		token, expiry, err := getKeycloakToken(username, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error authenticating with Keycloak: %v\n", err)
			os.Exit(1)
		}

		// Create config directory if it doesn't exist
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configDir = filepath.Join(home, ".config")
		if err = os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}
		configFile = filepath.Join(configDir, "slurm-rest-cli.yaml")

		// Load existing config if present
		existing := make(map[string]interface{})
		if _, err = os.Stat(configFile); err == nil {
			var data []byte
			data, err = os.ReadFile(configFile)
			if err == nil {
				yaml.Unmarshal(data, &existing)
			}
		}

		// Update only the relevant fields
		if existing["auth"] == nil {
			existing["auth"] = map[string]interface{}{}
		}
		authSection := existing["auth"].(map[string]interface{})
		authSection["jwt"] = map[string]interface{}{
			"token":  token,
			"expiry": expiry.Format(time.RFC3339),
		}
		// Merge or update keycloak section, preserving client_secret and other fields if present
		keycloakSection := map[string]interface{}{}
		if val, ok := authSection["keycloak"].(map[string]interface{}); ok {
			for k, v := range val {
				keycloakSection[k] = v
			}
		}
		// Only update these values if they were provided via flags
		if keycloakServerURL != "" {
			keycloakSection["server_url"] = keycloakServerURL
		}
		if keycloakRealm != "" {
			keycloakSection["realm"] = keycloakRealm
		}
		if keycloakClientID != "" {
			keycloakSection["client_id"] = keycloakClientID
		}
		keycloakSection["username"] = username
		authSection["keycloak"] = keycloakSection

		// Optionally update api.url and preferences if flags are set (preserve otherwise)
		if existing["api"] == nil {
			existing["api"] = map[string]interface{}{}
		}
		apiSection := existing["api"].(map[string]interface{})
		if apiURL != "" {
			apiSection["url"] = apiURL
		}
		if existing["preferences"] == nil {
			existing["preferences"] = map[string]interface{}{}
		}
		prefSection := existing["preferences"].(map[string]interface{})
		if outputFormat != "" {
			prefSection["output_format"] = outputFormat
		}
		config := existing

		// Create config directory if it doesn't exist
		home, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configDir = filepath.Join(home, ".config")
		if err = os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		// Write config to file
		configFile = filepath.Join(configDir, "slurm-rest-cli.yaml")
		yamlData, err = yaml.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling config: %v\n", err)
			os.Exit(1)
		}

		if err = os.WriteFile(configFile, yamlData, 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully logged in as %s\n", username)
		fmt.Printf("Configuration saved to %s\n", configFile)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)

	// Login command flags
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	loginCmd.Flags().StringVar(&keycloakServerURL, "keycloak-url", "", "Keycloak server URL")
	loginCmd.Flags().StringVar(&keycloakRealm, "realm", "", "Keycloak realm")
	loginCmd.Flags().StringVar(&keycloakClientID, "client-id", "", "Keycloak client ID")
	loginCmd.Flags().StringVar(&apiURL, "api-url", "", "Slurm REST API URL")
	loginCmd.Flags().StringVar(&outputFormat, "output-format", "", "Default output format")
	loginCmd.Flags().BoolVar(&noStoreCreds, "no-store-creds", false, "Don't store credentials in config file")

	// Bind flags to viper
	viper.BindPFlag("auth.keycloak.server_url", loginCmd.Flags().Lookup("keycloak-url"))
	viper.BindPFlag("auth.keycloak.realm", loginCmd.Flags().Lookup("realm"))
	viper.BindPFlag("auth.keycloak.client_id", loginCmd.Flags().Lookup("client-id"))
}

// getKeycloakToken fetches a JWT token from Keycloak using the password grant type
func getKeycloakToken(username, password string) (string, time.Time, error) {
	// Load config to get Keycloak settings
	config := viper.GetViper()
	
	// Get Keycloak server URL from config or flag
	serverURL := keycloakServerURL
	if serverURL == "" {
		if config.IsSet("auth.keycloak.server_url") {
			serverURL = config.GetString("auth.keycloak.server_url")
		}
	}
	if serverURL == "" {
		return "", time.Time{}, fmt.Errorf("Keycloak server URL not provided")
	}
	
	// Get realm from config or flag
	realm := keycloakRealm
	if realm == "" {
		if config.IsSet("auth.keycloak.realm") {
			realm = config.GetString("auth.keycloak.realm")
		}
	}
	if realm == "" {
		return "", time.Time{}, fmt.Errorf("Keycloak realm not provided")
	}
	
	// Get client ID from config or flag
	clientID := keycloakClientID
	if clientID == "" {
		if config.IsSet("auth.keycloak.client_id") {
			clientID = config.GetString("auth.keycloak.client_id")
		}
	}
	if clientID == "" {
		return "", time.Time{}, fmt.Errorf("Keycloak client ID not provided")
	}
	
	// Get client secret from config
	clientSecret := ""
	if config.IsSet("auth.keycloak.client_secret") {
		clientSecret = config.GetString("auth.keycloak.client_secret")
	}
	
	// Construct token URL
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", strings.TrimSuffix(serverURL, "/"), realm)
	
	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", clientID)
	data.Set("username", username)
	data.Set("password", password)
	
	// Add client secret if available
	if clientSecret != "" {
		data.Set("client_secret", clientSecret)
	}
	
	// Make the request
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to request token: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return "", time.Time{}, fmt.Errorf("authentication failed: %v", errorResp)
	}
	
	// Parse response
	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		TokenType        string `json:"token_type"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse token response: %v", err)
	}
	
	// Calculate expiry time
	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	
	// Return the token and expiry
	return tokenResp.AccessToken, expiry, nil
}
