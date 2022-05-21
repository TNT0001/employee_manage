package mount

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/container"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/handler"
	route "tungnt/emmployee_manage/internal/service/main_service/pkg/router"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/validator"
)

func MountAll(
	router *gin.Engine,
	logger *logrus.Logger,
	DB *gorm.DB,
	repositoryContainer *container.RepositoryContainer,
	serviceContainer *container.ServicesContainer,
) error {
	jsonSchemaValidator, err := validator.NewJsonSchemaValidator()
	if err != nil {
		return err
	}

	keycloakClient := keycloak.NewKeyCloakClient(logger)

	graphqlSchema, err := container.NewGraphQLSchema(repositoryContainer, serviceContainer, DB, logger, keycloakClient)
	if err != nil {
		log.Fatalln(err)
	}

	handler := handler.NewHandler(
		logger,
		DB,
		jsonSchemaValidator,
		graphqlSchema,
		repositoryContainer.UsersContainer,
		serviceContainer.AuthorizationsContainer,
		keycloakClient,
	)

	route.Mount(router, keycloakClient, handler)

	return nil
}
