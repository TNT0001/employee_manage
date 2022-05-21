package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"

	baseDto "tungnt/emmployee_manage/pkg/dto"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/validator"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

// BaseHTTPHandler base handler struct.
type BaseHTTPHandler struct {
	Logger    *logrus.Logger
	Validator *validator.JsonSchemaValidator
}

func (h BaseHTTPHandler) GetInputAsMap(c *gin.Context) (map[string]interface{}, error) {
	if c.ContentType() != "application/json" {
		return nil, wraperror.NewApiDisplayableError(
			http.StatusBadRequest, map[string]interface{}{
				"content_type": "expected content_type is application/json",
			}, errors.New("wrong content type"))
	}

	input := make(map[string]interface{})
	err := c.ShouldBindJSON(&input)

	return input, err
}

func (h *BaseHTTPHandler) SetGenericErrorResponse(c *gin.Context, finalError error) {
	// Retrieving the original error inside GraphQL's wrapper if there is one
	// If there is none, we keep the error coming from the graphql's engine
	originalError := finalError
	if _, ok := originalError.(gqlerrors.FormattedError); ok {
		err := originalError.(gqlerrors.FormattedError).OriginalError()
		if err != nil {
			originalError = err
		}

		if _, ok := originalError.(*gqlerrors.Error); ok {
			err := originalError.(*gqlerrors.Error).OriginalError
			if err != nil {
				originalError = err
			}
		}
	}

	apiError := &wraperror.ApiDisplayableError{}
	jsonError := &json.SyntaxError{}
	if errors.As(originalError, &apiError) {
		var debugInfo interface{} = nil
		if os.Getenv("DW_DEBUG") == "TRUE" {
			debugInfo = finalError
		}

		data := baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message:          apiError.Message(),
				DebugInformation: debugInfo,
			},
		}
		c.JSON(apiError.HttpStatus(), data)
		return
	} else if errors.Is(originalError, gorm.ErrRecordNotFound) || originalError.Error() == gorm.ErrRecordNotFound.Error() {
		data := baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message: originalError.Error(),
			},
		}
		c.JSON(http.StatusNotFound, data)
		return
	} else if errors.As(originalError, &jsonError) {
		data := &baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message: "Invalid json",
				Details: map[string]interface{}{
					"offset": jsonError.Offset,
					"error":  jsonError.Error(),
				},
			},
		}

		c.JSON(http.StatusBadRequest, data)
		return
	} else if strings.Contains(originalError.Error(), utils.KeycloakInvalidGrant) {
		data := baseDto.BaseErrorResponse{
			Error: &baseDto.ErrorResponse{
				Message: originalError.Error(),
			},
		}
		c.JSON(http.StatusBadRequest, data)
		return
	} else {
		h.SetInternalErrorResponse(c, finalError)
		return
	}
}

// This outputs a 500 error with a custom message (contrary to SetGenericErrorResponse that hides the real unhandled error)
func (h *BaseHTTPHandler) SetInternalErrorResponse(c *gin.Context, err error) {
	var debugInfo interface{} = nil
	if os.Getenv("DW_DEBUG") == "TRUE" {
		debugInfo = err
	} else {
		h.Logger.Errorf("Unexpected error: %+v", err)
	}

	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message:          utils.MessageInternalServerError,
			DebugInformation: debugInfo,
		},
	}

	c.JSON(http.StatusInternalServerError, data)
}

// This outputs a 400 error with a custom message (contrary to SetGenericErrorResponse that hides the real unhandled error)
func (h *BaseHTTPHandler) SetBadErrorResponse(c *gin.Context, message interface{}) {
	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message: message,
		},
	}

	c.JSON(http.StatusInternalServerError, data)
}

func (h *BaseHTTPHandler) SetJSONValidationErrorResponse(
	c *gin.Context,
	validationResults *gojsonschema.Result,
) {
	h.SetJSONValidationWithCustomErrorResponse(
		c,
		validationResults,
		func(result gojsonschema.ResultError) string {
			return ""
		},
	)
}

func (h *BaseHTTPHandler) SetJSONValidationWithCustomErrorResponse(
	c *gin.Context,
	validationResults *gojsonschema.Result,
	getError func(result gojsonschema.ResultError) string,
) {
	messages := map[string]string{}
	details := make([]map[string]interface{}, 0)
	for _, validationError := range validationResults.Errors() {
		field := h.Validator.GetErrorField(validationError)
		detail := h.Validator.GetErrorDetails(validationError)

		// Getting either a message defined by the handler,
		// or the default customized one
		message := getError(validationError)
		if message == "" {
			message = h.Validator.GetCustomErrorMessage(validationError)
		}

		messages[field] = message
		details = append(details, detail)
	}

	data := &baseDto.BaseErrorResponse{
		Error: &baseDto.ErrorResponse{
			Message: messages,
			Details: details,
		},
	}

	c.JSON(http.StatusBadRequest, data)
}
