package model

import (
	"encoding/json"
	"github.com/tendermint/tendermint/libs/db"
)

type UserID string
type Timestamp int32

type User struct {
	db        *db.DB

	ID        UserID     `json:"id"`
}

func CreateNewUser (db *db.DB, id UserID) *User {
	return &User{
		db: db,

		ID: id,
	}
}

func GetUser (db *db.DB, id UserID) *User {

	data := (*db).Get(generateDBid(userPrefix, []byte(id)))

	if data == nil {
		return nil
	}

	user := &User{ db: db }
	err := json.Unmarshal(data, &user)

	if err == nil {
		return nil
	}

	user.db = db

	return user
}

func (user User) GetID() UserID { return user.ID }

func (user *User) Save() {

	marshal, err := json.Marshal(user)

	if err == nil {
		(*user.db).SetSync(generateDBid(userPrefix, []byte(user.ID)), marshal)
	}
}