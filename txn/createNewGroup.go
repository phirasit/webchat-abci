package txn

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type createGroupHandler struct {
	TxnHandlerInterface
}

type CreateGroupType struct {
	User     model.UserID    `json:"user"`
	Group    model.GroupID   `json:"group"`
}

func (createGroupHandler) Check(db db.DB, data interface{}) types.ResponseCheckTx {

	_byte, _ := json.Marshal(data)

	txn := &CreateGroupType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}

	// check whether the group exists
	if model.IsGroupExists(&db, txn.Group) {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnauthorized,
			Log:  fmt.Sprintf(`group "%v" is already exists`, txn.Group),
		}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (createGroupHandler) Deliver(db db.DB, data interface{}) types.ResponseDeliverTx {

	_byte, _ := json.Marshal(data)

	txn := &CreateGroupType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	// add admin to the group
	model.CreateNewInGroup(&db, txn.User, txn.Group).Save()

	// create new group
	model.CreateNewGroup(&db, txn.User, txn.Group).Save()

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Log:  fmt.Sprintf(`New group "%v" is created`, txn.Group),
	}
}

