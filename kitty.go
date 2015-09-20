package main

import (
	"fmt"
	"github.com/cbonitz/gokitty/tomcat"
	"io/ioutil"
	"os"
	"strings"
)

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
	var hosts []*tomcat.Manager
	for _, line := range lines {
		trimmed := strings.Trim(line, "\r")
		if len(trimmed) > 0 {
			config, err := tomcat.CreateManager(trimmed)
			if err != nil {
				fmt.Println("Error: ", err.Error())
				os.Exit(1)
			} else {
				hosts = append(hosts, config)
			}
		}
	}
	var apps []*tomcat.Application
	for _, host := range hosts {
		fmt.Printf("Getting status for %s\n", host)
		hostApps, err := host.GetStatus()
		if err != nil {
			cat := *host
			fmt.Printf("Error on host %s: %s\n", cat.Host, err.Error())
		} else {
			apps = append(apps, hostApps...)
		}
	}
	for _, app := range apps {
		fmt.Println(app)
	}
}
