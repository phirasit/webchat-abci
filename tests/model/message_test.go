package model

import (
	"github.com/tendermint/tendermint/libs/db"
	"testing"
	"webchatABCI/model"
)

func TestMessage(t *testing.T) {
	d := db.NewDB("tmpDB", db.MemDBBackend, "/tmp")
	id := model.UserID("phirasit")
	name := model.GroupID("message1")
	messageText := model.MessageType("message")

	message := model.CreateNewMessage(&d, id, name, messageText)

	if message.GetMessage() != messageText {
		t.Errorf("A new message return %s, expect %s", message.GetMessage(), message)
	}

	if message.GetTimestamp() != 0 {
		t.Errorf("A new message timestamp is %d, should be 0", message.GetTimestamp())
	}

	if message.GetUser() != id {
		t.Errorf("A new message user return %s, expect %s", message.GetUser(), id)
	}
}
