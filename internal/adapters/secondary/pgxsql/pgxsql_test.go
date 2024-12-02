package pgxsql_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/mkadit/go-toybox/internal/utils/dbutil"
)

var (
	dbTester, _ = dbutil.NewDBTester(false)
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	d, err := dbTester.GetConnection(ctx, "init")
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	err = d.InsertDataTest()
	if err != nil {
		log.Fatal("failed to insert test data: %w", err)
	}
	dbTester.SetupDataTest()
	_ = d

	res := m.Run()

	os.Exit(res)
}
