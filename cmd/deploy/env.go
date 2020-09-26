package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("./app.deploy.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := string(b)
	projectID := os.Args[1]
	yaml := strings.Replace(lines, "##PROJECT_ID", projectID, 1)
	err = ioutil.WriteFile("./app.yaml", []byte(yaml), 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
