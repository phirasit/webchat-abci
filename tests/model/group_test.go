package model

import (
	"github.com/tendermint/tendermint/libs/db"
	"testing"
	"webchatABCI/model"
)

func TestGroup(t *testing.T) {
	d := db.NewDB("tmpDB", db.MemDBBackend, "/tmp")
	id := model.UserID("phirasit")
	gid := model.GroupID("group1")

	group := model.CreateNewGroup(&d, id, gid)

	if group.GetID() != gid {
		t.Errorf("GetName of a new group return %s, expect %s", group.GetID(), gid)
	}

	if group.GetAdmin() != id {
		t.Errorf("GetAdmin of a new group return %s, expect %s", group.GetAdmin(), id)
	}
}
