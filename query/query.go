package query

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
)

type QueryHandlerInterface interface {
	Query (db.DB, types.RequestQuery)   types.ResponseQuery
}

type QueryHandler struct {
	GetUnreadMessageHandler  getUnreadMessageHandler
	GetMessageHandler        getMessageHandler
}
const (
	// message
	GetUnreadMessage = "get_unread_message"
	GetMessage       = "get_message"
)

func CreateQueryHandler(timestamp *model.Timestamp) *QueryHandler {
	queryHandler := &QueryHandler {
		GetUnreadMessageHandler: getUnreadMessageHandler{},
		GetMessageHandler:       getMessageHandler{},
	}
	return queryHandler.Init(timestamp)
}

func (queryHandler QueryHandler) GetQueryHandler(t string) QueryHandlerInterface {

	switch t {
	case GetUnreadMessage: return queryHandler.GetUnreadMessageHandler
	case GetMessage:       return queryHandler.GetMessageHandler
	}

	return nil
}

func (queryHandler *QueryHandler) Init(timestamp *model.Timestamp) *QueryHandler {
	queryHandler.GetUnreadMessageHandler.SetMessageTimestamp(timestamp)
	return queryHandler
}
