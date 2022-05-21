package validator

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"tungnt/emmployee_manage/pkg/share/utils"
)

// gojsonschema handles the file handlers globally. So we have to cache
// this as singletons for the code to work properly in the tests
var cachedGoJsonSchemaInstances = map[string]*gojsonschema.Schema{}

type JsonSchemaValidator struct {
	basePath string
	schemas  map[string]*gojsonschema.Schema
}

func NewJsonSchemaValidator() (*JsonSchemaValidator, error) {
	validator := &JsonSchemaValidator{
		basePath: os.Getenv("DW_SCHEMAS_PATH") + "/" + os.Getenv("DW_SERVICE_NAME"),
		schemas:  make(map[string]*gojsonschema.Schema),
	}

	err := validator.loadDirSchemas("")
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (validator *JsonSchemaValidator) loadDirSchemas(path string) error {
	schemaFiles, err := os.ReadDir(validator.basePath + "/" + path)
	if err != nil {
		return err
	}

	for _, schemaFile := range schemaFiles {
		if schemaFile.Name() == ".gitkeep" {
			continue
		}

		schemaPath := path + "/" + schemaFile.Name()
		if schemaFile.IsDir() {
			if err = validator.loadDirSchemas(schemaPath); err != nil {
				return err
			}
			continue
		}

		goJsonSchemaPath := "file://" + validator.basePath + schemaPath
		var schema *gojsonschema.Schema
		var schemaExist bool
		if _, schemaExist = cachedGoJsonSchemaInstances[schemaPath]; !schemaExist {
			loader := gojsonschema.NewReferenceLoader(goJsonSchemaPath)
			schema, err = gojsonschema.NewSchema(loader)
			if err != nil {
				return err
			}
			cachedGoJsonSchemaInstances[schemaPath] = schema
		}
		validator.schemas[schemaPath] = schema
	}

	return nil
}

func (validator *JsonSchemaValidator) Validate(
	schemaFile string,
	data interface{},
) (*gojsonschema.Result, error) {
	if schemaFile[0] != '/' {
		schemaFile = "/" + schemaFile
	}

	if schema, existSchema := validator.schemas[schemaFile]; existSchema {
		result, err := schema.Validate(gojsonschema.NewGoLoader(data))
		if len(result.Errors()) == 0 {
			return nil, nil
		}

		return result, err
	} else {
		return nil, errors.New(fmt.Sprintf("schema not exist for file %s", schemaFile))
	}
}

func (validator *JsonSchemaValidator) GetErrorDetails(
	result gojsonschema.ResultError,
) map[string]interface{} {
	return map[string]interface{}{
		"context":     result.Context(),
		"description": result.Description(),
		"details":     result.Details(),
		"field":       result.Field(),
		"type":        result.Type(),
		"value":       result.Value(),
	}
}

func (validator *JsonSchemaValidator) GetErrorField(
	result gojsonschema.ResultError,
) string {
	field := result.Field()
	errorDetails := result.Details()
	if property, propertyExists := errorDetails["property"]; propertyExists {
		if propertyString, propertyIsString := property.(string); propertyIsString {
			field = field + "." + propertyString
		}
	}

	return field
}

func (validator *JsonSchemaValidator) GetCustomErrorMessage(
	result gojsonschema.ResultError,
) string {
	details := result.Details()
	format, formatExists := details["format"]
	min, minExists := details["min"]
	if result.Type() == "format" && formatExists {
		if format == "email" || format == "idn-email" {
			return utils.ErrorEmailFail
		}
		if format == "password" || format == "strong-password" {
			return utils.ErrorPasswordFail
		}
	}

	if result.Type() == "required" {
		return "MSGCM001"
	}

	if result.Type() == "string_gte" && minExists {
		if min == 1 {
			return utils.ErrorInputFail // Non-empty string
		}
	}

	_, maxExists := details["max"]
	if result.Type() == "string_lte" && maxExists {
		return utils.ErrorInputCharacterLimit
	}

	return utils.ErrorInputFail
}