package query

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type getMessageHandler struct {
	QueryHandlerInterface
}

type queryGetMessageRequest struct {
	User            model.UserID      `json:"user"`
	Group           model.GroupID     `json:"group"`
	LastMessageID   model.MessageID   `json:"last_message"`
	Limit           int32             `json:"limit"`
}

type queryGetMessageResponse struct {
	NumMessage      int               `json:"num_messages"`
	Messages        []model.Message   `json:"messages"`
	PrevMessage     model.MessageID   `json:"prev_message"`
	Timestamp       model.Timestamp   `json:"timestamp"`
}

func (getMessageHandler) Query(db db.DB, reqQuery types.RequestQuery) types.ResponseQuery {

	query := &queryGetMessageRequest{}
	err := json.Unmarshal(reqQuery.Data, &query)

	if err != nil {
		return types.ResponseQuery{
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

	var lastMessage *model.Message = nil

	if query.LastMessageID.IsNull() {
		groupLastMessageID := model.GetLastMessageID(&db, query.Group)
		lastMessage         = model.GetMessage(&db, groupLastMessageID)
	} else {
		lastMessage         = model.GetMessage(&db, query.LastMessageID)

		if lastMessage == nil {
			return types.ResponseQuery {
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf(`There is no message ID = %v`, query.LastMessageID),
			}
		}

		if lastMessage.GetGroup() != query.Group {
			return types.ResponseQuery {
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf(`The message %v is not in the group`, query.LastMessageID),
			}
		}
	}

	// get all unget message by
	response := queryGetMessageResponse{}

	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	for lastMessage != nil && query.Limit > 0 {
		response.Messages = append(response.Messages, *lastMessage)
		lastMessage = model.GetMessage(&db, lastMessage.GetPrevMessageID())

		query.Limit -= 1
	}

	response.NumMessage = len(response.Messages)

	// reverse the message
	for i, j := 0, response.NumMessage-1; i < j; i, j = i+1, j-1 {
		response.Messages[i], response.Messages[j] = response.Messages[j], response.Messages[i]
	}

	if lastMessage != nil {
		response.PrevMessage = lastMessage.GetID()
	} else {
		response.PrevMessage = model.MessageNil
	}

	resData, err := json.Marshal(response)

	if err != nil {
		return types.ResponseQuery{
			Code: code.CodeTypeUnknownError,
			Log:  err.Error(),
		}
	}

	return types.ResponseQuery{
		Code:   code.CodeTypeOK,
		Index:  -1, // TODO change to block height
		Key:    reqQuery.Data,
		Value:  resData,
		Log:    "OK",
	}
}
