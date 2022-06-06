package handlers

// swagger:model
type Error struct {
	// Code code error.
	Code string `json:"code"`

	// Message message error.
	Message string `json:"message"`
}
