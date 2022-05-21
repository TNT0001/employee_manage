package handler

import "github.com/sirupsen/logrus"

// ApplicationHTTPHandler base handler struct.
type ApplicationHTTPHandler struct {
	BaseHTTPHandler
}

// NewApplicationHTTPHandler returns ApplicationHTTPHandler instance.
func NewApplicationHTTPHandler(logger *logrus.Logger) *ApplicationHTTPHandler {
	return &ApplicationHTTPHandler{BaseHTTPHandler: BaseHTTPHandler{Logger: logger}}
}
