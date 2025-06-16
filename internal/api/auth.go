package api

import (
	"fmt"
	"net/http"
	"strings"
)

// TokenAuth represents authentication using X-SLURM-USER-TOKEN header.
type TokenAuth struct {
	Token string
}

// AddAuthHeaders adds the token authentication header to the request.
func (a *TokenAuth) AddAuthHeaders(req *http.Request) error {
	if a.Token == "" {
		return fmt.Errorf("authentication token is required")
	}
	req.Header.Set("X-SLURM-USER-TOKEN", a.Token)
	return nil
}

// UserTokenAuth represents authentication using X-SLURM-USER-NAME and X-SLURM-USER-TOKEN headers.
type UserTokenAuth struct {
	Username string
	Token    string
}

// AddAuthHeaders adds the username and token authentication headers to the request.
func (a *UserTokenAuth) AddAuthHeaders(req *http.Request) error {
	if a.Username == "" {
		return fmt.Errorf("username is required")
	}
	if a.Token == "" {
		return fmt.Errorf("authentication token is required")
	}
	req.Header.Set("X-SLURM-USER-NAME", a.Username)
	req.Header.Set("X-SLURM-USER-TOKEN", a.Token)
	return nil
}

// JWTAuth represents authentication using JWT Bearer token.
type JWTAuth struct {
	JWT string
}

// AddAuthHeaders adds the JWT Bearer authentication header to the request.
func (a *JWTAuth) AddAuthHeaders(req *http.Request) error {
	if a.JWT == "" {
		return fmt.Errorf("JWT token is required")
	}
	
	// Add Bearer prefix if not already present
	jwt := a.JWT
	if !strings.HasPrefix(strings.ToLower(jwt), "bearer ") {
		jwt = "Bearer " + jwt
	}
	
	req.Header.Set("Authorization", jwt)
	return nil
}

// NewAuthenticator creates an appropriate authenticator based on the provided credentials.
func NewAuthenticator(username, token, jwt string) Authenticator {
	// JWT authentication takes precedence
	if jwt != "" {
		return &JWTAuth{JWT: jwt}
	}
	
	// User + Token authentication
	if username != "" && token != "" {
		return &UserTokenAuth{Username: username, Token: token}
	}
	
	// Token-only authentication
	if token != "" {
		return &TokenAuth{Token: token}
	}
	
	// Default to no authentication (will likely fail)
	return &TokenAuth{Token: ""}
}
