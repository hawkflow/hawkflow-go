package hawkflow

import "fmt"

func createError(msg string) error {
	return fmt.Errorf("%s %s", msg, "Please see documentation at https://docs.hawkflow.ai/integration/index.html")
}
