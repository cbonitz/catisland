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

// StringGetter is a function that returns a string or error based on a
// Manager struct, e.g. via an http call
type StringGetter func(t *Manager) (result string, err error)

// NewManager creates a Manager
func NewManager(line string) (result *Manager, err error) {
	trimmedLine := strings.TrimSpace(line)
	items := strings.Split(trimmedLine, ";")
	if len(items) != 3 {
		return nil, errors.New("lines must be formatted as hostname;user;password, but found " + trimmedLine)
	}
	host := strings.TrimSpace(items[0])
	username := strings.TrimSpace(items[1])
	password := strings.TrimSpace(items[2])
	if len(host) == 0 || len(username) == 0 || len(password) == 0 {
		return nil, errors.New("host, username and passwords must be nonempty")
	}
	config := Manager{Host: host, username: username, password: password}
	return &config, nil
}

func (t Manager) String() string {
	return fmt.Sprintf("host '%s' username '%s'", t.Host, t.username)
}

// GetApplicationList gets a tring from the tomcat manager Application
// list
func GetApplicationList(t *Manager) (result string, err error) {
	// Tomcat manager wants a GET call with basic auth to a particular URL
	client := &http.Client{}
	req, err := http.NewRequest("GET", t.Host+"/manager/text/list", nil)
	req.SetBasicAuth(t.username, t.password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// It returns a status code 200 on success
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Got status '%s'", resp.Status)
	}
	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(rawBody), nil
}

// GetStatus gets the status of a Manager
func (t Manager) GetStatus(getter StringGetter) (result []*Application, err error) {
	text, err := getter(&t)
	if err != nil {
		return nil, fmt.Errorf(
			"While getting result from %s: %s",
			t.Host, err.Error())
	}
	lines := strings.Split(text, "\n")

	// The first line starts with OK
	if !strings.HasPrefix(lines[0], "OK") {
		return nil, fmt.Errorf("Non-OK response from %s: %s", t.Host, text)
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
