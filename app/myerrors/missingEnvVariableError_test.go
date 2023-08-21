package myerrors

import "testing"

func TestIs(t *testing.T) {
	missingEnvVarError := MissingEnvVariableError{"SUPERSECRET"}
	if missingEnvVarError.Error() != "Missing SUPERSECRET environment variable" {
		t.Error("Couldn't create an error with SUPERSECRET")
	}
}
