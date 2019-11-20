package benchmarks

import (
	"github.com/tendermint/tendermint/libs/db"
	"testing"
	"webchatABCI/model"
)

func BenchmarkUserSave(b *testing.B) {

	d := db.NewDB("benchmark", db.MemDBBackend, "cache")

	for i := 0; i < b.N; i += 1 {
		user := model.CreateNewUser(&d, []byte("user" + string(i)))
		user.Save()
	}
}
