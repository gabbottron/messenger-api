package datastore

import (
	"database/sql"
	"errors"
	"github.com/gabbottron/messenger-api/src/datatypes"
	"log"
)

const (
	InsertUserQuery string = `INSERT INTO msgr_user 
		(username, password) 
		VALUES($1, $2)
		RETURNING msgruserid, username`

	SelectUserQuery string = `SELECT msgruserid, username 
		FROM msgr_user WHERE msgruserid = $1`

	SelectUserByUsernameQuery string = `SELECT msgruserid, username 
		FROM msgr_user WHERE username = $1`
)

func InsertUserRecord(obj *datatypes.MsgrUserJSON) error {
	// attempt to insert the new user record
	err := db.QueryRow(InsertUserQuery, obj.Username, obj.Password).Scan(
		&obj.MsgrUserID, &obj.Username)

	if err != nil {
		log.Printf("Error inserting user record -> %s", err.Error())
		return err
	}

	return nil
}

func AuthenticateUser(username string, password string) (*datatypes.User, error) {
	// The userID from the database
	var id int
	// The hashed pw from the database
	var hashed_pass string

	err := db.QueryRow(`SELECT msgruserid, password FROM msgr_user WHERE username = $1`, username).Scan(&id, &hashed_pass)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No matching user!")
		}

		log.Printf("Error authenticating user -> %s", err.Error())
		return nil, errors.New("Unknown DB error!")
	}

	// Now compare the password provided
	if !ComparePasswords(hashed_pass, password) {
		log.Printf("Error authenticating user -> Passwords didn't match!")
		return nil, errors.New("Incorrect username/password combo!")
	}

	return &datatypes.User{
		UserID: id,
	}, nil
}
