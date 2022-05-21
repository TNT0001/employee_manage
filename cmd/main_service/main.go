package main

import (
	"log"
	"net/http"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/container"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/migrate"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/mount"
	"tungnt/emmployee_manage/pkg/infrastructure"
)

func main() {
	logger := infrastructure.NewLogger()
	DB, masterConn, slaveConn, err := infrastructure.NewDatabase(logger)
	if err != nil {
		log.Fatalln(err)
	}

	dbLogEntry := logger.WithField("service", "database")
	stopMasterLogger := infrastructure.LogPoolSize(masterConn, dbLogEntry.WithField("pool", "master"))
	defer stopMasterLogger()
	stopSlaveLogger := infrastructure.LogPoolSize(slaveConn, dbLogEntry.WithField("pool", "slave"))
	defer stopSlaveLogger()

	err = migrate.MigrateAll(logger, masterConn)
	if err != nil {
		logger.Fatalf("Failed to migrate the database: %v", err)
		return
	}

	err = migrate.CheckVersion(logger, masterConn)
	if err != nil {
		logger.Fatalf("The database is not up-to-date: %v", err)
		return
	}
	router := infrastructure.NewServer(logger, DB)

	repositoryContainer := container.NewRepositoryContainer(logger, DB)

	serviceContainer := container.NewServicesContainer(logger, DB, repositoryContainer)

	err = mount.MountAll(router, logger, DB, &repositoryContainer, &serviceContainer)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info("Start server at localhost:8088")
	if err := http.ListenAndServe(":8088", router); err != nil {
		log.Fatalln(err)
	}
}
