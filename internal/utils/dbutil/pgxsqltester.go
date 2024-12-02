package dbutil

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/mkadit/go-toybox/internal/adapters/secondary/pgxsql"
	"github.com/mkadit/go-toybox/internal/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestData struct {
	UserData models.User
}

type DBTester struct {
	ConnUrl   string
	Parallel  bool
	SleepTime time.Duration
	TestData  TestData
}

func DeployTestContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("go-toybox"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("SuiseiKawaii"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	return container, err

}

func NewDBTester(parallel bool) (*DBTester, error) {
	ctx := context.Background()
	container, err := DeployTestContainer(ctx)

	if err != nil {
		log.Fatalln("failed to load container:", err)
	}

	connURL, err := container.ConnectionString(ctx, "sslmode=disable")

	return &DBTester{
		ConnUrl:   connURL,
		Parallel:  parallel,
		SleepTime: time.Millisecond * 500,
		TestData:  CreateDataTest(),
	}, err

}

func (ms *DBTester) GetConnection(ctx context.Context, opts ...any) (*pgxsql.Adapter, error) {
	var conn string
	if len(opts) == 0 {
		return pgxsql.NewAdapterByURL(ms.ConnUrl, ctx)
	}

	ops := opts[0].(string)
	if ops == "init" {
		conn = ms.ConnUrl
	} else {
		container, err := NewDBTester(ms.Parallel)
		if err != nil {
			return nil, err
		}
		conn = container.ConnUrl
		time.Sleep(5 * time.Second)
	}
	db, err := pgxsql.NewAdapterByURL(conn, ctx)
	if err != nil {
		return db, err
	}

	err = db.MigrateDatabase()
	if err != nil {
		return db, err

	}
	return db, err

}

func (ms *DBTester) CheckParallel(t *testing.T) {
	if ms.Parallel {
		t.Parallel()
	}
}

func CreateDataTest() TestData {
	userData := models.User{
		Email:               "suisei@hololive.com",
		Username:            "suisei",
		PasswordHash:        "123123",
		IsActive:            false,
		IsVerified:          false,
		CreatedAt:           time.Time{},
		UpdatedAt:           time.Time{},
		LastLogin:           &time.Time{},
		OAuthAccounts:       []models.OAuthAccount{},
		Sessions:            []models.UserSession{},
		VerificationTokens:  []models.VerificationToken{},
		PasswordResetTokens: []models.PasswordResetToken{},
		URLs:                []models.URL{},
	}
	return TestData{
		UserData: userData,
	}
}

func (ms *DBTester) SetupDataTest(opts ...any) {
	var db *pgxsql.Adapter
	var err error

	ctx := context.Background()
	if len(opts) == 0 {
		db, err = ms.GetConnection(ctx)
	} else {
		db = opts[0].(*pgxsql.Adapter)
	}
	if err != nil {
		log.Fatal("failed to connect to db: %w", err)
	}

	_ = db.CreateUser(ms.TestData.UserData)

}
