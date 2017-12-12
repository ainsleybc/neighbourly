package db

import (
	"testing"

	r "github.com/dancannon/gorethink"
)

func TestCleanUp(t *testing.T) {

	t.Parallel()

	session, _ := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})

	dbName := "cleanUp"

	r.DBCreate(dbName).RunWrite(session)
	err := CleanUp(dbName)

	if err != nil {
		t.Errorf("\n%v", err)
	}

	resp, _ := r.DBList().Run(session)

	var row string
	for resp.Next(&row) {
		if row == dbName {
			t.Errorf("database %v not dropped", dbName)
		}
	}

	session.Close()

}

func TestSetup(t *testing.T) {

	t.Parallel()

	session, _ := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})

	var row string
	dbName := "setup"
	tables := []string{
		"users",
		"feeds",
		"addresses",
		"feedAddresses",
		"posts",
	}

	CleanUp(dbName)
	err := Setup(dbName)

	if err != nil {
		t.Errorf("\n%v", err)
	}

	resp, _ := r.DB(dbName).TableList().Run(session)

	for resp.Next(&row) {
		if !contains(tables, row) {
			t.Errorf("database %v not dropped", dbName)
		}
	}

	CleanUp(dbName)
	session.Close()

}

func contains(arr []string, table string) bool {
	for _, el := range arr {
		if string(el) == table {
			return true
		}
	}
	return false
}
