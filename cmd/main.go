package main

import (
	"fmt"

	"github.com/ppp3ppj/go-refactoring-workshop/config"
)

func main() {
    conf := config.ConfigGetting()
    _ = conf
    fmt.Printf("%s %s", conf.AppInfo.Name, conf.Database.Host)
    fmt.Println("Hello world")
}
