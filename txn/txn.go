package txn

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type TxnHandlerInterface interface {
	Check   (db.DB, interface{})   types.ResponseCheckTx
	Deliver (db.DB, interface{})   types.ResponseDeliverTx
}

type TxnHandler struct {
	CreateGroupHandler   createGroupHandler
	JoinGroupHandler     joinGroupHandler
	LeaveGroupHandler    leaveGroupHandler
	SendMessageHandler   sendMessageHandler
	ReadMessageHandler   readMessageHandler
}

const (
	// group
	CreateGroup      = "create_group"
	JoinGroup        = "join_group"
	LeaveGroup       = "leave_group"

	// message
	SendMessage      = "send_message"
	ReadMessage      = "read_message"
)

func CreateNewTxnHandler(timestamp *model.Timestamp) *TxnHandler {
	handler := &TxnHandler{
		CreateGroupHandler:  createGroupHandler{},
		JoinGroupHandler:    joinGroupHandler{},
		LeaveGroupHandler:   leaveGroupHandler{},
		SendMessageHandler:  sendMessageHandler{},
		ReadMessageHandler:  readMessageHandler{},
	}

	return handler.Init(timestamp)
}

func (txn TxnHandler) GetTxnHandler(t string) TxnHandlerInterface {

	switch t {

	case CreateGroup: return txn.CreateGroupHandler
	case JoinGroup:   return txn.JoinGroupHandler
	case LeaveGroup:  return txn.LeaveGroupHandler

	case SendMessage: return txn.SendMessageHandler
	case ReadMessage: return txn.ReadMessageHandler
	}

	return nil
}

func (txn *TxnHandler) Init(timestamp *model.Timestamp) *TxnHandler {
	txn.SendMessageHandler.Init(timestamp)
	return txn
}
