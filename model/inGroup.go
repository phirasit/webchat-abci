package model

import (
	"encoding/json"
	"github.com/tendermint/tendermint/libs/db"
)

type InGroup struct {
	db        *db.DB
	User      UserID      `json:"user"`
	Group     GroupID     `json:"group"`

	LastRead  Timestamp   `json:"last_read"`
}

func CreateNewInGroup (db *db.DB, user UserID, group GroupID) *InGroup {
	return &InGroup{
		db:       db,
		User:     user,
		Group:    group,
		LastRead: 0,
	}
}

func GetInGroup(db *db.DB, user UserID, group GroupID) *InGroup {

	data := (*db).Get(generateDBid(inGroupPrefix, getInGroupID(user, group)))

	if data == nil {
		return nil
	}

	inGroup := &InGroup { db: db }
	err := json.Unmarshal(data, &inGroup)

	if err != nil {
		return nil
	}

	return inGroup
}

func getInGroupID(user UserID, group GroupID) []byte {
	return []byte(string(user) + " " + string(group))
}

func IsUserInGroup(db *db.DB, user UserID, group GroupID) bool {
	return (*db).Has(generateDBid(inGroupPrefix, getInGroupID(user, group)))
}

func (inGroup *InGroup) UpdateLastRead(timestamp Timestamp) *InGroup {
	inGroup.LastRead = timestamp
	return inGroup
}

func (inGroup InGroup) GetTimestamp() Timestamp {
	return inGroup.LastRead
}

func (inGroup InGroup) Save() {

	data, err := json.Marshal(inGroup)

	if err == nil {
		(*inGroup.db).Set(generateDBid(inGroupPrefix, getInGroupID(inGroup.User, inGroup.Group)), data)
	}
}

func (inGroup InGroup) Delete() {
	(*inGroup.db).Delete(generateDBid(inGroupPrefix, getInGroupID(inGroup.User, inGroup.Group)))
}

