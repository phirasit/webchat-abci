package query

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type getUnreadMessageHandler struct {
	QueryHandlerInterface

	currentMessageTimestamp  *model.Timestamp
}

type queryUnreadMessageRequest struct {
	User       model.UserID   `json:"user"`
	Group      model.GroupID  `json:"group"`
}

type queryUnreadMessageResponse struct {
	NumMessage   int              `json:"num_messages"`
	Messages     []model.Message  `json:"messages"`
	Timestamp    model.Timestamp  `json:"timestamp"`
}

func (handler getUnreadMessageHandler) Query(db db.DB, reqQuery types.RequestQuery) types.ResponseQuery {

	query := &queryUnreadMessageRequest{}
	err := json.Unmarshal(reqQuery.Data, &query)

	if err != nil {
		return types.ResponseQuery {
			Code: code.CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}

	// check whether a user exists
	// check whether a group exists
	// check whether a user is in the group
	if !model.IsUserInGroup(&db, query.User, query.Group) {
		return types.ResponseQuery {
			Code: code.CodeTypeUnauthorized,
			Log:  fmt.Sprintf(`user "%v" is not in the group "%v"`, query.User, query.Group),
		}
	}

	lastReadTimestamp  := model.GetInGroup(&db, query.User, query.Group).GetTimestamp()
	groupLastMessageID := model.GetGroup(&db, query.Group).GetLastMessage()
	lastMessage        := model.GetMessage(&db, groupLastMessageID)

	// read all getUnread message by
	response := queryUnreadMessageResponse{
		Timestamp: *handler.currentMessageTimestamp,
	}

	for lastMessage != nil && lastMessage.GetTimestamp() > lastReadTimestamp {
		response.Messages = append(response.Messages, *lastMessage)
		lastMessage = model.GetMessage(&db, lastMessage.GetPrevMessageID())
	}

	response.NumMessage = len(response.Messages)

	// reverse the message
	for i, j := 0, response.NumMessage-1; i < j; i, j = i+1, j-1 {
		response.Messages[i], response.Messages[j] = response.Messages[j], response.Messages[i]
	}

	resData, err := json.Marshal(response)

	if err != nil {
		return types.ResponseQuery {
			Code: code.CodeTypeUnknownError,
			Log:  err.Error(),
		}
	}

	return types.ResponseQuery {
		Code:   code.CodeTypeOK,
		Index:  -1, // TODO change to block height
		Key:    reqQuery.Data,
		Value:  resData,
		Log:    "OK",
	}
}

func (handler *getUnreadMessageHandler) SetMessageTimestamp(timestamp *model.Timestamp) {
	handler.currentMessageTimestamp = timestamp
}
