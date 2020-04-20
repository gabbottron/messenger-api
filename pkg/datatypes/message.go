package datatypes

import (
	"time"
)

type MessageJSON struct {
	MessageID      int                 `json:"id"`
	Timestamp      *time.Time          `json:"timestamp,omitempty"`
	SenderID       int                 `json:"sender" binding:"required"`
	RecipientID    int                 `json:"recipient" binding:"required"`
	MessageContent *MessageContentJSON `json:"content" binding:"required"`
}

type MessageContentJSON struct {
	MessageContentID *int    `json:"messageContentId,omitempty"`
	MessageType      *string `json:"type" binding:"required"`
	// -- Text message fields
	MessageText *string `json:"text,omitempty"`
	// -- Image fields
	MessageImageURI  *string `json:"imageUri,omitempty"`
	MessageImageType *string `json:"imageType,omitempty"`
	// -- Video fields
	MessageVideoURI      *string `json:"videoUri,omitempty"`
	MessageVideoEncoding *string `json:"encoding,omitempty"`
	MessageVideoLength   *int    `json:"length,omitempty"`
}
