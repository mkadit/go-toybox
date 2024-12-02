package main

import (
	"context"
	"net/rpc"
	"time"

	"github.com/mkadit/go-toybox/common/config"
	"github.com/mkadit/go-toybox/internal/adapters/primary/http"
	"github.com/mkadit/go-toybox/internal/adapters/primary/rpc_server"
	"github.com/mkadit/go-toybox/internal/adapters/secondary/pgxsql"
	"github.com/mkadit/go-toybox/internal/applications/api"
	"github.com/mkadit/go-toybox/internal/applications/core/urls"
	"github.com/mkadit/go-toybox/internal/applications/core/users"
	logfile "github.com/mkadit/go-toybox/internal/logger"
	"github.com/mkadit/go-toybox/internal/models"
)

// Name description
func main() {
	go logfile.CreateLogger()
	time.Sleep(1 * time.Second)
	ctx := context.Background()
	conf, err := config.LoadConfig(config.ProjectRootPath)
	if err != nil {
		go logfile.LogFatal(err, config.ErrLoadEnv.Error())
	}

	// DB
	dbAdapter, err := pgxsql.NewAdapter(conf.Db, ctx)
	if err != nil {
		go logfile.LogFatal(err, models.ErrorConnectDB.Error())
	}

	logfile.LogEvent("migrating db")
	err = dbAdapter.MigrateDatabase()
	if err != nil {
		go logfile.LogErrorEvent(err, models.ErrMigrate.Error())

	}

	dbAdapter.InsertDataTest()
	defer func() {
		logfile.LogEvent("closing db connection")
		dbAdapter.CloseDbConnection()
	}()

	// CORE
	usersCore := users.New(conf.Auth.JwtSecret)
	urlsCore := urls.New()

	appAPI := api.NewApplication(dbAdapter, usersCore, urlsCore, conf.Email)
	_ = appAPI

	c := http.NewAdapter(appAPI, conf)

	// Start the RPC server in a separate goroutine
	err = rpc.RegisterName("DbAdapter", dbAdapter)
	if err != nil {
		go logfile.LogErrorEvent(err, models.ErrorConnectDB.Error())
	}
	go rpc_server.Start(conf.Rpc)

	c.SetupRouter()

	err = c.Setup()
	if err != nil {
		logfile.LogFatal(err, models.ErrListenAndServer.Error())
	}
	// _ = c

}
