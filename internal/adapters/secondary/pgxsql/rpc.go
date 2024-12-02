package pgxsql

import (
	"fmt"

	"github.com/mkadit/go-toybox/internal/models"
)

func (ad *Adapter) RPCDropMigration(args *models.RPCMigrationArgs, reply *string) error {
	fmt.Println("DROP MIGRATION")
	if err := ad.DropMigration(); err != nil {
		return err
	}
	*reply = "Database migration dropped successfully"
	return nil
}

func (ad *Adapter) RPCMigrateDatabase(args *models.RPCMigrationArgs, reply *string) error {
	if err := ad.MigrateDatabase(); err != nil {
		return err
	}
	*reply = "Database migration completed successfully"
	return nil
}
