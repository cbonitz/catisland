package tomcat

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var okResponse = `OK - Listed applications for virtual host localhost
/:running:0:ROOT
/examples:running:0:examples
/host-manager:running:0:host-manager
/manager:running:1:manager
/docs:running:0:docs`

func TestNewManagerValidInput(t *testing.T) {
	// create form valid but un-trimmed line
	manager, err := NewManager(" http://example.com;user;password ")
	if manager == nil {
		t.Error("Constructor must return a value.")
	}
	if err != nil {
		t.Error("Creating a manager from a valid config line must not return an error")
	}
	if manager.password != "password" || manager.Host != "http://example.com" {
		t.Error("Line must be trimmed")
	}
	if manager.username != "user" {
		t.Error("User must be parsed")
	}
}

func TestMalformedLine(t *testing.T) {
	assertError(t, "this is not what we want", "malformed line")
}

func TestMissingArgs(t *testing.T) {
	assertError(t, ";;", "empty arguments")
	assertError(t, ";user;password", "empty arguments")
	assertError(t, "http://example.com;;password", "empty arguments")
	assertError(t, "http://example.com;user;", "empty arguments")
}

func TestTrimArgs(t *testing.T) {
	assertError(t, ";;", "empty arguments")
	assertError(t, " ;user;password", "empty arguments")
	assertError(t, "http://example.com; ;password", "empty arguments")
	assertError(t, "http://example.com;user; ", "empty arguments")
}

func TestGetterErrorHandling(t *testing.T) {
	manager, _ := NewManager("http://example.com;username;password")
	res, err := manager.GetStatus(func(m *Manager) (result string, err error) {
		return "", errors.New("oops")
	})
	if res != nil {
		t.Error("Manager should not return a result when loading data fails")
	}
	if err == nil {
		t.Error("Manager should return an error when loading data fails")
	}
	if !strings.Contains(err.Error(), "http://example.com") {
		t.Error("Error should contain hostname when data loading fails")
	}
}

func TestResultParsing(t *testing.T) {
	manager, _ := NewManager("http://example.com;username;password")
	res, err := manager.GetStatus(func(m *Manager) (result string, err error) {
		// success response as shown in
		// https://tomcat.apache.org/tomcat-8.0-doc/manager-howto.html#List_Currently_Deployed_Applications
		responseBody := okResponse
		return responseBody, nil
	})
	if err != nil {
		t.Error("Should not have gotten error on valid response: " + err.Error())
	}
	if len(res) != 5 {
		t.Errorf(
			"5 result lines should result in a result of length 5, but was %d",
			len(res))
	}
}

func TestErrorParsing(t *testing.T) {
	manager, _ := NewManager("http://example.com;username;password")
	// error response as described in
	// https://tomcat.apache.org/tomcat-8.0-doc/manager-howto.html#List_Currently_Deployed_Applications
	responseBody := `FAIL - something went wrong
and this is the reason`
	res, err := manager.GetStatus(func(m *Manager) (result string, err error) {
		return responseBody, nil
	})
	if err == nil {
		t.Error("Should have gotten error on response indicating failure.")
	}
	if !strings.Contains(err.Error(), responseBody) {
		t.Error("Error should be passed on")
	}
	if !strings.Contains(err.Error(), "http://example.com") {
		t.Error("Error message should contain hostname.")
	}
	if res != nil {
		t.Error("Should not have gotten a result on response indicating failure")
	}
}

func TestIntegrationSuccess(t *testing.T) {
	successServer := createServer(200, okResponse)
	defer successServer.Close()
	m, _ := NewManager(successServer.URL + ";user;password")
	apps, err := m.GetStatus(GetApplicationList)
	if err != nil {
		t.Error("Should have parsed result, but got error. ", err.Error())
	}
	if len(apps) != 5 {
		t.Error("Should have gotten 5 results but got ", len(apps))
	}
}

func TestIntegrationFail(t *testing.T) {
	successServer := createServer(404, "404 not found")
	defer successServer.Close()
	m, _ := NewManager(successServer.URL + ";user;password")
	_, err := m.GetStatus(GetApplicationList)
	if err == nil {
		t.Error("Should have gotten error for status 404")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Status code should be in error message, but was '%s'", err.Error())
	}
}

func createServer(statusCode int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, response)
	}))
}

func assertError(t *testing.T, line string, attempt string) {
	manager, err := NewManager(line)
	if manager != nil {
		t.Error("Constructor should return nil on " + attempt)
	}
	if err == nil {
		t.Error("Error must be returned on " + attempt)
	}
}
