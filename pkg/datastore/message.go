package datastore

import (
	//"database/sql"
	"errors"
	"github.com/gabbottron/messenger-api/pkg/datatypes"
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

	InsertMessageImageQuery string = `INSERT INTO message_image 
		(messageid, imageuri, imagetype) 
		VALUES($1, $2, $3)
		RETURNING imageuri, imagetype`

	InsertMessageVideoQuery string = `INSERT INTO message_video 
		(messageid, videouri, videoencoding, videolength) 
		VALUES($1, $2, $3, $4)
		RETURNING videouri, videoencoding, videolength`

	// In case we want to see a full conversation...
	SelectMessagesQuery string = `SELECT m.messageid, mt.messagetext, mi.imageuri, mi.imagetype
		FROM message AS m 
		LEFT JOIN message_text AS mt ON mt.messageid = m.messageid 
		LEFT JOIN message_image AS mi ON mi.messageid = m.messageid 
		WHERE ( m.senderid = $1 AND m.recipientid = $2 ) OR ( m.senderid = $2 AND m.recipientid = $1) 
		ORDER BY m.createdate ASC;`

	SelectMessagesToRecipientQuery string = `SELECT m.messageid, m.createdate, 
		m.senderid, m.recipientid, m.messagecontenttype, mt.messagetext, mi.imageuri, 
		mi.imagetype, mv.videouri, mv.videoencoding, mv.videolength   
		FROM message AS m 
		LEFT JOIN message_text AS mt ON mt.messageid = m.messageid 
		LEFT JOIN message_image AS mi ON mi.messageid = m.messageid
		LEFT JOIN message_video AS mv ON mv.messageid = m.messageid 
		WHERE m.senderid = $1 AND m.recipientid = $2 AND m.messageid >= $3
		ORDER BY m.createdate ASC
		LIMIT $4;`
)

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
			&obj.MessageContent.MessageImageURI, &obj.MessageContent.MessageImageType,
			&obj.MessageContent.MessageVideoURI, &obj.MessageContent.MessageVideoEncoding,
			&obj.MessageContent.MessageVideoLength)

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
		//log.Println("Content type: text")
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
		//log.Println("Content type: image")
		err = db.QueryRow(InsertMessageImageQuery, obj.MessageID, obj.MessageContent.MessageImageURI, obj.MessageContent.MessageImageType).Scan(&obj.MessageContent.MessageImageURI, &obj.MessageContent.MessageImageType)

		if err != nil {
			log.Printf("Insert message image: QueryRow.Scan() -> %s", err.Error())

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
	case "video":
		//log.Println("Content type: video")
		err = db.QueryRow(InsertMessageVideoQuery, obj.MessageID, obj.MessageContent.MessageVideoURI, obj.MessageContent.MessageVideoEncoding, obj.MessageContent.MessageVideoLength).Scan(&obj.MessageContent.MessageVideoURI, &obj.MessageContent.MessageVideoEncoding, &obj.MessageContent.MessageVideoLength)

		if err != nil {
			log.Printf("Insert message image: QueryRow.Scan() -> %s", err.Error())

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
	default:
		log.Println("unknown type!")
		return errors.New("Unknown message type!")
	}

	return nil
}
