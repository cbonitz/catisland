package tomcat

import (
	"fmt"
	"strings"
)

// An Application is something running in a Tomcat
type Application struct {
	Host  string
	Path  string
	State string
}

// createApplication parses a Tomcat manager status line into an Application struct
func createApplication(manager Manager, line string) (app *Application, err error) {
	parts := strings.Split(line, ":")
	return &Application{Host: manager.Host, Path: parts[0], State: parts[1]}, nil
}

func (a Application) String() string {
	return fmt.Sprintf("%s%s (%s)", a.Host, a.Path, a.State)
}
