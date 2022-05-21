package handler

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	authContainer "tungnt/emmployee_manage/internal/modules/authorizations/pkg/container"
	"tungnt/emmployee_manage/internal/modules/users/pkg/container"
	"tungnt/emmployee_manage/pkg/share/handler"
	"tungnt/emmployee_manage/pkg/share/keycloak"
	"tungnt/emmployee_manage/pkg/share/validator"
)

type Handler struct {
	handler.BaseHTTPHandler
	Graphql              graphql.Schema
	DB                   *gorm.DB
	UserRepoContainer    container.UserContainer
	AuthServiceContainer authContainer.AuthorizationServiceContainer
	Client               *keycloak.Client
}

func NewHandler(
	logger *logrus.Logger,
	db *gorm.DB,
	schemaValidator *validator.JsonSchemaValidator,
	graphql graphql.Schema,
	userRepoContainer container.UserContainer,
	authServiceContainer authContainer.AuthorizationServiceContainer,
	client *keycloak.Client,
) *Handler {
	baseHanler := handler.BaseHTTPHandler{
		Logger:    logger,
		Validator: schemaValidator,
	}

	return &Handler{
		baseHanler,
		graphql,
		db,
		userRepoContainer,
		authServiceContainer,
		client,
	}
}
