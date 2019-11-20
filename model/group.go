package model

import (
	"encoding/json"
	"github.com/tendermint/tendermint/libs/db"
)

type GroupID string

type Group struct {
	db     *db.DB

	ID           GroupID    `json:"id"`
	Admin        UserID     `json:"admin"`
	LastMessage  MessageID  `json:"last_message"`
}

func CreateNewGroup (db *db.DB, user UserID, id GroupID) *Group {
	return &Group {
		db:           db,

		ID:           id,
		Admin:        user,
		LastMessage:  MessageNil,
	}
}

func GetGroup(db *db.DB, groupID GroupID) *Group {
	groupData := (*db).Get(generateDBid(groupPrefix, []byte(groupID)))

	if groupData == nil {
		return nil
	}

	group := &Group { db: db }
	err := json.Unmarshal(groupData, group)

	if err != nil {
		return nil
	}

	return group
}

func GetLastMessageID(db *db.DB, groupID GroupID) MessageID {
	group := GetGroup(db, groupID)
	if group == nil {
		return MessageNil
	} else {
		return group.GetLastMessage()
	}
}

func (group Group) GetID() GroupID { return group.ID }
func (group Group) GetAdmin() UserID { return group.Admin }
func (group Group) GetLastMessage() MessageID { return group.LastMessage }

func (group *Group) Save() {

	marshal, err := json.Marshal(group)

	if err == nil {
		(*group.db).SetSync(generateDBid(groupPrefix, []byte(group.ID)), marshal)
	}
}

func (group *Group) SetLastMessage(messageID MessageID) *Group {
	group.LastMessage = messageID
	return group
}

func IsGroupExists(db *db.DB, group GroupID) bool {
	return (*db).Has(generateDBid(groupPrefix, []byte(group)))
}
