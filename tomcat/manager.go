package tomcat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Manager is data type for tomcat manager configuration
type Manager struct {
	Host     string
	username string
	password string
}

// CreateManager creates a Manager
func CreateManager(trimmedLine string) (result *Manager, err error) {
	items := strings.Split(trimmedLine, ";")
	if len(items) != 3 {
		return nil, errors.New("lines must be formatted as hostname;user;password, but found " + trimmedLine)
	}
	config := Manager{Host: items[0], username: items[1], password: items[2]}
	return &config, nil
}

func (t Manager) String() string {
	passwordOutput := t.password != ""
	return fmt.Sprintf("host '%s' username '%s' password? %t", t.Host, t.username, passwordOutput)
}

// GetStatus gets the status of a Manager
func (t Manager) GetStatus() (result []*Application, err error) {
	// Tomcat manager wants a GET call with basic auth to a particular URL
	client := &http.Client{}
	req, err := http.NewRequest("GET", t.Host+"/manager/text/list", nil)
	req.SetBasicAuth(t.username, t.password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// It returns a status code 200 on success
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got status %d - %s", resp.StatusCode, resp.Status)
	}
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	text := string(rawBody)
	lines := strings.Split(text, "\n")

	// The first line starts with OK
	if !strings.HasPrefix(lines[0], "OK") {
		return nil, errors.New("Got non-OK response: " + lines[0])
	}

	// Parse applications and return them
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		app, err := createApplication(t, line)
		if err != nil {
			return nil, err
		}
		result = append(result, app)
	}
	return result, nil
}
