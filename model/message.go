package model

import (
	"encoding/json"
	"github.com/tendermint/tendermint/libs/db"
	"strconv"
)

type MessageID string
type MessageType string

func (id MessageID) IsNull() bool { return id == MessageNil }

type Message struct {
	db         *db.DB

	ID         MessageID     `json:"id"`
	Group      GroupID       `json:"group"`
	User       UserID        `json:"user"`
	Message    MessageType   `json:"message"`
	Time       string        `json:"time"`
	Timestamp  Timestamp     `json:"timestamp"`

	PrevID     MessageID     `json:"prevID"`
}

var MessageNil = MessageID("")

func GenerateNewMessageID(seed int) MessageID {
	return MessageID(string(seed))
}

func CreateNewMessage(db *db.DB, user UserID, group GroupID, message MessageType, time string) *Message {
	return &Message {
		db:          db,

		ID:          "", // TODO find appropriate id generator
		Group:       group,
		User:        user,
		Message:     message,
		Time:        time,
		Timestamp:   0,

		PrevID:      MessageNil,
	}
}

func (message *Message) GenerateID(seed Timestamp) *Message {
	message.ID = MessageID(strconv.Itoa(int(seed)))
	return message
}

func GetMessage(db *db.DB, msgID MessageID) *Message {

	data := (*db).Get(generateDBid(messagePrefix, []byte(msgID)))

	if data == nil {
		return nil
	}

	message := &Message{ db: db }
	err := json.Unmarshal(data, &message)

	if err != nil {
		return nil
	}

	message.db = db

	return message
}

func (message Message) GetMessage() MessageType { return message.Message }
func (message Message) GetUser() UserID { return message.User }
func (message Message) GetID() MessageID { return message.ID }
func (message Message) GetGroup() GroupID { return message.Group }
func (message Message) GetTime() string { return message.Time }
func (message Message) GetTimestamp() Timestamp { return message.Timestamp }
func (message Message) GetPrevMessageID() MessageID { return message.PrevID }

func (message *Message) Save() {

	marshal, err := json.Marshal(message)

	if err == nil {
		(*message.db).SetSync(generateDBid(messagePrefix, []byte(message.ID)), marshal)
	}
}

func (message *Message) SetPrevID(messageID MessageID) *Message {
	message.PrevID = messageID
	return message
}

func (message *Message) SetTimestamp(timestamp Timestamp) *Message {
	message.Timestamp = timestamp
	return message
}

func (message Message) Validate() bool {
	return IsUserInGroup(message.db, message.User, message.Group)
}

