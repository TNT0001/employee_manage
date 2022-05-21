package keycloak

import (
	"github.com/Nerzal/gocloak/v11"
	"github.com/sirupsen/logrus"
	"os"
)

type Client struct {
	Logger *logrus.Logger
	gocloak.GoCloak
}

func NewKeyCloakClient(logger *logrus.Logger) *Client {
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_ROOT_URL"))
	return &Client{
		Logger:  logger,
		GoCloak: client,
	}
}
