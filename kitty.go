package main

import "fmt"
import "os"
import "io/ioutil"
import "strings"

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
	var hosts []map[string]string
	for _, line := range lines {
		trimmed := strings.Trim(line, "\r")
		if len(trimmed) > 0 {
			items := strings.Split(trimmed, ";")
			if len(items) != 3 {
				fmt.Println("lines must be formatted as hostname;user;password, but found", trimmed)
				os.Exit(1)
			}
			config := make(map[string]string)
			config["host"] = items[0]
			config["user"] = items[1]
			config["password"] = items[2]
			fmt.Printf("host '%s' user '%s' password '%s'\n", config["host"], config["user"], config["password"])
			hosts = append(hosts, config)
		}
	}
}
