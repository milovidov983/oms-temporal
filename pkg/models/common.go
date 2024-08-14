package models

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
	CommentRequest struct {
		Comment string `json:"comment"`
	}
	ReasonRequest struct {
		Reason string `json:"reason"`
	}
)
