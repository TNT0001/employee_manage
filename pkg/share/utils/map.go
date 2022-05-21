package utils

import "github.com/mitchellh/mapstructure"

func MapToStruct(
	input map[string]interface{},
	output interface{},
) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		ZeroFields:  true,
		Result:      output,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
