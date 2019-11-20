package client

import (
	"github.com/tendermint/tendermint/libs/db"
	"webchatABCI/model"
	"webchatABCI/query"
	"webchatABCI/txn"
)

type State struct {
	db             db.DB
	txnHandler     *txn.TxnHandler
	queryHandler   *query.QueryHandler

	ChainID     string  `json:"chain_id"`
	Size        int64   `json:"size"`
	Height      int64   `json:"height"`
	AppHash     []byte  `json:"app_hash"`

	MessageTimestamp  model.Timestamp  `json:"message_timestamp"`
}

