package txn

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type joinGroupHandler struct {
	TxnHandlerInterface
}

type joinGroupType struct {
	User    model.UserID     `json:"user"`
	Group   model.GroupID    `json:"group"`
}

func (joinGroupHandler) Check(db db.DB, data interface{}) types.ResponseCheckTx {

	_byte, _ := json.Marshal(data)

	txn := &joinGroupType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}

	if model.GetGroup(&db, txn.Group) == nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnauthorized,
			Log:  fmt.Sprintf(`The group "%v" does not exist`, txn.Group),
		}
	}

	// check whether the user is already in the group
	if model.IsUserInGroup(&db, txn.User, txn.Group) {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnauthorized,
			Log:  fmt.Sprintf(`The user "%v" is already in the group "%v"`, txn.User, txn.Group),
		}
	}

	return types.ResponseCheckTx{
		Code: code.CodeTypeOK,
		Log: fmt.Sprintf(`The user "%v" has joined the group "%v"`, txn.User, txn.Group),
	}
}

func (joinGroupHandler) Deliver(db db.DB, data interface{}) types.ResponseDeliverTx {

	_byte, _ := json.Marshal(data)

	txn := &joinGroupType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	// add the user into the desired group
	model.CreateNewInGroup(&db, txn.User, txn.Group).Save()

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Log:  fmt.Sprintf(`User %v has joined group "%v"`, txn.User, txn.Group),
	}
}

