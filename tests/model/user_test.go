package model

import (
	"github.com/tendermint/tendermint/libs/db"
	"testing"
	"webchatABCI/model"
)

func TestUser(t *testing.T) {
	d := db.NewDB("tmpDB", db.MemDBBackend, "/tmp")
	id := model.UserID("phirasit")

	user := model.CreateNewUser(&d, id)

	if user.GetID() != id {
		t.Errorf("GetID of a new user return %s, expect %s", user.GetID(), id)
	}
}
