package main

import "fmt"
import "os"
import "io/ioutil"
import "net/http"
import "strings"
import "errors"

// TomcatManager is data type for tomcat manager configuration
type TomcatManager struct {
	host     string
	username string
	password string
}

// String functon for TomcatManager
func (t TomcatManager) String() string {
	passwordOutput := t.password != ""
	return fmt.Sprintf("host '%s' username '%s' password? %t\n", t.host, t.username, passwordOutput)
}

// CreateManager creates a TomcatManager
func CreateManager(trimmedLine string) (result *TomcatManager, err error) {
	items := strings.Split(trimmedLine, ";")
	if len(items) != 3 {
		return nil, errors.New("lines must be formatted as hostname;user;password, but found " + trimmedLine)
	}
	config := TomcatManager{host: items[0] + "manager/text/list", username: items[1], password: items[2]}
	fmt.Printf("%s\n", config)
	return &config, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gokitty <config-file>")
		os.Exit(1)
	}

	configFile := os.Args[1]
	fmt.Println("Reading from config file", configFile)
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	var hosts []*TomcatManager
	for _, line := range lines {
		trimmed := strings.Trim(line, "\r")
		if len(trimmed) > 0 {
			config, err := CreateManager(trimmed)
			if err != nil {
				fmt.Println("Error: ", err.Error())
				os.Exit(1)
			} else {
				hosts = append(hosts, config)
			}
		}
	}

	client := &http.Client{}
	for _, host := range hosts {
		req, err := http.NewRequest("GET", host.host, nil)
		req.SetBasicAuth(host.username, host.password)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error on request.", err)
			os.Exit(1)
		}
		text, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(text))
	}
}
