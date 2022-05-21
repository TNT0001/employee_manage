package infrastructure

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source"

	// import source file
	_ "github.com/golang-migrate/migrate/v4/source/file"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func GetGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func NewDatabase(logger *logrus.Logger) (*gorm.DB, *sql.DB, *sql.DB, error) {
	gormConfig := GetGormConfig()
	gormConfig.Logger = NewGormLogger(logger)

	masterConnection, err := sql.Open("mysql", os.Getenv("MASTER_CONNECTION_STRING"))
	if err != nil {
		return nil, nil, nil, err
	}

	slaveConnection, err := sql.Open("mysql", os.Getenv("SLAVE_CONNECTION_STRING"))
	if err != nil {
		return nil, nil, nil, err
	}

	masterConnection.SetConnMaxIdleTime(5 * time.Minute)
	slaveConnection.SetConnMaxIdleTime(5 * time.Minute)

	masterConnection.SetConnMaxLifetime(20 * time.Minute)
	slaveConnection.SetConnMaxLifetime(20 * time.Minute)

	masterConnection.SetMaxIdleConns(50)
	slaveConnection.SetMaxIdleConns(50)

	masterConnection.SetMaxOpenConns(100)
	slaveConnection.SetMaxOpenConns(100)

	masterDriver := mysqlDriver.New(mysqlDriver.Config{
		Conn:                      masterConnection,
		SkipInitializeWithVersion: true,
		DontSupportForShareClause: true,
	})

	slaveDriver := mysqlDriver.New(mysqlDriver.Config{
		Conn:                      slaveConnection,
		SkipInitializeWithVersion: true,
		DontSupportForShareClause: true,
	})

	db, err := gorm.Open(masterDriver, gormConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	resolver := dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{slaveDriver},
	})

	err = db.Use(resolver)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = ping(db); err != nil {
		return nil, nil, nil, err
	}

	return db, masterConnection, slaveConnection, nil
}

func ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func CloseDB(logger *logrus.Logger, masterConn, slaveConn *sql.DB) error {
	logger.Info("Closing master connection")
	if err := masterConn.Close(); err != nil {
		logger.Errorf("Fail to close master connection")
		return err
	}
	logger.Info("Closing master connection done!")

	logger.Info("Closing slave connection")
	if err := slaveConn.Close(); err != nil {
		logger.Errorf("Fail to close slave connection")
		return err
	}
	logger.Info("Closing slave connection done!")

	return nil
}

func logStat(entry *logrus.Entry, stats sql.DBStats) {
	entry.WithFields(logrus.Fields{
		"InUser":   stats.InUse,
		"Idle":     stats.Idle,
		"OpenConn": stats.OpenConnections,
	}).Info()
}

func LogPoolSize(db *sql.DB, logger *logrus.Entry) func() {
	stop := make(chan struct{})
	curStat := db.Stats()
	go func() {
		for {
			select {
			case <-stop:
				logStat(logger, db.Stats())
				return
			case <-time.After(1 * time.Second):
				stat := db.Stats()
				if curStat.InUse != stat.InUse || curStat.OpenConnections != stat.OpenConnections || curStat.Idle != stat.Idle {
					logStat(logger, stat)
					curStat = stat
				}
				time.Sleep(1 * time.Second)
			}

		}
	}()

	return func() {
		stop <- struct{}{}
	}
}

func Migrate(
	master *sql.DB,
	logger *logrus.Logger,
	module string,
	migrationTable string,
) error {
	connection, err := master.Conn(context.Background())
	if err != nil {
		return err
	}

	mysql.DefaultMigrationsTable = migrationTable
	driver, err := mysql.WithConnection(
		context.Background(),
		connection,
		&mysql.Config{},
	)
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(getMigrationsPath(module), "mysql", driver)
	if err != nil {
		return err
	}
	m.Log = NewMigrateLogger(logger)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func getMigrationsPath(module string) string {
	return "file://" + os.Getenv("EM_MIGRATION_PATH") + "/" + module
}

func CheckDatabaseVersion(
	masterDb *sql.DB,
	logger *logrus.Logger,
	module string,
	migrationsTable string,
) error {
	connection, err := masterDb.Conn(context.Background())
	if err != nil {
		return err
	}

	mysql.DefaultMigrationsTable = migrationsTable
	driver, err := mysql.WithConnection(
		context.Background(),
		connection,
		&mysql.Config{},
	)
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(getMigrationsPath(module), "mysql", driver)
	if err != nil {
		return err
	}
	m.Log = NewMigrateLogger(logger)

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}
	if dirty {
		return migrate.ErrDirty{
			Version: int(version),
		}
	}

	fileMigration, err := source.Open(getMigrationsPath(module))
	if err != nil {
		return err
	}

	if _, err := fileMigration.Next(version); err == os.ErrNotExist {
		// no new version, return
		return nil
	}

	return connection.Close()
}
