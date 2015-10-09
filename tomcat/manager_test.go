package tomcat

import (
	"testing"
)

func Test_CreateManager(t *testing.T) {
	// create form valid but un-trimmed line
	manager, err := CreateManager(" http://localhost:8080;user;password ")
	if manager == nil {
		t.Error("Constructor must return a value.")
	}
	if err != nil {
		t.Error("Creating a manager from a valid config line must not return an error")
	}
	if manager.password != "password" || manager.Host != "http://localhost:8080" {
		t.Error("Line must be trimmed")
	}
	if manager.username != "user" {
		t.Error("User must be parsed")
	}
}

func Test_malformed_line(t *testing.T) {
	assertError(t, "this is not what we want", "malformed line")
}

func Test_missing_args(t *testing.T) {
	assertError(t, ";;", "empty arguments")
	assertError(t, ";user;password", "empty arguments")
	assertError(t, "http://localhost:8080;;password", "empty arguments")
	assertError(t, "http://localhost:8080;user;", "empty arguments")
}

func Test_trim_args(t *testing.T) {
	assertError(t, ";;", "empty arguments")
	assertError(t, " ;user;password", "empty arguments")
	assertError(t, "http://localhost:8080; ;password", "empty arguments")
	assertError(t, "http://localhost:8080;user; ", "empty arguments")
}

func assertError(t *testing.T, line string, attempt string) {
	manager, err := CreateManager(line)
	if manager != nil {
		t.Error("Constructor should return nil on " + attempt)
	}
	if err == nil {
		t.Error("Error must be returned on " + attempt)
	}
}
