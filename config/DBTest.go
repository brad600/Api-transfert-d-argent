package sqlstore

import (
	"POD_TRANSFERT_ARGENT_API/store"
	"os"
	"testing"
)

func TestSqlStore(t *testing.T) {
	dbPath := "/tmp/tets.db"
	s, err := New(dbPath, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	defer os.RemoveAll(dbPath)

	store.TestStore(s, t)
	store.TestStoreConcurrentTransfer(s, t)
}
