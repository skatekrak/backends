package api

type JSONError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
} // @Name JSONError

type JSONMessage struct {
	Message string `json:"message,omitempty"`
} // @Name JSONMessage
