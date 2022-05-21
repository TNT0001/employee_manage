package migrate

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"tungnt/emmployee_manage/pkg/infrastructure"
)

var migrateTable = map[string]string{
	"users":          "schema_migrations_users",
	"jobs":           "schema_migrations_jobs",
	"authorizations": "schema_migrations_authorizations",
	"organizations":  "schema_migrations_organizations",
}

func MigrateAll(
	logger *logrus.Logger,
	masterConn *sql.DB,
) error {
	for modules, tableName := range migrateTable {
		if err := infrastructure.Migrate(
			masterConn,
			logger,
			modules,
			tableName,
		); err != nil {
			logger.Infof("Fail to migrate modules: %s, Error: %s", modules, err.Error())
			return err
		}
	}
	return nil
}

func CheckVersion(
	logger *logrus.Logger,
	masterConn *sql.DB,
) error {
	for modules, tableName := range migrateTable {
		if err := infrastructure.CheckDatabaseVersion(
			masterConn,
			logger,
			modules,
			tableName,
		); err != nil {
			logger.Infof("The migrations for module %v are not up-to-date: %s", modules, err)
			return err
		}
	}
	return nil
}
