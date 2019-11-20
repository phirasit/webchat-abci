package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/version"
	"webchatABCI/query"
	"webchatABCI/txn"

	"github.com/tendermint/tendermint/libs/db"
)

var (
	stateKey                         = []byte("webchat")
	ProtocolVersion version.Protocol = 0x1
)

func loadState(db db.DB) State {
	stateBytes := db.Get(stateKey)
	var cState State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &cState)
		if err != nil {
			panic(err)
		}
	} else {
		cState = State{
			ChainID: "chain0",
			Size:    0,
			Height:  0,
			AppHash: []byte{},

			MessageTimestamp: 0,
		}
	}
	cState.db = db

	return cState
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

var _ types.Application = (*WebChatApplication)(nil)

type WebChatApplication struct {
	types.BaseApplication

	state State
}

func (app *WebChatApplication) GetState() *State {
	return &app.state
}

type Txn struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
	Nonce []byte      `json:"nonce"`
}

func (app *WebChatApplication) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{
		Data:       fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:    version.ABCIVersion,
		AppVersion: ProtocolVersion.Uint64(),
	}
}

func NewWebChatApplication() *WebChatApplication {
	webChat := loadState(db.NewMemDB())

	webChat.txnHandler = txn.CreateNewTxnHandler(&webChat.MessageTimestamp)
	webChat.queryHandler = query.CreateQueryHandler(&webChat.MessageTimestamp)

	return &WebChatApplication{state: webChat}
}

func (app *WebChatApplication) CheckTx(tx []byte) types.ResponseCheckTx {

	obj := &Txn{}
	err := json.Unmarshal(tx, obj)

	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}

	txnHandler := app.state.txnHandler.GetTxnHandler(obj.Type)

	if txnHandler == nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnknownError,
			Log:  fmt.Sprintf(`No Txn handler match "%v"`, obj.Type),
		}
	}

	return txnHandler.Check(app.state.db, obj.Data)
}

func (app *WebChatApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {

	obj := &Txn{}
	err := json.Unmarshal(tx, obj)

	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	txnHandler := app.state.txnHandler.GetTxnHandler(obj.Type)

	if txnHandler == nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeUnknownError,
			Log:  fmt.Sprintf(`No Txn handler match "%v"`, obj.Type),
		}
	}

	return txnHandler.Deliver(app.state.db, obj.Data)
}

func (app *WebChatApplication) Commit() types.ResponseCommit {

	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)

	app.state.AppHash = appHash
	app.state.Height += 1
	saveState(app.state)

	return types.ResponseCommit{Data: appHash}

}

func (app *WebChatApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {

	data := reqQuery.GetData()

	obj := &Txn{}
	err := json.Unmarshal(data, obj)

	if err != nil {
		return types.ResponseQuery{Code: code.CodeTypeEncodingError}
	}

	queryHandler := app.state.queryHandler.GetQueryHandler(obj.Type)

	if queryHandler == nil {
		return types.ResponseQuery {
			Code: code.CodeTypeUnknownError,
			Log:  fmt.Sprintf(`No Query handler match "%v"`, obj.Type),
		}
	}

	reqQuery.Data, _ = json.Marshal(obj.Data)

	return queryHandler.Query(app.state.db, reqQuery)
}
