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

/*
type MessageTextJSON struct {
	MessageText *string `json:"messageText"`
}

type MessageImageJSON struct {
	MessageImageURI  *string `json:"messageImageUri"`
	MessageImageType *string `json:"messageImageType"`
}

type MessageVideoJSON struct {
	MessageVideoURI      *string `json:"messageVideoUri"`
	MessageVideoEncoding *string `json:"messageVideoEncoding"`
	MessageVideoLength   *int    `json:"messageVideoLength"`
}
*/

type MessageContentJSON struct {
	MessageContentID *int    `json:"messageContentId,omitempty"`
	MessageType      *string `json:"type" binding:"required"`
	// -- Text message fields
	MessageText *string `json:"text,omitempty"`
	// -- Image fields
	MessageImageURI  *string `json:"messageImageUri,omitempty"`
	MessageImageType *string `json:"messageImageType,omitempty"`
	// -- Video fields
	MessageVideoURI      *string `json:"messageVideoUri,omitempty"`
	MessageVideoEncoding *string `json:"messageVideoEncoding,omitempty"`
	MessageVideoLength   *int    `json:"messageVideoLength,omitempty"`
}
