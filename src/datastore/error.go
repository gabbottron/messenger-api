package datastore

import (
	"fmt"
	"regexp"
	"strings"
)

// Error codes for this package
// This way you can use the default string if you like, or pass your own!
const (
	ErrorNoFieldsToUpdate                   = iota
	ErrorNoFieldsToUpdateString      string = "No fields to update!"
	ErrorNotAuthorized                      = iota
	ErrorNotAuthorizedString         string = "You are not authorized!"
	ErrorForbidden                          = iota
	ErrorForbiddenString             string = "You don't have permission to create/modify this resource!"
	ErrorMissingRequiredFields              = iota
	ErrorMissingRequiredFieldsString string = "The request was missing required fields!"
	ErrorConstraintViolation                = iota
	ErrorConstraintViolationString   string = "The data in the request violated a database constraint!"
	ErrorRecordNotFound                     = iota
	ErrorRecordNotFoundString               = "Record not found!"

	ErrorUnknownDatabaseError              = iota
	ErrorUnknownDatabaseErrorString string = "There was a database error performing the request!"

	// The postgres error text matching strings
	PostgresViolatesConstraintMatchString string = "violates"
	// This will match the duplicate key error test from postgres and return the constraint name
	PostgresViolatesConstraintDuplicateKeyMatchString string = "^.*duplicate +key.*violates +unique.*constraint +\"(.*?)\""

	// pq: null value in column "username" violates not-null constraint
	PostgresViolatesConstraintNullValueMatchString string = "^.*null +value.*in +column +\"(.*?)\" +violates +not-null +constraint"
)

type DatastoreError struct {
	Text string
	Code int
}

func (e *DatastoreError) Error() string {
	return fmt.Sprintf("Code: %d: %s", e.Code, e.Text)
}

func IsDatastoreError(err error) bool {
	if _, ok := err.(*DatastoreError); ok {
		return true
	}
	return false
}

// Postgres error helpers -------------

// Is the error a constraint violation?
func IsConstraintViolation(err error) bool {
	if strings.Contains(err.Error(), PostgresViolatesConstraintMatchString) {
		return true
	}
	return false
}

// Check for null value violation (special)
// pq: null value in column "username" violates not-null constraint
func IsSpecialNullViolationCase(err error) (bool, string) {
	re := regexp.MustCompile(PostgresViolatesConstraintNullValueMatchString)

	match := re.FindStringSubmatch(err.Error())

	if len(match) < 2 {
		return false, ""
	} else {
		switch match[1] {
		case "username":
			return true, "You must provide a username!"
		default:
			return true, "The record violated a not-null constraint: " + match[1]
		}
	}
}

// Check for duplicate key violation
// ex: pq: duplicate key value violates unique constraint "msgr_user_username_key"
func IsDuplicateKeyViolation(err error) (bool, string) {
	re := regexp.MustCompile(PostgresViolatesConstraintDuplicateKeyMatchString)

	match := re.FindStringSubmatch(err.Error())

	if len(match) < 2 {
		return false, ""
	} else {
		switch match[1] {
		case "msgr_user_username_key":
			return true, "Someone else has already chosen that username!"
		default:
			return true, "The record violated a unique constraint: " + match[1]
		}
	}
}
