package txn

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"

	"webchatABCI/model"
)

type readMessageHandler struct {
	TxnHandlerInterface
}

type readMessageType struct {
	User       model.UserID      `json:"user"`
	Group      model.GroupID     `json:"group"`
	Timestamp  model.Timestamp   `json:"timestamp"`
}

func (readMessageHandler) Check(db db.DB, data interface{}) types.ResponseCheckTx {

	_byte, _ := json.Marshal(data)

	txn := &readMessageType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}

	// check whether the user is in the group
	if !model.IsUserInGroup(&db, txn.User, txn.Group) {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnauthorized,
			Log:  fmt.Sprintf(`The user "%v" is not in the group "%v"`, txn.User, txn.Group),
		}
	}

	// check for the lastRead timestamp not to be lower than the current one
	timestamp := model.GetInGroup(&db, txn.User, txn.Group).GetTimestamp()

	if txn.Timestamp < timestamp {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnauthorized,
			Log:  "The new last read timestamp must not be lower than the current one",
		}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (readMessageHandler) Deliver(db db.DB, data interface{}) types.ResponseDeliverTx {

	_byte, _ := json.Marshal(data)

	txn := &readMessageType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	// update the lastRead timestamp
	inGroup := model.GetInGroup(&db, txn.User, txn.Group)

	if inGroup == nil {
		inGroup = model.CreateNewInGroup(&db, txn.User, txn.Group)
	}

	inGroup.UpdateLastRead(txn.Timestamp).Save()

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Log:  "The last read message is updated",
	}
}

