package models

// GenericResponse represents a generic API response.
type GenericResponse struct {
	Errors []string `json:"errors,omitempty"`
}
