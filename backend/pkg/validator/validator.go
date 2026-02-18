package validator

import "errors"

func Require(value string, field string) error {
	if value == "" {
		return errors.New(field + " is required")
	}
	return nil
}
