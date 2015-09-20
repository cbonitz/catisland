package tomcat

import (
	"fmt"
	"strings"
)

// An Application is something running in a Tomcat
type Application struct {
	host  string
	path  string
	state string
}

// createApplication parses a Tomcat manager status line into an Application struct
func createApplication(manager Manager, line string) (app *Application, err error) {
	fmt.Printf(line)
	parts := strings.Split(line, ":")
	return &Application{host: manager.Host, path: parts[0], state: parts[1]}, nil
}

func (a Application) String() string {
	return fmt.Sprintf("%s%s (%s)", a.host, a.path, a.state)
}
