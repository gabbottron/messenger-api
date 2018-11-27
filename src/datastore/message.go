package datastore

import (
	//"database/sql"
	//"errors"
	"fmt"
	"github.com/gabbottron/messenger-api/src/datatypes"
	"log"
	"strconv"
)

const (
	InsertMessageQuery string = `INSERT INTO message 
		(senderid, recipientid, messagecontenttype) 
		VALUES($1, $2, $3)
		RETURNING messageid, senderid, recipientid, messagecontenttype, createdate`

	InsertMessageTextQuery string = `INSERT INTO message_text 
		(messageid, messagetext) 
		VALUES($1, $2)
		RETURNING messagetext`

	// In case we want to see a full conversation...
	SelectMessagesQuery string = `SELECT m.messageid, mt.messagetext, mi.imageuri 
		FROM message AS m 
		LEFT JOIN message_text AS mt ON mt.messageid = m.messageid 
		LEFT JOIN message_image AS mi ON mi.messageid = m.messageid 
		WHERE ( m.senderid = $1 AND m.recipientid = $2 ) OR ( m.senderid = $2 AND m.recipientid = $1) 
		ORDER BY m.createdate ASC;`

	SelectMessagesToRecipientQuery string = `SELECT m.messageid, m.createdate, 
		m.senderid, m.recipientid, m.messagecontenttype, mt.messagetext, mi.imageuri 
		FROM message AS m 
		LEFT JOIN message_text AS mt ON mt.messageid = m.messageid 
		LEFT JOIN message_image AS mi ON mi.messageid = m.messageid 
		WHERE m.senderid = $1 AND m.recipientid = $2 AND m.messageid >= $3
		ORDER BY m.createdate ASC
		LIMIT $4;`
)

/*

type MessageJSON struct {
	MessageID      int                 `json:"id"`
	Timestamp      *time.Time          `json:"timestamp,omitempty"`
	SenderID       int                 `json:"sender" binding:"required"`
	RecipientID    int                 `json:"recipient" binding:"required"`
	MessageContent *MessageContentJSON `json:"content" binding:"required"`
}
*/

func GetMessages(sender int, recipient int, start int, limit int) ([]datatypes.MessageJSON, error) {
	results := make([]datatypes.MessageJSON, 0)

	rows, err := db.Query(SelectMessagesToRecipientQuery, sender, recipient, start, limit)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var obj datatypes.MessageJSON
		obj.MessageContent = new(datatypes.MessageContentJSON)
		err = rows.Scan(&obj.MessageID, &obj.Timestamp, &obj.SenderID, &obj.RecipientID,
			&obj.MessageContent.MessageType, &obj.MessageContent.MessageText,
			&obj.MessageContent.MessageImageURI)

		// If the scan was successful, load the row
		if err == nil {
			results = append(results, obj)
			count++
			//log.Println(temp_row)
		}
	}

	// Show the count of rows
	log.Println("Rows returned: " + strconv.Itoa(count))

	if err = rows.Err(); err != nil {
		// Abnormal termination of the rows loop
		// close should be called automatically in this case
		log.Println(err)
	}

	return results, nil
}

func InsertMessageRecord(obj *datatypes.MessageJSON) error {
	// TODO: Should wrap these queries in a transaction probably, in case
	//       the content type insert fails!

	log.Println("Insert message record...")

	// attempt to insert the new message record
	err := db.QueryRow(InsertMessageQuery, obj.SenderID, obj.RecipientID, obj.MessageContent.MessageType).Scan(
		&obj.MessageID, &obj.SenderID, &obj.RecipientID, &obj.MessageContent.MessageType, &obj.Timestamp)

	if err != nil {
		log.Printf("Insert message: QueryRow.Scan() -> %s", err.Error())

		if IsConstraintViolation(err) {
			is_dup, msg := IsDuplicateKeyViolation(err)
			if is_dup {
				return &DatastoreError{msg, ErrorConstraintViolation}
			}
			// Handle the db constraint error
			return &DatastoreError{ErrorConstraintViolationString, ErrorConstraintViolation}
		}
		return err
	}

	// Now insert the correct content type
	switch *obj.MessageContent.MessageType {
	case "text":
		fmt.Println("text")
		err = db.QueryRow(InsertMessageTextQuery, obj.MessageID, obj.MessageContent.MessageText).Scan(&obj.MessageContent.MessageText)

		if err != nil {
			log.Printf("Insert message text: QueryRow.Scan() -> %s", err.Error())

			if IsConstraintViolation(err) {
				is_dup, msg := IsDuplicateKeyViolation(err)
				if is_dup {
					return &DatastoreError{msg, ErrorConstraintViolation}
				}
				// Handle the db constraint error
				return &DatastoreError{ErrorConstraintViolationString, ErrorConstraintViolation}
			}
			return err
		}
	case "image":
		fmt.Println("image")
	case "video":
		fmt.Println("video")
	default:
		fmt.Println("unknown type!")
	}

	return nil
}
