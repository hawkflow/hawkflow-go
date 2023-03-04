package hawkflow

import (
	"fmt"
	"regexp"
)

func validateApiKey(apiKey string) error {
	if apiKey == "" {
		return createError("No API Key set.")
	}

	if len(apiKey) > 50 {
		return createError("Invalid API Key format.")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", apiKey); m == false {
		return createError("Invalid API Key format.")
	}

	return nil
}

func validateTime(r *request) error {
	if err := validateProcess(r.Process); err != nil {
		return err
	}
	if err := validateMeta(r.Meta); err != nil {
		return err
	}
	if err := validateUID(r.UID); err != nil {
		return err
	}

	return nil
}

func validateException(r *request) error {
	if err := validateProcess(r.Process); err != nil {
		return err
	}
	if err := validateMeta(r.Meta); err != nil {
		return err
	}
	if err := validateExceptionMessage(r.ExceptionMessage); err != nil {
		return err
	}

	return nil
}

func validateMetric(r *request) error {
	if err := validateProcess(r.Process); err != nil {
		return err
	}
	if err := validateMeta(r.Meta); err != nil {
		return err
	}
	if err := validateMetricItems(r.Items); err != nil {
		return err
	}

	return nil
}

func validateProcess(process string) error {
	if process == "" {
		return createError("No process set.")
	}

	if len(process) > 250 {
		return createError("Process parameter exceeded max length of 250 characters.")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", process); m == false {
		return createError("Process parameter contains unsupported characters.")
	}

	return nil
}

func validateMeta(meta string) error {
	if len(meta) > 500 {
		return createError("Meta parameter exceeded max length of 500 characters.")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z0-9_-]*$", meta); m == false {
		return createError("Meta parameter contains unsupported characters.")
	}

	return nil
}

func validateUID(uid string) error {
	if len(uid) > 50 {
		return createError("UID parameter exceeded max length of 50 characters.")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z0-9_-]*$", uid); m == false {
		return createError("UID parameter contains unsupported characters.")
	}

	return nil
}

func validateExceptionMessage(exceptionMessage string) error {
	if len(exceptionMessage) > 15000 {
		return createError("ExceptionMessage parameter exceeded max length of 15000 characters.")
	}

	return nil
}

func validateMetricItems(items map[string]float64) error {
	if len(items) == 0 {
		return createError("No items set.")
	}

	for k := range items {
		if len(k) > 50 {
			return createError(fmt.Sprintf("Item key %s exceeded max length of 50 characters.", k))
		}
	}

	return nil
}
