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
	fmt.Printf("Getting results from %d hosts.\n", len(hosts))
	// Get the running applications from the Tomcats
	errs := make(chan error, 100)
	results := make(chan []*tomcat.Application, 100)
	limiter := make(chan int, 5)
	go func() {
		for _, host := range hosts {
			limiter <- 1
			go func(host *tomcat.Manager) {
				hostApps, err := host.GetStatus(tomcat.GetApplicationList)
				if err != nil {
					errs <- err
				} else {
					results <- hostApps
				}
				<-limiter
			}(host)
		}
	}()
	appMap := make(map[string]*tomcat.Application)
	for i := range hosts {
		select {
		case err := <-errs:
			fmt.Printf("Error: %s\n", err.Error())
		case hostApps := <-results:
			for _, app := range hostApps {
				appRepresentation := app.String()
				if _, ok := appMap[appRepresentation]; !ok {
					appMap[appRepresentation] = app
				}
			}
		}
		if i%100 == 99 {
			fmt.Print(".")
		}
	}
	fmt.Println("Finished collecting results.")
	apps := []*tomcat.Application{}
	// Show them to the user
	for key, app := range appMap {
		fmt.Println(key)
		apps = append(apps, app)
	}
	fmt.Printf("%d apps total.\n", len(apps))
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
