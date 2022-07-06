package api

type JSONError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
} // @Name JSONError
