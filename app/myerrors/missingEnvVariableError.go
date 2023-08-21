package myerrors

import "fmt"

type MissingEnvVariableError struct {
	variableName string
}

func (e *MissingEnvVariableError) Error() string {
	return fmt.Sprintf("Missing %s environment variable", e.variableName)
}

func NewMissingEnvVariableError(name string) error {
	return &MissingEnvVariableError{name}
}
