# Cat Island - A Tomcat Application Overview Aggregator
*Cat Island* consolidates the application list of multiple Tomcat managers into a single JSON file.

According to Wikipedia, cats outnumber humans by a ratio of 6:1 on the island of [Aoshima](https://en.wikipedia.org/wiki/Aoshima,_Ehime).
If something similar can be said about Tomcat installations where you work, this tool might be for you.

## Features
* Consolidate applications from multiple tomcat servers into one list
* Parallel queries
* Optional JSON output

## Installation
`go get github.com/cbonitz/catisland`

## Usage
`catisland <config-file> [json-output-file]`
