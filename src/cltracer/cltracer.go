package main

import (
	"flag"
	"fmt"
	"util"
	"time"
	mrand "math/rand"
//	"regexp"
//	"strconv"
)

const APP_VERSION = "0.3"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

var fileToParse *string = flag.String("f", "", "File to parse")

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}
	
	if *fileToParse == "" {
		fmt.Println("No file to parse. Please provide a file to parse")
		return
	}
	
	content := util.Parse(*fileToParse)
	
	scene := util.ParseFile(content)
	
	today := time.Now()
	epoc := today.Unix()
	
	mrand.Seed(epoc)
	
	scene.Render(epoc)
	
	fmt.Println("Fin")
}
