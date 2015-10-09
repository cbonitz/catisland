package main

import (
	"encoding/json"
	"fmt"
	"github.com/cbonitz/catisland/tomcat"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	// Parse trivial command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: catisland <config-file> [json-file]")
		os.Exit(1)
	}

	// Read the configuration
	configFile := os.Args[1]
	fmt.Println("Reading from config file", configFile)
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	// Each line should describe a tomcat manager
	lines := strings.Split(string(content), "\n")
	var hosts []*tomcat.Manager
	for _, line := range lines {
		trimmed := strings.Trim(line, "\r")
		if len(trimmed) > 0 {
			config, err := tomcat.NewManager(trimmed)
			if err != nil {
				fmt.Println("Error: ", err.Error())
				os.Exit(1)
			} else {
				hosts = append(hosts, config)
			}
		}
	}

	// Get the running applications from the Tomcats
	var apps []*tomcat.Application
	for _, host := range hosts {
		fmt.Printf("Getting status for %s\n", host)
		hostApps, err := host.GetStatus(tomcat.GetApplicationList)
		if err != nil {
			cat := *host
			fmt.Printf("Error on host %s: %s\n", cat.Host, err.Error())
		} else {
			apps = append(apps, hostApps...)
		}
	}

	// Show them to the user
	for _, app := range apps {
		fmt.Println(app)
	}

	// JSON-serialize to file, if desired
	if len(os.Args) == 3 {
		jsonSerialized, _ := json.MarshalIndent(apps, "", "  ")
		err := ioutil.WriteFile(os.Args[2], []byte(jsonSerialized), 0777)
		if err != nil {
			fmt.Printf("Error writing JSON: %s\n", err)
			os.Exit(1)
		}
	}
}
