package benchmarks

import (
	"github.com/tendermint/tendermint/libs/db"
	"testing"
	"webchatABCI/model"
)


func BenchmarkMessageSave(b *testing.B) {

	d := db.NewDB("benchmark", db.MemDBBackend, "cache")

	for i := 0; i < b.N; i += 1 {
		user := model.CreateNewMessage(&d, []byte("user" + string(i)), []byte("group1"), []byte("hello world"))
		user.Save()
	}
}
