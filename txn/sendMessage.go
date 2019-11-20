package txn

import (
	"encoding/json"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type sendMessageHandler struct {
	TxnHandlerInterface

	messageTimestamp   *model.Timestamp
}

type sendMessageRequest struct {
	User        model.UserID      `json:"user"`
	Group       model.GroupID     `json:"group"`
	Message     model.MessageType `json:"message"`
	Time        string            `json:"time"`
}

func (sendMessageHandler) Check(db db.DB, data interface{}) types.ResponseCheckTx {

	_byte, _ := json.Marshal(data)

	txn := &sendMessageRequest{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}

	message := model.CreateNewMessage(&db, txn.User, txn.Group, txn.Message, txn.Time)

	if message.Validate() {
		return types.ResponseCheckTx{Code: code.CodeTypeOK}
	} else {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnauthorized,
			Log:  "This message is invalid",
		}
	}

}

func (handler sendMessageHandler) Deliver(db db.DB, data interface{}) types.ResponseDeliverTx {

	_byte, _ := json.Marshal(data)

	txn := &sendMessageRequest{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	// create a new message
	message := model.CreateNewMessage(&db, txn.User, txn.Group, txn.Message, txn.Time)
	message.SetTimestamp(*handler.messageTimestamp)

	if message.Validate() {

		// increase the message counter
		*handler.messageTimestamp += 1

		// generate MessageID
		message.GenerateID(*handler.messageTimestamp)

		// set previous messageID
		message.SetPrevID(model.GetLastMessageID(&db, message.GetGroup()))

		// change the last message of the group
		model.GetGroup(&db, message.GetGroup()).SetLastMessage(message.GetID()).Save()

		// save the message
		message.Save()

		return types.ResponseDeliverTx {
			Code: code.CodeTypeOK,
			Log:  "The message is delivered",
			Tags: []common.KVPair {
				{
					Key:   []byte("message.Group"),
					Value: []byte(message.GetGroup()),
				},
			},
		}
	} else {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeUnauthorized,
			Log:  "This message is not authorized",
		}
	}

}

func (handler *sendMessageHandler) Init(time *model.Timestamp) {
	handler.messageTimestamp = time
}
