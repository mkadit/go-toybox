package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkadit/go-toybox/common/config"
	logfile "github.com/mkadit/go-toybox/internal/logger"
	"github.com/mkadit/go-toybox/internal/models"
)

// Adapter implements the DbPort interface
type Adapter struct {
	db         *pgxpool.Pool
	ctx        context.Context
	dbSource   string
	dbName     string
	dbHost     string
	dbUser     string
	dbPassword string
	dbPort     string
}

// NewAdapter creates a new Adapter
func NewAdapter(conf models.DbConfiguration, ctx context.Context) (*Adapter, error) {
	logfile.LogEvent("connecting to db")
	source := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	conn, err := pgxpool.New(ctx, source)
	if err != nil {
		return &Adapter{}, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return &Adapter{}, err
	}

	logfile.LogEvent("connected to db")
	return &Adapter{
		db:         conn,
		ctx:        ctx,
		dbSource:   source,
		dbName:     conf.Database,
		dbHost:     conf.Host,
		dbUser:     conf.User,
		dbPassword: conf.Password,
		dbPort:     conf.Port,
	}, nil
}

// NewTestAdapterMSQL Msql secondary adapter using repository pattert for testing
func NewAdapterByURL(connURL string, ctx context.Context) (*Adapter, error) {
	conn, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return &Adapter{}, models.ErrorConnectDB

	}
	return &Adapter{
		db:       conn,
		dbSource: connURL,
		ctx:      ctx,
	}, nil
}

func (ad Adapter) MigrateDatabase() error {
	var dsn string
	if ad.dbSource != "" {
		// dsn = fmt.Sprintf("pgx5%s", ad.dbSource[8:])
		parts := strings.SplitN(ad.dbSource, "://", 2)
		parts[0] = "pgx5"
		dsn = strings.Join(parts, "://")
	} else {

		dsn = fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable", ad.dbUser, ad.dbPassword, ad.dbHost, ad.dbPort, ad.dbName)
	}

	migrationsLoc := "file://" + config.ProjectRootPath + "/common/schema/"
	m, err := migrate.New(migrationsLoc, dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err.Error() != models.ErrMigrateNoChange.Error() {
		if strings.Contains(err.Error(), models.ErrDirtyMigration.Error()) {
			err1 := m.Drop()
			if err1 != nil {
				err = errors.Join(err, err1)
			}
		}
		return err
	}
	return nil
}

func (ad Adapter) DropMigration() error {
	var dsn string
	if ad.dbSource != "" {
		// dsn = fmt.Sprintf("pgx5%s", ad.dbSource[8:])
		parts := strings.SplitN(ad.dbSource, "://", 2)
		parts[0] = "pgx5"
		dsn = strings.Join(parts, "://")
	} else {

		dsn = fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable", ad.dbUser, ad.dbPassword, ad.dbHost, ad.dbPort, ad.dbName)
	}

	migrationsLoc := "file://" + config.ProjectRootPath + "/common/schema/"
	m, err := migrate.New(migrationsLoc, dsn)
	if err != nil {
		return err
	}
	err = m.Drop()
	if err != nil {
		return err
	}
	return nil
}

func (ad Adapter) CloseDbConnection() {
	ad.db.Close()
}

func (ad Adapter) GetDB() *pgxpool.Pool {
	return ad.db

}

// InsertDataTest Insert test data from /common/script/data.sql
func (d Adapter) InsertDataTest() error {
	sqlFile := config.ProjectRootPath + "/common/script/data.sql"
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		return err
	}

	// Split the SQL content into lines
	lines := strings.Split(string(sqlContent), "\n")

	tx, err := d.db.Begin(d.ctx)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback(d.ctx)
	}()

	var currentStatement strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and full-line comments
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		// Add line to current statement
		currentStatement.WriteString(line + " ")

		// Check if the statement is complete (ends with a semicolon)
		if strings.HasSuffix(line, ";") {
			stmt := strings.TrimSpace(currentStatement.String())

			// Execute the statement
			if strings.HasPrefix(stmt, "INSERT") {
				_, err := tx.Exec(d.ctx, stmt)
				if err != nil {
					return fmt.Errorf("error executing statement %q: %v", stmt, err)
				}
			} else {
				log.Printf("Unrecognized statement: %s\n", stmt)
			}

			// Reset the current statement
			currentStatement.Reset()
		}
	}

	if err = tx.Commit(d.ctx); err != nil {
		return err
	}
	return nil
}
