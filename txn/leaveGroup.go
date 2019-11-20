package txn

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type leaveGroupHandler struct {
	TxnHandlerInterface
}

type leaveGroupType struct {
	User    model.UserID     `json:"user"`
	Group   model.GroupID    `json:"group"`
}

func (leaveGroupHandler) Check(db db.DB, data interface{}) types.ResponseCheckTx {

	_byte, _ := json.Marshal(data)

	txn := &leaveGroupType{}
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

	return types.ResponseCheckTx{
		Code: code.CodeTypeOK,
		Log: fmt.Sprintf(`The user "%v" has left the group "%v"`, txn.User, txn.Group),
	}
}

func (leaveGroupHandler) Deliver(db db.DB, data interface{}) types.ResponseDeliverTx {

	_byte, _ := json.Marshal(data)

	txn := &leaveGroupType{}
	err := json.Unmarshal(_byte, &txn)

	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	// remove the user from the group
	inGroup := model.GetInGroup(&db, txn.User, txn.Group)
	if inGroup != nil {
		inGroup.Delete()
	}

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Log:  fmt.Sprintf(`User %v has left group "%v"`, txn.User, txn.Group),
	}
}

